package repositories

import (
	"dss/internal/domain/entities"

	"github.com/jmoiron/sqlx"
)

type sqliteAlternativeRepository struct {
	db *sqlx.DB
}

func NewAlternativeRepository(db *sqlx.DB) AlternativeRepository {
	return &sqliteAlternativeRepository{db: db}
}

func (r *sqliteAlternativeRepository) FindAll() ([]*entities.Alternative, error) {
	var alternatives []*entities.Alternative
	err := r.db.Select(&alternatives, "SELECT * FROM alternatives ORDER BY id")
	return alternatives, err
}

func (r *sqliteAlternativeRepository) FindByID(id int64) (*entities.Alternative, error) {
	var alt entities.Alternative
	err := r.db.Get(&alt, "SELECT * FROM alternatives WHERE id = ?", id)
	return &alt, err
}

func (r *sqliteAlternativeRepository) Create(alt *entities.Alternative) (int64, error) {
	res, err := r.db.Exec("INSERT INTO alternatives (name, description) VALUES (?, ?)", alt.Name, alt.Description)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteAlternativeRepository) Update(alt *entities.Alternative) error {
	_, err := r.db.Exec("UPDATE alternatives SET name = ?, description = ? WHERE id = ?", alt.Name, alt.Description, alt.ID)
	return err
}

func (r *sqliteAlternativeRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM alternatives WHERE id = ?", id)
	return err
}
