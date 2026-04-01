package entities

type Evaluation struct {
	AlternativeID int64   `db:"alternative_id"`
	CriterionID   int64   `db:"criterion_id"`
	Value         float64 `db:"value"`
}
