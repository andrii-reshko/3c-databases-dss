package app

import (
	"dss/internal/domain/analytics"
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
	RuleRepo           repositories.RuleRepository
	EvaluationRepo     repositories.EvaluationRepository
	AlternativeHandler *http.AlternativeHandler
	CriterionHandler   *http.CriterionHandler
	EvaluationHandler  *http.EvaluationHandler
	RankingHandler     *http.RankingHandler
	ImportHandler      *http.ImportHandler
	RuleHandler        *http.RuleHandler
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

	altRepo := repositories.NewAlternativeRepository(db)
	critRepo := repositories.NewCriterionRepository(db)
	evalRepo := repositories.NewEvaluationRepository(db)
	ruleRepo := repositories.NewRuleRepository(db)

	if err := persistence.SeedData(altRepo, critRepo, evalRepo); err != nil {
		return nil, err
	}

	votingSvc := analytics.NewVotingService()
	expertSvc := analytics.NewExpertEvaluationService()

	return &Container{
		DB:                 db,
		AlternativeRepo:    altRepo,
		CriterionRepo:      critRepo,
		EvaluationRepo:     evalRepo,
		RuleRepo:           ruleRepo,
		AlternativeHandler: http.NewAlternativeHandler(altRepo),
		CriterionHandler:   http.NewCriterionHandler(critRepo),
		EvaluationHandler:  http.NewEvaluationHandler(altRepo, critRepo, evalRepo),
		RankingHandler:     http.NewRankingHandler(altRepo, critRepo, evalRepo, ruleRepo),
		ImportHandler:      http.NewImportHandler(critRepo, evalRepo, votingSvc, expertSvc),
		RuleHandler:        http.NewRuleHandler(ruleRepo, critRepo),
	}, nil
}

func (c *Container) Close() {
	err := c.DB.Close()
	if err != nil {
		log.Fatalf("Error closing DB: %v", err)
	}
}
