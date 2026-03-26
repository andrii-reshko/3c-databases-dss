package entities

type Alternative struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func NewAlternative(name, description string) *Alternative {
	return &Alternative{
		Name:        name,
		Description: description,
	}
}
