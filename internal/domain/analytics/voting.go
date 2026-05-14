package analytics

import (
	"math"
)

type VotingStrategy interface {
	GetName() string
	CalculateWeights(rankings map[string][]int64, criteriaIDs []int64) map[int64]float64
}

// --- Helper Functions ---

func makeInitScores(criteriaIDs []int64) map[int64]float64 {
	scores := make(map[int64]float64)
	for _, id := range criteriaIDs {
		scores[id] = 0
	}
	return scores
}

func findPosition(ranking []int64, item int64) int {
	for i, v := range ranking {
		if v == item {
			return i
		}
	}
	return len(ranking) // Assume lowest priority if not ranked
}

func normalizeWeights(scores map[int64]float64) map[int64]float64 {
	sum := 0.0
	for _, score := range scores {
		sum += score
	}

	weights := make(map[int64]float64)
	if sum == 0 {
		// Distribute equally if all scores are 0
		n := len(scores)
		if n == 0 {
			return weights
		}
		eq := 1.0 / float64(n)
		for id := range scores {
			weights[id] = eq
		}
		return weights
	}

	for id, score := range scores {
		weights[id] = score / sum
	}

	return weights
}

// --- Strategies ---

type BordaStrategy struct{}

func (s *BordaStrategy) GetName() string { return "Borda" }
func (s *BordaStrategy) CalculateWeights(rankings map[string][]int64, criteriaIDs []int64) map[int64]float64 {
	scores := makeInitScores(criteriaIDs)
	m := len(criteriaIDs)
	for _, ranking := range rankings {
		for i, critID := range ranking {
			// score is (m - 1 - position)
			scores[critID] += float64(m - 1 - i)
		}
	}
	return normalizeWeights(scores)
}

type CopelandStrategy struct{}

func (s *CopelandStrategy) GetName() string { return "Copeland" }
func (s *CopelandStrategy) CalculateWeights(rankings map[string][]int64, criteriaIDs []int64) map[int64]float64 {
	scores := makeInitScores(criteriaIDs)
	
	// For each pair (A, B), count how many experts prefer A over B
	for i := 0; i < len(criteriaIDs); i++ {
		for j := i + 1; j < len(criteriaIDs); j++ {
			a := criteriaIDs[i]
			b := criteriaIDs[j]

			aWins := 0
			bWins := 0

			for _, ranking := range rankings {
				posA := findPosition(ranking, a)
				posB := findPosition(ranking, b)
				if posA < posB {
					aWins++
				} else if posB < posA {
					bWins++
				}
			}

			if aWins > bWins {
				scores[a] += 1
				scores[b] -= 1
			} else if bWins > aWins {
				scores[b] += 1
				scores[a] -= 1
			}
		}
	}

	// Shift scores to be positive to avoid negative weights
	minScore := math.MaxFloat64
	for _, score := range scores {
		if score < minScore {
			minScore = score
		}
	}

	// Shift by |minScore| + 1 so the lowest score is 1
	shift := 0.0
	if minScore <= 0 {
		shift = math.Abs(minScore) + 1
	}

	for id := range scores {
		scores[id] += shift
	}

	return normalizeWeights(scores)
}

type SimpsonStrategy struct{}

func (s *SimpsonStrategy) GetName() string { return "Simpson" }
func (s *SimpsonStrategy) CalculateWeights(rankings map[string][]int64, criteriaIDs []int64) map[int64]float64 {
	scores := makeInitScores(criteriaIDs)
	
	for i, a := range criteriaIDs {
		minWins := math.MaxFloat64
		for j, b := range criteriaIDs {
			if i == j {
				continue
			}
			aWins := 0
			for _, ranking := range rankings {
				posA := findPosition(ranking, a)
				posB := findPosition(ranking, b)
				if posA < posB {
					aWins++
				}
			}
			if float64(aWins) < minWins {
				minWins = float64(aWins)
			}
		}
		if len(criteriaIDs) > 1 {
			scores[a] = minWins
		} else {
			scores[a] = 1 // Fallback
		}
	}
	
	return normalizeWeights(scores)
}

type RelativeMajorityStrategy struct{}

func (s *RelativeMajorityStrategy) GetName() string { return "Relative Majority" }
func (s *RelativeMajorityStrategy) CalculateWeights(rankings map[string][]int64, criteriaIDs []int64) map[int64]float64 {
	scores := makeInitScores(criteriaIDs)
	
	for _, ranking := range rankings {
		if len(ranking) > 0 {
			firstChoice := ranking[0]
			scores[firstChoice] += 1
		}
	}
	
	return normalizeWeights(scores)
}

// --- Voting Service Context ---

type VotingService struct{}

func NewVotingService() *VotingService {
	return &VotingService{}
}

// Execute applies the given voting strategy
func (s *VotingService) Execute(strategy VotingStrategy, rankings map[string][]int64, criteriaIDs []int64) map[int64]float64 {
	if strategy == nil {
		strategy = &BordaStrategy{} // Default
	}
	return strategy.CalculateWeights(rankings, criteriaIDs)
}

func (s *VotingService) GetStrategy(name string) VotingStrategy {
	switch name {
	case "copeland":
		return &CopelandStrategy{}
	case "simpson":
		return &SimpsonStrategy{}
	case "relative_majority":
		return &RelativeMajorityStrategy{}
	case "borda":
		fallthrough
	default:
		return &BordaStrategy{}
	}
}
