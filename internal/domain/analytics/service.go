package analytics

import "dss/internal/domain/entities"

// AnalyticsService TODO implement later with actual decision-making logic
type AnalyticsService interface {
	RankAlternatives(
		alternatives []*entities.Alternative,
		criteria []*entities.Criterion,
		evaluations []*entities.Evaluation,
	) ([]*entities.Alternative, error)
}
