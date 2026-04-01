package entities

type Criterion struct {
	ID          int64   `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	InputMin    float64 `db:"input_min"`
	InputMax    float64 `db:"input_max"`
	OutputMin   float64 `db:"output_min"`
	OutputMax   float64 `db:"output_max"`
	Weight      float64 `db:"weight"`
}

func (c *Criterion) Normalize(value float64) float64 {
	if c.InputMax == c.InputMin {
		// OutputMax < OutputMin -> minimization
		if c.OutputMax < c.OutputMin {
			if value < c.InputMin {
				return c.OutputMin
			}
			return c.OutputMax
		}
		if value > c.InputMin {
			return c.OutputMax // -> maximization
		}
		return c.OutputMin
	}

	t := clamp(norm(value, c.InputMin, c.InputMax), 0, 1)

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
