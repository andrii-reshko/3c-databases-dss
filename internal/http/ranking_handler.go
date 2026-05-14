package http

import (
	"dss/internal/domain/analytics"
	"dss/internal/domain/entities"
	"dss/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RankingHandler struct {
	altRepo  repositories.AlternativeRepository
	critRepo repositories.CriterionRepository
	evalRepo repositories.EvaluationRepository
	ruleRepo repositories.RuleRepository
}

func NewRankingHandler(altRepo repositories.AlternativeRepository, critRepo repositories.CriterionRepository, evalRepo repositories.EvaluationRepository, ruleRepo repositories.RuleRepository) *RankingHandler {
	return &RankingHandler{
		altRepo:  altRepo,
		critRepo: critRepo,
		evalRepo: evalRepo,
		ruleRepo: ruleRepo,
	}
}

type RankingView struct {
	Alternatives   []*entities.Alternative
	Criteria       []*entities.Criterion
	Scores         []*analytics.AlternativeScore
	StrategyName   string
	Strategy       string
	Charts         map[int64]float64
	HasEvaluations bool
}

func (h *RankingHandler) ShowRanking(c *gin.Context) {
	strategyName := c.DefaultQuery("strategy", "additive")

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

	// Apply custom weights for sensitivity analysis
	for _, crit := range criteria {
		weightStr := c.Query("weight_" + strconv.FormatInt(crit.ID, 10))
		if weightStr != "" {
			if w, err := strconv.ParseFloat(weightStr, 64); err == nil {
				crit.Weight = w
			}
		}
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

	scoringService := analytics.NewScoringService()
	strategy := scoringService.GetStrategy(strategyName)

	rules, err := h.ruleRepo.FindAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}

	scores, err := scoringService.Rank(alternatives, criteria, evalMap, strategy, rules)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}

	view := RankingView{
		Alternatives:   alternatives,
		Criteria:       criteria,
		Scores:         scores,
		StrategyName:   strategy.GetName(),
		Strategy:       strategyName,
		HasEvaluations: len(evaluations) > 0,
	}

	c.HTML(http.StatusOK, "ranking.html", gin.H{"title": "Ranking", "view": view})
}
