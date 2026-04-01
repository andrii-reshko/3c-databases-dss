package http

import (
	"dss/internal/domain/entities"
	"dss/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlternativeHandler struct {
	repo repositories.AlternativeRepository
}

func NewAlternativeHandler(repo repositories.AlternativeRepository) *AlternativeHandler {
	return &AlternativeHandler{repo: repo}
}

func (h *AlternativeHandler) ListAlternatives(c *gin.Context) {
	alternatives, err := h.repo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "alternatives.html", gin.H{
		"title":        "Alternatives",
		"alternatives": alternatives,
	})
}

func (h *AlternativeHandler) ShowAlternativeForm(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" { // New alternative form
		c.HTML(http.StatusOK, "alternative_form.html", gin.H{
			"title":  "Add New Alternative",
			"action": "/alternatives/new",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"title": "Error", "message": "Invalid ID"})
		return
	}

	alt, err := h.repo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"title": "Error", "message": "Alternative not found"})
		return
	}

	c.HTML(http.StatusOK, "alternative_form.html", gin.H{
		"title":       "Edit Alternative",
		"action":      "/alternatives/edit/" + idStr,
		"alternative": alt,
	})
}

func (h *AlternativeHandler) CreateAlternative(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")

	var alt entities.Alternative
	alt.Name = name
	alt.Description = description

	if _, err := h.repo.Create(&alt); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/")
}

func (h *AlternativeHandler) UpdateAlternative(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var alt entities.Alternative
	alt.ID = id
	alt.Name = c.PostForm("name")
	alt.Description = c.PostForm("description")

	if err := h.repo.Update(&alt); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/alternatives/")
}

func (h *AlternativeHandler) DeleteAlternative(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.repo.Delete(id); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/alternatives/")
}
