package analytics

import (
	"dss/internal/domain/entities"
	"log"
	"sort"
)

type AlternativeScore struct {
	Alternative       *entities.Alternative
	Score             float64
	RelativeScore     float64
	CriteriaScores    map[int64]float64
	CriteriaScoresExt []struct {
		CriterionName string
		Score         float64
	}
}

type ScoringService struct{}

func NewScoringService() *ScoringService {
	return &ScoringService{}
}

func (s *ScoringService) Rank(
	alternatives []*entities.Alternative,
	criteria []*entities.Criterion,
	evaluations map[int64]map[int64]float64,
	strategy ScoringStrategy,
) ([]*AlternativeScore, error) {
	intermediate := make([]*AlternativeScore, 0, len(alternatives))

	for _, alt := range alternatives {
		x, w := s.buildNormalizedArrays(alt.ID, criteria, evaluations)
		criteriaScores := s.buildCriteriaScores(criteria, x, w)

		criteriaScoresExt := make([]struct {
			CriterionName string
			Score         float64
		}, len(criteria))
		for i, crit := range criteria {
			criteriaScoresExt[i] = struct {
				CriterionName string
				Score         float64
			}{
				CriterionName: crit.Name,
				Score:         criteriaScores[crit.ID],
			}
		}

		score, err := strategy.Eval(x, w)
		if err != nil {
			log.Println(err)
			score = 0
		}

		intermediate = append(intermediate, &AlternativeScore{
			Alternative:       alt,
			Score:             score,
			CriteriaScores:    criteriaScores,
			CriteriaScoresExt: criteriaScoresExt,
		})
	}

	maxValue := 0.0
	for _, r := range intermediate {
		if r.Score > maxValue {
			maxValue = r.Score
		}
	}

	scores := make([]*AlternativeScore, 0, len(intermediate))
	for _, r := range intermediate {
		rel := 0.0
		if maxValue > 0 {
			rel = r.Score / maxValue
		} else {
			rel = 0
		}
		scores = append(scores, &AlternativeScore{
			Alternative:       r.Alternative,
			Score:             r.Score,
			CriteriaScores:    r.CriteriaScores,
			RelativeScore:     rel * 100,
			CriteriaScoresExt: r.CriteriaScoresExt,
		})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	for i, score := range scores {
		log.Printf("%d. %d: %.4f\n", i+1, score.Alternative.ID, score.Score)
	}

	return scores, nil
}

func (s *ScoringService) buildNormalizedArrays(altID int64, criteria []*entities.Criterion, evaluations map[int64]map[int64]float64) ([]float64, []float64) {
	x := make([]float64, len(criteria))
	w := make([]float64, len(criteria))
	var wsum float64
	for i, crit := range criteria {
		row := evaluations[altID]
		if row == nil {
			x[i] = 0
		} else {
			x[i] = crit.Normalize(row[crit.ID])
		}
		w[i] = crit.Weight
		if w[i] <= 0 {
			w[i] = 1.0
		}
		wsum += w[i]
	}
	for i, _ := range w {
		w[i] /= wsum
	}
	return x, w
}

func (s *ScoringService) buildCriteriaScores(criteria []*entities.Criterion, x, w []float64) map[int64]float64 {
	scores := make(map[int64]float64)
	for i, crit := range criteria {
		scores[crit.ID] = x[i] * w[i]
	}
	return scores
}

func (s *ScoringService) GetStrategy(name string) ScoringStrategy {
	switch name {
	case "multiplicative":
		return &MultiplicativeStrategy{}
	case "minimax":
		return &CautiousStrategy{}
	default:
		return &AdditiveStrategy{}
	}
}
