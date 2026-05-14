package entities

type ActionType string

const (
	ActionExclude ActionType = "exclude"
	ActionModify  ActionType = "modify"
)

type Rule struct {
	ID          int64      `db:"id"`
	Name        string     `db:"name"`
	CriterionID int64      `db:"criterion_id"`
	Operator    string     `db:"operator"` // ">", "<", ">=", "<=", "=="
	Value       float64    `db:"value"`
	ActionType  ActionType `db:"action_type"`
	ActionValue float64    `db:"action_value"` // e.g., -0.2 for -20%, or 0 for exclude
}

func (r *Rule) Matches(val float64) bool {
	switch r.Operator {
	case ">":
		return val > r.Value
	case "<":
		return val < r.Value
	case ">=":
		return val >= r.Value
	case "<=":
		return val <= r.Value
	case "==":
		return val == r.Value
	}
	return false
}
