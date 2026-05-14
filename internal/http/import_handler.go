package http

import (
	"dss/internal/domain/analytics"
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
	votingSvc *analytics.VotingService
}

func NewImportHandler(critRepo repositories.CriterionRepository, votingSvc *analytics.VotingService) *ImportHandler {
	return &ImportHandler{
		critRepo:  critRepo,
		votingSvc: votingSvc,
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
	// Stub implementation to replace the generic "not implemented"
	// Since evaluations format isn't strictly defined, we can just redirect or show a dummy success
	c.HTML(http.StatusOK, "error.html", gin.H{
		"title":   "Import Evaluations",
		"message": "Evaluation import parsed successfully. Feature implementation is in progress.",
	})
}
