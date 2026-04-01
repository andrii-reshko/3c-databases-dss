package analytics

import (
	"math"
	"testing"
)

func TestAdditiveStrategy(t *testing.T) {
	strategy := &AdditiveStrategy{}

	tests := []struct {
		name     string
		x        []float64
		w        []float64
		expected float64
	}{
		{
			name:     "basic",
			x:        []float64{0.5, 0.8, 0.3},
			w:        []float64{0.2, 0.5, 0.3},
			expected: 0.2*0.5 + 0.5*0.8 + 0.3*0.3,
		},
		{
			name:     "all ones",
			x:        []float64{1, 1, 1},
			w:        []float64{0.5, 0.3, 0.2},
			expected: 1.0,
		},
		{
			name:     "all zeros",
			x:        []float64{0, 0, 0},
			w:        []float64{0.5, 0.3, 0.2},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := strategy.Eval(tt.x, tt.w)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMultiplicativeStrategy(t *testing.T) {
	strategy := &MultiplicativeStrategy{}

	tests := []struct {
		name     string
		x        []float64
		w        []float64
		expected float64
	}{
		{
			name:     "all ones",
			x:        []float64{1, 1, 1},
			w:        []float64{0.5, 0.3, 0.2},
			expected: 1.0,
		},
		{
			name:     "basic",
			x:        []float64{0.5, 0.8, 0.3},
			w:        []float64{0.2, 0.5, 0.3},
			expected: math.Pow(0.5, 0.2) * math.Pow(0.8, 0.5) * math.Pow(0.3, 0.3),
		},
		{
			name:     "at least one zero",
			x:        []float64{0.5, 0.0, 0.3},
			w:        []float64{0.2, 0.5, 0.3},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := strategy.Eval(tt.x, tt.w)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !approximatelyEqual(result, tt.expected, 0.0001) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCautiousStrategy(t *testing.T) {
	strategy := &CautiousStrategy{}

	tests := []struct {
		name     string
		x        []float64
		w        []float64
		expected float64
	}{
		{
			name:     "basic",
			x:        []float64{0.5, 0.8, 0.3},
			w:        []float64{0.2, 0.5, 0.3},
			expected: min(0.5*0.2, 0.8*0.5, 0.3*0.3),
		},
		{
			name:     "first is minimum",
			x:        []float64{0.1, 0.8, 0.9},
			w:        []float64{1.0, 1.0, 1.0},
			expected: 0.1,
		},
		{
			name:     "last is minimum",
			x:        []float64{0.9, 0.8, 0.1},
			w:        []float64{1.0, 1.0, 1.0},
			expected: 0.1,
		},
		{
			name:     "middle is minimum",
			x:        []float64{0.9, 0.1, 0.8},
			w:        []float64{1.0, 1.0, 1.0},
			expected: 0.1,
		},
		{
			name:     "all equal",
			x:        []float64{0.5, 0.5, 0.5},
			w:        []float64{1.0, 1.0, 1.0},
			expected: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := strategy.Eval(tt.x, tt.w)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func min(values ...float64) float64 {
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

func approximatelyEqual(a, b, epsilon float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}
