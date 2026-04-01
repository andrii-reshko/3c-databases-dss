package http

import (
	"dss/internal/domain/analytics"
	"dss/internal/domain/entities"
	"dss/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RankingHandler struct {
	altRepo  repositories.AlternativeRepository
	critRepo repositories.CriterionRepository
	evalRepo repositories.EvaluationRepository
}

func NewRankingHandler(altRepo repositories.AlternativeRepository, critRepo repositories.CriterionRepository, evalRepo repositories.EvaluationRepository) *RankingHandler {
	return &RankingHandler{
		altRepo:  altRepo,
		critRepo: critRepo,
		evalRepo: evalRepo,
	}
}

type RankingView struct {
	Alternatives []*entities.Alternative
	Criteria     []*entities.Criterion
	Scores       []*analytics.AlternativeScore
	StrategyName string
	Strategy     string
	Charts       map[int64]float64
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
	scores, err := scoringService.Rank(alternatives, criteria, evalMap, strategy)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "message": err.Error()})
		return
	}

	view := RankingView{
		Alternatives: alternatives,
		Criteria:     criteria,
		Scores:       scores,
		StrategyName: strategy.GetName(),
		Strategy:     strategyName,
	}

	c.HTML(http.StatusOK, "ranking.html", gin.H{"title": "Ranking", "view": view})
}
