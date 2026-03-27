package app

import (
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
	if err := persistence.SeedData(db); err != nil {
		return nil, err
	}

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
