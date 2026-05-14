package repositories

import (
	"dss/internal/domain/entities"

	"github.com/jmoiron/sqlx"
)

type sqliteRuleRepository struct {
	db *sqlx.DB
}

func NewRuleRepository(db *sqlx.DB) RuleRepository {
	return &sqliteRuleRepository{db: db}
}

func (r *sqliteRuleRepository) FindAll() ([]*entities.Rule, error) {
	var rules []*entities.Rule
	err := r.db.Select(&rules, "SELECT * FROM rules ORDER BY id")
	return rules, err
}

func (r *sqliteRuleRepository) FindByID(id int64) (*entities.Rule, error) {
	var rule entities.Rule
	err := r.db.Get(&rule, "SELECT * FROM rules WHERE id = ?", id)
	return &rule, err
}

func (r *sqliteRuleRepository) Create(rule *entities.Rule) (int64, error) {
	res, err := r.db.Exec(`
		INSERT INTO rules (name, criterion_id, operator, value, action_type, action_value) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		rule.Name, rule.CriterionID, rule.Operator, rule.Value, rule.ActionType, rule.ActionValue,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteRuleRepository) Update(rule *entities.Rule) error {
	_, err := r.db.Exec(`
		UPDATE rules SET name = ?, criterion_id = ?, operator = ?, value = ?, action_type = ?, action_value = ? 
		WHERE id = ?`,
		rule.Name, rule.CriterionID, rule.Operator, rule.Value, rule.ActionType, rule.ActionValue, rule.ID,
	)
	return err
}

func (r *sqliteRuleRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM rules WHERE id = ?", id)
	return err
}
