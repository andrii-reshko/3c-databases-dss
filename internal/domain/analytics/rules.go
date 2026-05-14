package analytics

import "dss/internal/domain/entities"

type RuleContext struct {
	Excluded     map[int64]bool
	Multipliers  map[int64]float64
	AppliedRules map[int64][]string // AltID -> list of applied rule names
}

type RuleEngine struct{}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

// Evaluate applies all threshold and logical rules to the raw evaluation matrix.
func (e *RuleEngine) Evaluate(
	alternatives []*entities.Alternative,
	evaluations map[int64]map[int64]float64,
	rules []*entities.Rule,
) *RuleContext {
	ctx := &RuleContext{
		Excluded:     make(map[int64]bool),
		Multipliers:  make(map[int64]float64),
		AppliedRules: make(map[int64][]string),
	}

	for _, alt := range alternatives {
		ctx.Multipliers[alt.ID] = 1.0 // Default multiplier is 1 (no change)

		rawVals := evaluations[alt.ID]
		if rawVals == nil {
			continue
		}

		for _, rule := range rules {
			val, exists := rawVals[rule.CriterionID]
			if !exists {
				continue
			}

			if rule.Matches(val) {
				ctx.AppliedRules[alt.ID] = append(ctx.AppliedRules[alt.ID], rule.Name)

				if rule.ActionType == entities.ActionExclude {
					// Rule triggered a cutoff threshold
					ctx.Excluded[alt.ID] = true
				} else if rule.ActionType == entities.ActionModify {
					// Rule triggered a score modification (e.g., -0.2 = reduce by 20%)
					ctx.Multipliers[alt.ID] *= (1.0 + rule.ActionValue)
				}
			}
		}
	}

	return ctx
}
