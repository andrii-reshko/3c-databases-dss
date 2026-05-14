package persistence

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func NewDB(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sqlx.DB) error {
	alternativeSchema := `
    CREATE TABLE IF NOT EXISTS alternatives (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT
    );`

	criterionSchema := `
    CREATE TABLE IF NOT EXISTS criteria (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT,
        input_min REAL NOT NULL DEFAULT 0,
        input_max REAL NOT NULL DEFAULT 1,
        output_min REAL NOT NULL DEFAULT 0,
        output_max REAL NOT NULL DEFAULT 1,
        weight REAL NOT NULL DEFAULT 0
    );`

	evaluationSchema := `
    CREATE TABLE IF NOT EXISTS evaluations (
        alternative_id INTEGER NOT NULL,
        criterion_id INTEGER NOT NULL,
        value REAL NOT NULL,
        PRIMARY KEY (alternative_id, criterion_id),
        FOREIGN KEY (alternative_id) REFERENCES alternatives(id) ON DELETE CASCADE,
        FOREIGN KEY (criterion_id) REFERENCES criteria(id) ON DELETE CASCADE
    );`

	ruleSchema := `
    CREATE TABLE IF NOT EXISTS rules (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        criterion_id INTEGER NOT NULL,
        operator TEXT NOT NULL,
        value REAL NOT NULL,
        action_type TEXT NOT NULL,
        action_value REAL NOT NULL,
        FOREIGN KEY (criterion_id) REFERENCES criteria(id) ON DELETE CASCADE
    );`

	if _, err := db.Exec(alternativeSchema); err != nil {
		return err
	}
	if _, err := db.Exec(criterionSchema); err != nil {
		return err
	}
	if _, err := db.Exec(evaluationSchema); err != nil {
		return err
	}
	if _, err := db.Exec(ruleSchema); err != nil {
		return err
	}

	return nil
}
