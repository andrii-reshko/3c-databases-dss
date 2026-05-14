package http

import (
	"dss/internal/domain/entities"
	"dss/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RuleHandler struct {
	ruleRepo repositories.RuleRepository
	critRepo repositories.CriterionRepository
}

func NewRuleHandler(ruleRepo repositories.RuleRepository, critRepo repositories.CriterionRepository) *RuleHandler {
	return &RuleHandler{
		ruleRepo: ruleRepo,
		critRepo: critRepo,
	}
}

type RuleView struct {
	Rules   []*entities.Rule
	CritMap map[int64]*entities.Criterion
}

func (h *RuleHandler) ShowRules(c *gin.Context) {
	criteria, err := h.critRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}

	critMap := make(map[int64]*entities.Criterion)
	for _, cr := range criteria {
		critMap[cr.ID] = cr
	}

	rules, err := h.ruleRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}

	view := RuleView{
		Rules:   rules,
		CritMap: critMap,
	}

	c.HTML(http.StatusOK, "rules.html", gin.H{"title": "Expert Logic", "view": view})
}

func (h *RuleHandler) ShowRuleForm(c *gin.Context) {
	idStr := c.Param("id")
	criteria, err := h.critRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": "Could not load criteria"})
		return
	}

	if idStr == "" {
		// New Rule form
		c.HTML(http.StatusOK, "rule_form.html", gin.H{
			"title":    "Add New Rule",
			"action":   "/rules/new",
			"criteria": criteria,
			"rule": &entities.Rule{
				Operator:   ">",
				ActionType: "modify",
			},
		})
		return
	}

	// Edit Rule form
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"title": "Error", "message": "Invalid ID"})
		return
	}

	rule, err := h.ruleRepo.FindByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"title": "Error", "message": "Rule not found"})
		return
	}

	c.HTML(http.StatusOK, "rule_form.html", gin.H{
		"title":    "Edit Rule",
		"action":   "/rules/edit/" + idStr,
		"criteria": criteria,
		"rule":     rule,
	})
}

func (h *RuleHandler) CreateRule(c *gin.Context) {
	critID, _ := strconv.ParseInt(c.PostForm("criterion_id"), 10, 64)
	val, _ := strconv.ParseFloat(c.PostForm("value"), 64)
	actionVal, _ := strconv.ParseFloat(c.PostForm("action_value"), 64)

	rule := &entities.Rule{
		Name:        c.PostForm("name"),
		CriterionID: critID,
		Operator:    c.PostForm("operator"),
		Value:       val,
		ActionType:  entities.ActionType(c.PostForm("action_type")),
		ActionValue: actionVal,
	}

	if _, err := h.ruleRepo.Create(rule); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/rules")
}

func (h *RuleHandler) UpdateRule(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	critID, _ := strconv.ParseInt(c.PostForm("criterion_id"), 10, 64)
	val, _ := strconv.ParseFloat(c.PostForm("value"), 64)
	actionVal, _ := strconv.ParseFloat(c.PostForm("action_value"), 64)

	rule := &entities.Rule{
		ID:          id,
		Name:        c.PostForm("name"),
		CriterionID: critID,
		Operator:    c.PostForm("operator"),
		Value:       val,
		ActionType:  entities.ActionType(c.PostForm("action_type")),
		ActionValue: actionVal,
	}

	if err := h.ruleRepo.Update(rule); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/rules")
}

func (h *RuleHandler) DeleteRule(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.ruleRepo.Delete(id); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, "/rules")
}
