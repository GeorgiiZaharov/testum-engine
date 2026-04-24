package db

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type MigrationConfig struct {
	Dir     string
	Dialect string
}

func RunMigrations(db *sql.DB, cfg MigrationConfig) error {
	if err := goose.SetDialect(cfg.Dialect); err != nil {
		return err
	}

	if err := goose.Up(db, cfg.Dir); err != nil {
		return err
	}

	return nil
}
