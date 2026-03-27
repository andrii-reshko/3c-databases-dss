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
	InputMin    float64       `db:"input_min"`
	InputMax    float64       `db:"input_max"`
	OutputMin   float64       `db:"output_min"`
	OutputMax   float64       `db:"output_max"`
}

func NewCriterion(name string, cType CriterionType, description string) *Criterion {
	return &Criterion{
		Name:        name,
		Description: description,
		Type:        cType,
		InputMin:    0,
		InputMax:    1,
		OutputMin:   0,
		OutputMax:   1,
	}
}

func (c *Criterion) Normalize(value float64) float64 {
	if c.InputMax == c.InputMin {
		if c.Type == TypeMinimize {
			if value < c.InputMin {
				return c.OutputMax
			}
			return c.OutputMin
		}
		if value > c.InputMin {
			return c.OutputMax
		}
		return c.OutputMin
	}

	t := clamp(norm(value, c.InputMin, c.InputMax), 0, 1)

	if c.Type == TypeMinimize {
		return lerp(t, c.OutputMax, c.OutputMin)
	}

	return lerp(t, c.OutputMin, c.OutputMax)
}

func norm(value, min, max float64) float64 {
	return (value - min) / (max - min)
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}
