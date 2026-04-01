package repositories

import (
	"dss/internal/domain/entities"

	"github.com/jmoiron/sqlx"
)

type sqliteEvaluationRepository struct {
	db *sqlx.DB
}

func NewEvaluationRepository(db *sqlx.DB) EvaluationRepository {
	return &sqliteEvaluationRepository{db: db}
}

func (r *sqliteEvaluationRepository) FindAll() ([]*entities.Evaluation, error) {
	var evaluations []*entities.Evaluation
	err := r.db.Select(&evaluations, "SELECT * FROM evaluations")
	return evaluations, err
}

func (r *sqliteEvaluationRepository) UpsertBatch(evaluations []*entities.Evaluation) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO evaluations (alternative_id, criterion_id, value) VALUES (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, e := range evaluations {
		_, err := stmt.Exec(e.AlternativeID, e.CriterionID, e.Value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
