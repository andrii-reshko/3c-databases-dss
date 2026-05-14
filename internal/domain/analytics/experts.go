package analytics

import (
	"math"
	"sort"
)

type AggregationMethod string

const (
	ArithmeticMean AggregationMethod = "arithmetic_mean"
	Median         AggregationMethod = "median"
	GeometricMean  AggregationMethod = "geometric_mean"
)

type ExpertEvaluationService struct{}

func NewExpertEvaluationService() *ExpertEvaluationService {
	return &ExpertEvaluationService{}
}

// AggregateEvaluations aggregates raw evaluations from multiple experts into a single matrix
// expertEvals: map of expertID -> alternativeID -> criterionID -> score
func (s *ExpertEvaluationService) AggregateEvaluations(expertEvals map[string]map[int64]map[int64]float64, method AggregationMethod) map[int64]map[int64]float64 {
	// First, restructure data to collect all expert scores for each (alt, crit) pair
	collected := make(map[int64]map[int64][]float64)

	for _, altMatrix := range expertEvals {
		for altID, critMap := range altMatrix {
			if _, ok := collected[altID]; !ok {
				collected[altID] = make(map[int64][]float64)
			}
			for critID, score := range critMap {
				collected[altID][critID] = append(collected[altID][critID], score)
			}
		}
	}

	result := make(map[int64]map[int64]float64)

	for altID, critMap := range collected {
		result[altID] = make(map[int64]float64)
		for critID, scores := range critMap {
			var finalScore float64
			switch method {
			case ArithmeticMean:
				finalScore = s.calculateArithmeticMean(scores)
			case Median:
				finalScore = s.calculateMedian(scores)
			case GeometricMean:
				finalScore = s.calculateGeometricMean(scores)
			default:
				finalScore = s.calculateArithmeticMean(scores)
			}
			result[altID][critID] = finalScore
		}
	}

	return result
}

func (s *ExpertEvaluationService) calculateArithmeticMean(scores []float64) float64 {
	if len(scores) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range scores {
		sum += v
	}
	return sum / float64(len(scores))
}

func (s *ExpertEvaluationService) calculateMedian(scores []float64) float64 {
	if len(scores) == 0 {
		return 0
	}
	
	// Create a copy to avoid sorting the original slice
	sorted := make([]float64, len(scores))
	copy(sorted, scores)
	sort.Float64s(sorted)

	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2.0
	}
	return sorted[n/2]
}

func (s *ExpertEvaluationService) calculateGeometricMean(scores []float64) float64 {
	if len(scores) == 0 {
		return 0
	}
	
	product := 1.0
	count := 0.0
	for _, v := range scores {
		if v > 0 { // Geometric mean requires positive numbers
			product *= v
			count++
		}
	}
	
	if count == 0 {
		return 0
	}
	
	return math.Pow(product, 1.0/count)
}
