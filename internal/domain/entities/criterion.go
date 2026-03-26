package entities

type CriterionType string

const (
	TypeMaximize CriterionType = "maximize"
	TypeMinimize CriterionType = "minimize"
)

type Criterion struct {
	ID          int64         `db:"id"`
	Name        string        `db:"name"`
	Type        CriterionType `db:"type"`
	Description string        `db:"description"`
}

func NewCriterion(name string, cType CriterionType, description string) *Criterion {
	return &Criterion{
		Name:        name,
		Type:        cType,
		Description: description,
	}
}
