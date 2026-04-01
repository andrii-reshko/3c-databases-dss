package http

import (
	"dss/internal/domain/entities"
	"dss/internal/repositories"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EvaluationHandler struct {
	altRepo  repositories.AlternativeRepository
	critRepo repositories.CriterionRepository
	evalRepo repositories.EvaluationRepository
}

func NewEvaluationHandler(altRepo repositories.AlternativeRepository, critRepo repositories.CriterionRepository, evalRepo repositories.EvaluationRepository) *EvaluationHandler {
	return &EvaluationHandler{
		altRepo:  altRepo,
		critRepo: critRepo,
		evalRepo: evalRepo,
	}
}

type MatrixView struct {
	Alternatives []*entities.Alternative
	Criteria     []*entities.Criterion
	Values       map[int64]map[int64]float64
}

func (h *EvaluationHandler) ShowMatrix(c *gin.Context) {
	alternatives, err := h.altRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	criteria, err := h.critRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	evaluations, err := h.evalRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}

	evalMap := make(map[int64]map[int64]float64)
	for _, alt := range alternatives {
		evalMap[alt.ID] = make(map[int64]float64)
	}
	for _, e := range evaluations {
		if row, ok := evalMap[e.AlternativeID]; ok {
			row[e.CriterionID] = e.Value
		}
	}

	view := MatrixView{
		Alternatives: alternatives,
		Criteria:     criteria,
		Values:       evalMap,
	}

	c.HTML(http.StatusOK, "index.html", gin.H{"title": "Evaluation Matrix", "view": view})
}

func (h *EvaluationHandler) UpdateMatrix(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"title": "Error", "message": "Invalid form data"})
		return
	}

	var evaluations []*entities.Evaluation
	for key, values := range c.Request.PostForm {
		if len(values) == 0 {
			continue
		}

		var altID, critID int64
		if _, err := fmt.Sscanf(key, "eval-%d-%d", &altID, &critID); err != nil {
			continue
		}

		value, err := strconv.ParseFloat(values[0], 64)
		if err != nil {
			value = 0
		}

		evaluations = append(evaluations, &entities.Evaluation{
			AlternativeID: altID,
			CriterionID:   critID,
			Value:         value,
		})
	}

	if err := h.evalRepo.UpsertBatch(evaluations); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/ranking")
}
