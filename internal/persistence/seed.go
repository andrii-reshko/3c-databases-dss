package persistence

import (
	"dss/internal/domain/entities"
	"dss/internal/repositories"
)

func SeedData(alt repositories.AlternativeRepository, crit repositories.CriterionRepository, eval repositories.EvaluationRepository) error {

	alternatives := []*entities.Alternative{
		entities.NewAlternative("Розробка MVP вебдодатку", "США, 50 USD/h, опубліковано 15 хв тому"),
		entities.NewAlternative("Інтеграція платіжної системи Stripe", "Велика Британія, 30 USD/h, опубліковано 2 години тому"),
		entities.NewAlternative("Виправлення багів на WordPress-сайті", "Індія, 10 USD/h, опубліковано 5 хв тому"),
		entities.NewAlternative("Розробка мобільного додатку на Flutter", "Україна, 25 USD/h, опубліковано 1 годину тому"),
	}

	criteria := []*entities.Criterion{
		{
			ID:          1,
			Name:        "Час публікації",
			Description: "Час від публікації, h",
			InputMin:    1,
			InputMax:    24,
			OutputMin:   1,
			OutputMax:   0,
			Weight:      5,
		},
		{
			ID:          2,
			Name:        "Відповідність експертизі",
			Description: "Збіг з навичками команди, %",
			InputMin:    0,
			InputMax:    100,
			OutputMin:   0,
			OutputMax:   1,
			Weight:      6,
		},
		{
			ID:          3,
			Name:        "Ставка",
			Description: "USD/h",
			InputMin:    10,
			InputMax:    60,
			OutputMin:   0,
			OutputMax:   1,
			Weight:      6,
		},
		{
			ID:          4,
			Name:        "Конкуренція",
			Description: "Кількість поданих заявок",
			InputMin:    0,
			InputMax:    50,
			OutputMin:   1,
			OutputMax:   0,
			Weight:      6,
		},
		{
			ID:          5,
			Name:        "Активні інтерв'ю",
			Description: "Кількість фрилансерів на інтерв'ю",
			InputMin:    0,
			InputMax:    5,
			OutputMin:   1,
			OutputMax:   0,
			Weight:      8,
		},
		{
			ID:          6,
			Name:        "Відсоток найму",
			Description: "Кількість закритих контрактів, %",
			InputMin:    0,
			InputMax:    100,
			OutputMin:   0,
			OutputMax:   1,
			Weight:      3,
		},
		{
			ID:          7,
			Name:        "Середня виплачена ставка",
			Description: "Середня ставка клієнта, USD/h",
			InputMin:    10,
			InputMax:    60,
			OutputMin:   0,
			OutputMax:   1,
			Weight:      4,
		},
		{
			ID:          8,
			Name:        "Негативні відгуки",
			Description: "Відсоток негативних відгуків",
			InputMin:    0,
			InputMax:    100,
			OutputMin:   1,
			OutputMax:   0,
			Weight:      7,
		},
		{
			ID:          9,
			Name:        "Тривалість проєкту",
			Description: "Тривалість у місяцях",
			InputMin:    1,
			InputMax:    12,
			OutputMin:   0,
			OutputMax:   1,
			Weight:      6,
		},
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
		{3, 1, 20},  // Pub time: 20h (давно)
		{3, 2, 100}, // Expertise match: 40%
		{3, 3, 20},  // Rate: 20 USD/h (низька ставка)
		{3, 4, 50},  // Proposals: 50 (висока конкуренція)
		{3, 5, 10},  // Active interviews: 5
		{3, 6, 100}, // Hire %: 20%
		{3, 7, 25},  // Avg rate: 25 USD/h
		{3, 8, 0},   // Negative: 0%
		{3, 9, 1},   // Duration: 1 місяць
		// alt 4 - Хороша пропозиція: нещодавно опублікована, середня ставка, помірна конкуренція
		{4, 1, 1},  // Pub time: 1h
		{4, 2, 80}, // Expertise match: 80%
		{4, 3, 25}, // Rate: 25 USD/h
		{4, 4, 20}, // Proposals: 20
		{4, 5, 2},  // Active interviews: 2
		{4, 6, 70}, // Hire %: 70%
		{4, 7, 28}, // Avg rate: 28 USD/h
		{4, 8, 10}, // Negative: 10%
		{4, 9, 4},  // Duration: 4 місяці
	}

	if stored, _ := alt.FindAll(); len(stored) == 0 {
		for _, e := range alternatives {
			if _, err := alt.Create(e); err != nil {
				return err
			}
		}
	}

	if stored, _ := crit.FindAll(); len(stored) == 0 {
		for _, e := range criteria {
			if _, err := crit.Create(e); err != nil {
				return err
			}
		}
	}

	if stored, _ := eval.FindAll(); len(stored) == 0 {
		batch := make([]*entities.Evaluation, 0, len(evalData))
		for _, e := range evalData {
			batch = append(batch, &entities.Evaluation{
				AlternativeID: int64(e.altID),
				CriterionID:   int64(e.critID),
				Value:         e.value,
			})
		}
		if err := eval.UpsertBatch(batch); err != nil {
			return err
		}
	}

	return nil
}
