package analytics

import (
	"errors"
	"math"
)

type ScoringStrategy interface {
	GetName() string
	Eval(x, w []float64) (float64, error)
}

type AdditiveStrategy struct{}

func (s *AdditiveStrategy) GetName() string { return "Additive" }
func (s *AdditiveStrategy) Eval(x, w []float64) (float64, error) {
	if len(x) != len(w) {
		return 0, errors.New("length mismatch between x and w")
	}

	result := 0.0
	for i := range x {
		result += w[i] * x[i]
	}

	return result, nil
}

type MultiplicativeStrategy struct{}

func (s *MultiplicativeStrategy) GetName() string { return "Multiplicative" }
func (s *MultiplicativeStrategy) Eval(x, w []float64) (float64, error) {
	if len(x) != len(w) {
		return 0, errors.New("length mismatch between x and w")
	}

	result := 1.0
	for i := range x {
		if x[i] <= 0 {
			return 0, nil
		}
		result *= math.Pow(x[i], w[i])
	}

	return result, nil
}

type CautiousStrategy struct{}

func (s *CautiousStrategy) GetName() string { return "Cautious" }
func (s *CautiousStrategy) Eval(x, w []float64) (float64, error) {
	if len(x) != len(w) {
		return 0, errors.New("length mismatch between x and w")
	}

	if len(x) == 0 {
		return 0, errors.New("empty input")
	}

	minVal := w[0] * x[0]

	for i := 1; i < len(x); i++ {
		val := w[i] * x[i]
		if val < minVal {
			minVal = val
		}
	}

	return minVal, nil
}
