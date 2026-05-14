package analytics

import (
	"log"
	"math"
	"sort"
)

type AggregationMethod string

const (
	WeightedMean AggregationMethod = "weighted_mean"
	SumOfRanks   AggregationMethod = "sum_of_ranks"
	KemenySnell  AggregationMethod = "kemeny_snell"
)

type ExpertEvaluationService struct{}

func NewExpertEvaluationService() *ExpertEvaluationService {
	return &ExpertEvaluationService{}
}

// AggregateEvaluations aggregates raw evaluations from multiple experts into a single matrix.
// expertEvals: map of expertID -> alternativeID -> criterionID -> score
func (s *ExpertEvaluationService) AggregateEvaluations(expertEvals map[string]map[int64]map[int64]float64, method AggregationMethod) map[int64]map[int64]float64 {
	
	// Collect all criteria IDs and Alternative IDs
	critSet := make(map[int64]bool)
	altSet := make(map[int64]bool)
	for _, altMap := range expertEvals {
		for altID, critMap := range altMap {
			altSet[altID] = true
			for critID := range critMap {
				critSet[critID] = true
			}
		}
	}

	var altIDs []int64
	for id := range altSet {
		altIDs = append(altIDs, id)
	}
	sort.Slice(altIDs, func(i, j int) bool { return altIDs[i] < altIDs[j] })

	var critIDs []int64
	for id := range critSet {
		critIDs = append(critIDs, id)
	}

	result := make(map[int64]map[int64]float64)
	for _, altID := range altIDs {
		result[altID] = make(map[int64]float64)
	}

	for _, critID := range critIDs {
		switch method {
		case WeightedMean:
			s.applyWeightedMean(expertEvals, altIDs, critID, result)
		case SumOfRanks:
			s.applySumOfRanks(expertEvals, altIDs, critID, result)
		case KemenySnell:
			s.applyKemenySnell(expertEvals, altIDs, critID, result)
		default:
			s.applyWeightedMean(expertEvals, altIDs, critID, result)
		}
	}

	return result
}

func (s *ExpertEvaluationService) applyWeightedMean(expertEvals map[string]map[int64]map[int64]float64, altIDs []int64, critID int64, result map[int64]map[int64]float64) {
	// Default to equal weights if no specific weights are provided
	weight := 1.0 / float64(len(expertEvals))
	if len(expertEvals) == 0 {
		weight = 0
	}

	for _, altID := range altIDs {
		sum := 0.0
		for _, altMap := range expertEvals {
			if critMap, ok := altMap[altID]; ok {
				sum += critMap[critID] * weight
			}
		}
		result[altID][critID] = sum
	}
}

func (s *ExpertEvaluationService) applySumOfRanks(expertEvals map[string]map[int64]map[int64]float64, altIDs []int64, critID int64, result map[int64]map[int64]float64) {
	m := float64(len(expertEvals))
	n := float64(len(altIDs))

	rankSums := make(map[int64]float64)
	for _, altID := range altIDs {
		sum := 0.0
		for _, altMap := range expertEvals {
			if critMap, ok := altMap[altID]; ok {
				sum += critMap[critID]
			}
		}
		rankSums[altID] = sum
		result[altID][critID] = sum
	}

	// Calculate Concordance Coefficient W
	// W = 12 * S / (m^2 * (n^3 - n))
	// where S is the sum of squared deviations of rank sums from the mean rank sum
	if m > 0 && n > 1 {
		meanSum := m * (n + 1) / 2.0
		S := 0.0
		for _, altID := range altIDs {
			dev := rankSums[altID] - meanSum
			S += dev * dev
		}
		W := (12.0 * S) / (m * m * (n*n*n - n))
		log.Printf("[SumOfRanks] Criterion %d: Concordance Coefficient W = %.4f\n", critID, W)
	}
}

func (s *ExpertEvaluationService) applyKemenySnell(expertEvals map[string]map[int64]map[int64]float64, altIDs []int64, critID int64, result map[int64]map[int64]float64) {
	// Kemeny-Snell median ranking
	// For each pair of alternatives (A, B), an expert gives:
	// 1 if A is preferred to B (rankA < rankB)
	// -1 if B is preferred to A (rankA > rankB)
	// 0 if tied

	n := len(altIDs)
	if n == 0 {
		return
	}

	expertsCount := len(expertEvals)
	if expertsCount == 0 {
		return
	}

	// Precompute expert pair matrices
	expertMatrices := make([][][]int, 0, expertsCount)
	for _, altMap := range expertEvals {
		matrix := make([][]int, n)
		for i := 0; i < n; i++ {
			matrix[i] = make([]int, n)
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				valI := altMap[altIDs[i]][critID]
				valJ := altMap[altIDs[j]][critID]
				// Assuming lower value means better rank
				if valI < valJ {
					matrix[i][j] = 1
				} else if valI > valJ {
					matrix[i][j] = -1
				} else {
					matrix[i][j] = 0
				}
			}
		}
		expertMatrices = append(expertMatrices, matrix)
	}

	// Generate all permutations of ranks (1 to n)
	var bestPerm []int
	minDistance := math.MaxFloat64

	var generatePermutations func([]int, int)
	generatePermutations = func(arr []int, k int) {
		if k == len(arr) {
			// Calculate distance for this permutation
			dist := 0.0
			for _, expMat := range expertMatrices {
				for i := 0; i < n; i++ {
					for j := 0; j < n; j++ {
						if i == j {
							continue
						}
						// arr contains the rank for alternative i
						// lower rank is better
						var permVal int
						if arr[i] < arr[j] {
							permVal = 1
						} else if arr[i] > arr[j] {
							permVal = -1
						} else {
							permVal = 0
						}

						diff := float64(permVal - expMat[i][j])
						dist += math.Abs(diff)
					}
				}
			}
			
			if dist < minDistance {
				minDistance = dist
				bestPerm = make([]int, n)
				copy(bestPerm, arr)
			}
			return
		}

		for i := k; i < len(arr); i++ {
			arr[k], arr[i] = arr[i], arr[k]
			generatePermutations(arr, k+1)
			arr[k], arr[i] = arr[i], arr[k]
		}
	}

	initialPerm := make([]int, n)
	for i := 0; i < n; i++ {
		initialPerm[i] = i + 1
	}

	// Limit n to prevent excessive permutations (e.g. n <= 8 -> 40k perms)
	if n <= 8 {
		generatePermutations(initialPerm, 0)
	} else {
		// Fallback to simple sum of ranks if too many alternatives
		log.Printf("[KemenySnell] Too many alternatives (%d), falling back to SumOfRanks\n", n)
		s.applySumOfRanks(expertEvals, altIDs, critID, result)
		return
	}

	// Apply the best permutation as the result ranks
	for i, altID := range altIDs {
		result[altID][critID] = float64(bestPerm[i])
	}
}
