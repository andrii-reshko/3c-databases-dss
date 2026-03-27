package http

import (
	"dss/internal/domain/entities"
	"dss/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CriterionHandler struct {
	repo repositories.CriterionRepository
}

func NewCriterionHandler(repo repositories.CriterionRepository) *CriterionHandler {
	return &CriterionHandler{repo: repo}
}

func (h *CriterionHandler) ListCriteria(c *gin.Context) {
	criteria, err := h.repo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "criteria.html", gin.H{
		"title":    "Criteria",
		"criteria": criteria,
	})
}

func (h *CriterionHandler) ShowCriterionForm(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.HTML(http.StatusOK, "criterion_form.html", gin.H{
			"title":  "Add New Criterion",
			"action": "/criteria/new",
			"criterion": &entities.Criterion{
				InputMin:  0,
				InputMax:  1,
				OutputMin: 0,
				OutputMax: 1,
			},
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"title": "Error", "message": "Invalid ID"})
		return
	}

	criterion, err := h.repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"title": "Error", "message": "Criterion not found"})
		return
	}

	c.HTML(http.StatusOK, "criterion_form.html", gin.H{
		"title":     "Edit Criterion",
		"action":    "/criteria/edit/" + idStr,
		"criterion": criterion,
	})
}

func (h *CriterionHandler) CreateCriterion(c *gin.Context) {
	var criterion entities.Criterion
	criterion.Name = c.PostForm("name")
	criterion.Type = entities.CriterionType(c.PostForm("type"))
	criterion.Description = c.PostForm("description")
	criterion.InputMin, _ = strconv.ParseFloat(c.PostForm("input_min"), 64)
	criterion.InputMax, _ = strconv.ParseFloat(c.PostForm("input_max"), 64)
	criterion.OutputMin, _ = strconv.ParseFloat(c.PostForm("output_min"), 64)
	criterion.OutputMax, _ = strconv.ParseFloat(c.PostForm("output_max"), 64)

	if _, err := h.repo.Create(&criterion); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/criteria/")
}

func (h *CriterionHandler) UpdateCriterion(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var criterion entities.Criterion
	criterion.ID = id
	criterion.Name = c.PostForm("name")
	criterion.Type = entities.CriterionType(c.PostForm("type"))
	criterion.Description = c.PostForm("description")
	criterion.InputMin, _ = strconv.ParseFloat(c.PostForm("input_min"), 64)
	criterion.InputMax, _ = strconv.ParseFloat(c.PostForm("input_max"), 64)
	criterion.OutputMin, _ = strconv.ParseFloat(c.PostForm("output_min"), 64)
	criterion.OutputMax, _ = strconv.ParseFloat(c.PostForm("output_max"), 64)

	if err := h.repo.Update(&criterion); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/criteria/")
}

func (h *CriterionHandler) DeleteCriterion(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.Delete(id); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/criteria/")
}
