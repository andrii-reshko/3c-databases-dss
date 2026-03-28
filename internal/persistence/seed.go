package persistence

import (
	"dss/internal/domain/entities"

	"github.com/jmoiron/sqlx"
)

func SeedData(db *sqlx.DB) error {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM alternatives").Scan(&count)
	if count > 0 {
		return nil
	}

	alternatives := []*entities.Alternative{
		entities.NewAlternative("Розробка MVP вебдодатку", "США, 50 USD/h, опубліковано 15 хв тому"),
		entities.NewAlternative("Інтеграція платіжної системи Stripe", "Велика Британія, 30 USD/h, опубліковано 2 години тому"),
		entities.NewAlternative("Виправлення багів на WordPress-сайті", "Індія, 10 USD/h, опубліковано 5 хв тому"),
	}

	criteria := []*entities.Criterion{
		{ID: 1, Name: "Час публікації", Type: entities.TypeMinimize, Description: "Час від публікації, h", InputMin: 0, InputMax: 24, OutputMin: 1, OutputMax: 0, Weight: 0.20},
		{ID: 2, Name: "Відповідність експертизі", Type: entities.TypeMaximize, Description: "Збіг з навичками команди, %", InputMin: 0, InputMax: 100, OutputMin: 0, OutputMax: 1, Weight: 0.12},
		{ID: 3, Name: "Ставка", Type: entities.TypeMaximize, Description: "USD/h", InputMin: 10, InputMax: 60, OutputMin: 0, OutputMax: 1, Weight: 0.10},
		{ID: 4, Name: "Конкуренція", Type: entities.TypeMinimize, Description: "Кількість поданих заявок", InputMin: 0, InputMax: 50, OutputMin: 1, OutputMax: 0, Weight: 0.15},
		{ID: 5, Name: "Активні інтерв'ю", Type: entities.TypeMinimize, Description: "Кількість фрилансерів на інтерв'ю", InputMin: 0, InputMax: 5, OutputMin: 1, OutputMax: 0, Weight: 0.12},
		{ID: 6, Name: "Відсоток найму", Type: entities.TypeMaximize, Description: "Кількість закритих контрактів, %", InputMin: 0, InputMax: 100, OutputMin: 0, OutputMax: 1, Weight: 0.10},
		{ID: 7, Name: "Середня виплачена ставка", Type: entities.TypeMaximize, Description: "Середня ставка клієнта, USD/h", InputMin: 10, InputMax: 60, OutputMin: 0, OutputMax: 1, Weight: 0.08},
		{ID: 8, Name: "Негативні відгуки", Type: entities.TypeMinimize, Description: "Відсоток негативних відгуків", InputMin: 0, InputMax: 100, OutputMin: 1, OutputMax: 0, Weight: 0.08},
		{ID: 9, Name: "Тривалість проєкту", Type: entities.TypeMaximize, Description: "Тривалість у місяцях", InputMin: 1, InputMax: 12, OutputMin: 0, OutputMax: 1, Weight: 0.05},
	}

	for _, alt := range alternatives {
		_, err := db.Exec("INSERT INTO alternatives (name, description) VALUES (?, ?)",
			alt.Name, alt.Description)
		if err != nil {
			return err
		}
	}

	for _, c := range criteria {
		_, err := db.Exec("INSERT INTO criteria (name, type, description, input_min, input_max, output_min, output_max, weight) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			c.Name, c.Type, c.Description, c.InputMin, c.InputMax, c.OutputMin, c.OutputMax, c.Weight)
		if err != nil {
			return err
		}
	}

	evalData := []struct {
		altID  int
		critID int
		value  float64
	}{
		// alt 1 - Відмінна пропозиція: швидко опублікована, висока ставка, досвідчений клієнт
		{1, 1, 0.5}, // Pub time: 0.5h (опубліковано нещодавно)
		{1, 2, 95},  // Expertise match: 95%
		{1, 3, 50},  // Rate: 50 USD/h (висока ставка!)
		{1, 4, 15},  // Proposals: 15 (низька конкуренція)
		{1, 5, 1},   // Active interviews: 1
		{1, 6, 85},  // Hire %: 85%
		{1, 7, 45},  // Avg rate: 45 USD/h
		{1, 8, 5},   // Negative: 5% (мало негативу)
		{1, 9, 5},   // Duration: 5 місяців
		// alt 2 - Середня пропозиція
		{2, 1, 3},  // Pub time: 3h
		{2, 2, 70}, // Expertise match: 70%
		{2, 3, 35}, // Rate: 35 USD/h
		{2, 4, 30}, // Proposals: 30
		{2, 5, 3},  // Active interviews: 3
		{2, 6, 50}, // Hire %: 50%
		{2, 7, 30}, // Avg rate: 30 USD/h
		{2, 8, 20}, // Negative: 20%
		{2, 9, 3},  // Duration: 3 місяці
		// alt 3 - Слабка пропозиція: давно опублікована, низька ставка, багато конкурентів
		{3, 1, 20}, // Pub time: 20h (давно)
		{3, 2, 40}, // Expertise match: 40%
		{3, 3, 20}, // Rate: 20 USD/h (низька ставка)
		{3, 4, 50}, // Proposals: 50 (висока конкуренція)
		{3, 5, 5},  // Active interviews: 5
		{3, 6, 20}, // Hire %: 20%
		{3, 7, 25}, // Avg rate: 25 USD/h
		{3, 8, 40}, // Negative: 40% (багато негативу)
		{3, 9, 1},  // Duration: 1 місяць
	}

	for _, e := range evalData {
		_, err := db.Exec("INSERT INTO evaluations (alternative_id, criterion_id, value) VALUES (?, ?, ?)",
			e.altID, e.critID, e.value)
		if err != nil {
			return err
		}
	}

	return nil
}
