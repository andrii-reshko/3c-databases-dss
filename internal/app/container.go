package app

import (
	"dss/internal/domain/entities"
	"dss/internal/http"
	"dss/internal/persistence"
	"dss/internal/repositories"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

type Container struct {
	DB                 *sqlx.DB
	AlternativeRepo    repositories.AlternativeRepository
	CriterionRepo      repositories.CriterionRepository
	EvaluationRepo     repositories.EvaluationRepository
	AlternativeHandler *http.AlternativeHandler
	CriterionHandler   *http.CriterionHandler
	EvaluationHandler  *http.EvaluationHandler
}

func NewContainer() (*Container, error) {
	getwd, _ := os.Getwd()
	log.Printf("cwd: %s", getwd)
	db, err := persistence.NewDB(filepath.Join(getwd, "dss.db"))
	log.Printf("db: %v", db)
	if err != nil {
		return nil, err
	}

	if err := persistence.Migrate(db); err != nil {
		return nil, err
	}

	seedData(db)

	altRepo := repositories.NewAlternativeRepository(db)
	critRepo := repositories.NewCriterionRepository(db)
	evalRepo := repositories.NewEvaluationRepository(db)

	return &Container{
		DB:                 db,
		AlternativeRepo:    altRepo,
		CriterionRepo:      critRepo,
		EvaluationRepo:     evalRepo,
		AlternativeHandler: http.NewAlternativeHandler(altRepo),
		CriterionHandler:   http.NewCriterionHandler(critRepo),
		EvaluationHandler:  http.NewEvaluationHandler(altRepo, critRepo, evalRepo),
	}, nil
}

func seedData(db *sqlx.DB) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM alternatives").Scan(&count)
	if count > 0 {
		return
	}

	alternatives := []*entities.Alternative{
		entities.NewAlternative("Розробка MVP вебдодатку", "США, 50 USD/h, опубліковано 15 хв тому"),
		entities.NewAlternative("Інтеграція платіжної системи Stripe", "Велика Британія, 30 USD/h, опубліковано 2 години тому"),
		entities.NewAlternative("Виправлення багів на WordPress-сайті", "Індія, 10 USD/h, опубліковано 5 хв тому"),
	}

	criteria := []*entities.Criterion{
		entities.NewCriterion("Час публікації", entities.TypeMinimize, "Час від публікації в годинах"),
		entities.NewCriterion("Відповідність експертизі", entities.TypeMaximize, "Збіг з навичками команди"),
		entities.NewCriterion("Ставка", entities.TypeMaximize, "USD/hour"),
		entities.NewCriterion("Конкуренція", entities.TypeMinimize, "Кількість поданих заявок"),
		entities.NewCriterion("Активні інтерв'ю", entities.TypeMinimize, "Кількість фрилансерів на інтерв'ю"),
		entities.NewCriterion("Відсоток найму", entities.TypeMaximize, "% успішно закритих вакансій"),
		entities.NewCriterion("Середня виплачена ставка", entities.TypeMaximize, "Історична ставка клієнта"),
		entities.NewCriterion("Негативні відгуки", entities.TypeMinimize, "Кількість негативних відгуків"),
		entities.NewCriterion("Тривалість проєкту", entities.TypeMaximize, "Короткостроковий або довгостроковий"),
	}

	for _, alt := range alternatives {
		_, _ = db.Exec("INSERT INTO alternatives (name, description) VALUES (?, ?)",
			alt.Name, alt.Description)
	}

	for _, c := range criteria {
		_, _ = db.Exec("INSERT INTO criteria (name, type, description) VALUES (?, ?, ?)",
			c.Name, c.Type, c.Description)
	}

	evalData := []struct {
		altID  int
		critID int
		value  float64
	}{
		// alt 1
		{1, 1, 0.01},  // Pub time
		{1, 2, 1.0},   // Expertise match
		{1, 4, 0.9},   // Rate
		{1, 5, 0.6},   // Proposals
		{1, 6, 0.3},   // Active interviews
		{1, 7, 0.75},  // Hire rate
		{1, 8, 0.6},   // Avg paid rate
		{1, 9, 0.2},   // Negative reviews
		{1, 10, 0.75}, // Duration
		// alt 2
		{2, 1, 0.5},
		{2, 2, 0.7},
		{2, 4, 0.56},
		{2, 5, 0.73},
		{2, 6, 0.2},
		{2, 7, 0.6},
		{2, 8, 0.4},
		{2, 9, 0.3},
		{2, 10, 0.5},
		// alt 3
		{3, 1, 0.003},
		{3, 2, 0.2},
		{3, 4, 0.16},
		{3, 5, 0.17},
		{3, 6, 0.5},
		{3, 7, 0.4},
		{3, 8, 0.2},
		{3, 9, 0.5},
		{3, 10, 0.1},
	}

	for _, e := range evalData {
		_, _ = db.Exec("INSERT INTO evaluations (alternative_id, criterion_id, value) VALUES (?, ?, ?)",
			e.altID, e.critID, e.value)
	}
}
