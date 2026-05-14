package http

import (
	"dss/internal/domain/analytics"
	"dss/internal/domain/entities"
	"dss/internal/repositories"
	"encoding/csv"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ImportHandler struct {
	critRepo  repositories.CriterionRepository
	evalRepo  repositories.EvaluationRepository
	votingSvc *analytics.VotingService
	expertSvc *analytics.ExpertEvaluationService
}

func NewImportHandler(critRepo repositories.CriterionRepository, evalRepo repositories.EvaluationRepository, votingSvc *analytics.VotingService, expertSvc *analytics.ExpertEvaluationService) *ImportHandler {
	return &ImportHandler{
		critRepo:  critRepo,
		evalRepo:  evalRepo,
		votingSvc: votingSvc,
		expertSvc: expertSvc,
	}
}

func (h *ImportHandler) ImportVoting(c *gin.Context) {
	importFilePath := "last_votes.csv"
	var records [][]string

	fileHeader, err := c.FormFile("import_file")
	var originalFilename string

	if err == nil {
		// Новий файл був завантажений
		originalFilename = fileHeader.Filename
		if err := c.SaveUploadedFile(fileHeader, importFilePath); err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": "Could not save uploaded file"})
			return
		}
	} else {
		// Спроба використати попередньо збережений файл
		originalFilename = c.PostForm("original_filename")
		if originalFilename == "" {
			originalFilename = "previous_upload.csv"
		}
	}

	methodStr := c.PostForm("voting_method")
	if methodStr == "" {
		methodStr = "borda"
	}

	importFile, err := os.Open(importFilePath)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"title": "Error", "message": "No file uploaded and no previous file found"})
		return
	}
	defer importFile.Close()

	reader := csv.NewReader(importFile)
	records, err = reader.ReadAll()
	if err != nil || len(records) < 2 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"title": "Error", "message": "Invalid CSV file"})
		return
	}

	rankings := make(map[string][]int64)
	for i, row := range records {
		if i == 0 {
			continue // Skip header
		}
		if len(row) < 2 {
			continue
		}
		expertID := row[0]
		var rankList []int64
		for _, valStr := range row[1:] {
			valStr = strings.TrimSpace(valStr)
			if valStr == "" {
				continue
			}
			critID, err := strconv.ParseInt(valStr, 10, 64)
			if err == nil {
				rankList = append(rankList, critID)
			}
		}
		rankings[expertID] = rankList
	}

	allCriteria, err := h.critRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": "Could not fetch criteria"})
		return
	}

	var allIDs []int64
	for _, crit := range allCriteria {
		allIDs = append(allIDs, crit.ID)
	}

	strategy := h.votingSvc.GetStrategy(methodStr)
	newWeights := h.votingSvc.Execute(strategy, rankings, allIDs)

	// Шлях А: відсутні множаться на нуль (просто записуємо отримані ваги, неоцінені матимуть 0)
	for _, crit := range allCriteria {
		weight, exists := newWeights[crit.ID]
		if !exists {
			weight = 0
		}
		crit.Weight = weight
		err := h.critRepo.Update(crit)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": "Failed to update weights: " + err.Error()})
			return
		}
	}

	c.Redirect(http.StatusSeeOther, "/criteria?last_file="+url.QueryEscape(originalFilename)+"&last_method="+url.QueryEscape(methodStr))
}

func (h *ImportHandler) ImportEvaluations(c *gin.Context) {
	fileHeader, err := c.FormFile("import_file")
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"title": "Error", "message": "No file uploaded"})
		return
	}

	methodStr := c.PostForm("aggregation_method")
	if methodStr == "" {
		methodStr = string(analytics.WeightedMean)
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": "Could not open file"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil || len(records) < 2 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"title": "Error", "message": "Invalid CSV file"})
		return
	}

	// Parsing CSV: First row: "ExpertID", "AlternativeID", CritID1, CritID2...
	if len(records[0]) < 3 {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"title": "Error", "message": "Invalid CSV format (need at least Expert, Alt, and 1 Criterion)"})
		return
	}

	var critIDs []int64
	for _, header := range records[0][2:] {
		id, err := strconv.ParseInt(strings.TrimSpace(header), 10, 64)
		if err != nil {
			continue
		}
		critIDs = append(critIDs, id)
	}

	expertEvals := make(map[string]map[int64]map[int64]float64)

	for i, row := range records {
		if i == 0 {
			continue
		}
		if len(row) < len(critIDs)+2 {
			continue
		}

		expertID := strings.TrimSpace(row[0])
		altID, err := strconv.ParseInt(strings.TrimSpace(row[1]), 10, 64)
		if err != nil {
			continue
		}

		if expertEvals[expertID] == nil {
			expertEvals[expertID] = make(map[int64]map[int64]float64)
		}
		if expertEvals[expertID][altID] == nil {
			expertEvals[expertID][altID] = make(map[int64]float64)
		}

		for j, valStr := range row[2 : 2+len(critIDs)] {
			val, err := strconv.ParseFloat(strings.TrimSpace(valStr), 64)
			if err == nil {
				critID := critIDs[j]
				expertEvals[expertID][altID][critID] = val
			}
		}
	}

	aggregated := h.expertSvc.AggregateEvaluations(expertEvals, analytics.AggregationMethod(methodStr))

	var batch []*entities.Evaluation
	for altID, critMap := range aggregated {
		for critID, score := range critMap {
			batch = append(batch, &entities.Evaluation{
				AlternativeID: altID,
				CriterionID:   critID,
				Value:         score,
			})
		}
	}

	if err := h.evalRepo.UpsertBatch(batch); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}
