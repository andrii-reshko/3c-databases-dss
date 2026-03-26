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
        type TEXT NOT NULL CHECK(type IN ('maximize', 'minimize')),
        description TEXT
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

	if _, err := db.Exec(alternativeSchema); err != nil {
		return err
	}
	if _, err := db.Exec(criterionSchema); err != nil {
		return err
	}
	if _, err := db.Exec(evaluationSchema); err != nil {
		return err
	}

	return nil
}
