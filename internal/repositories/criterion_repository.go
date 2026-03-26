package repositories

import (
	"dss/internal/domain/entities"

	"github.com/jmoiron/sqlx"
)

type sqliteCriterionRepository struct {
	db *sqlx.DB
}

func NewCriterionRepository(db *sqlx.DB) CriterionRepository {
	return &sqliteCriterionRepository{db: db}
}

func (r *sqliteCriterionRepository) FindAll() ([]*entities.Criterion, error) {
	var criteria []*entities.Criterion
	err := r.db.Select(&criteria, "SELECT * FROM criteria ORDER BY id")
	return criteria, err
}

func (r *sqliteCriterionRepository) FindByID(id int64) (*entities.Criterion, error) {
	var criterion entities.Criterion
	err := r.db.Get(&criterion, "SELECT * FROM criteria WHERE id = ?", id)
	return &criterion, err
}

func (r *sqliteCriterionRepository) Create(c *entities.Criterion) (int64, error) {
	res, err := r.db.Exec("INSERT INTO criteria (name, type, description) VALUES (?, ?, ?)", c.Name, c.Type, c.Description)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteCriterionRepository) Update(c *entities.Criterion) error {
	_, err := r.db.Exec("UPDATE criteria SET name = ?, type = ?, description = ? WHERE id = ?", c.Name, c.Type, c.Description, c.ID)
	return err
}

func (r *sqliteCriterionRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM criteria WHERE id = ?", id)
	return err
}
