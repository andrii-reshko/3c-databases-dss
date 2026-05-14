package entities

import (
	"math"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		c        *Criterion
		input    float64
		expected float64
	}{
		{
			name: "Pub time sub 1h",
			c: &Criterion{
				InputMin:  1,
				InputMax:  24,
				OutputMin: 1,
				OutputMax: 0,
			},
			input:    0.5, // sub 1h is ok
			expected: 1,
		},
		{
			name: "Pub time 12h",
			c: &Criterion{
				InputMin:  1,
				InputMax:  24,
				OutputMin: 1,
				OutputMax: 0,
			},
			input:    12,
			expected: 0.521, // 1 - norm(12) = 1 - (12-1)/(24-1) = 1 - 11/23 = 12/23 = ~0.5217
		},
		{
			name: "Pub time 24h plus",
			c: &Criterion{
				InputMin:  1,
				InputMax:  24,
				OutputMin: 1,
				OutputMax: 0,
			},
			input:    36, // 24h plus is bad
			expected: 0,
		},
		{
			name: "Interviewing 0",
			c: &Criterion{
				InputMin:  5,
				InputMax:  5,
				OutputMin: 1,
				OutputMax: 0,
			},
			input:    0,
			expected: 1, // no active interviews - ok
		},
		{
			name: "Interviewing 4",
			c: &Criterion{
				InputMin:  5,
				InputMax:  5,
				OutputMin: 1,
				OutputMax: 0,
			},
			input:    4, // 4 is still ok
			expected: 1,
		},
		{
			name: "Interviewing 5",
			c: &Criterion{
				InputMin:  5,
				InputMax:  5,
				OutputMin: 1,
				OutputMax: 0,
			},
			input:    5, // 5 active interviews - immediately bad
			expected: 0,
		},
		{
			name: "Duration - short term",
			c: &Criterion{
				InputMin:  1,
				InputMax:  6,
				OutputMin: 0.5,
				OutputMax: 1,
			},
			input:    0.1, // sub 1 month is so-so
			expected: 0.5,
		},
		{
			name: "Duration - very long term",
			c: &Criterion{
				InputMin:  1,
				InputMax:  6,
				OutputMin: 0.5,
				OutputMax: 1,
			},
			input:    12, // 6 months plus is good
			expected: 1,
		},
		{
			name: "Duration - semi long term",
			c: &Criterion{
				InputMin:  1,
				InputMax:  6,
				OutputMin: 0.5,
				OutputMax: 1,
			},
			input:    3.5,
			expected: 0.75,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.c.Normalize(tc.input)
			if math.Abs(result-tc.expected) > 0.001 {
				t.Errorf("Normalize(%v) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}
