package bootstrap

import (
	"testum-engine/app/internal/adapter/db"
)

type Config struct {
	DBOptions  db.DBOptions
	Migrations string
}

func Setup(cfg Config) (*db.DB, func(), error) {
	database, err := db.NewDB(cfg.DBOptions)
	if err != nil {
		return nil, nil, err
	}

	sqlDB := database.DB.DB

	if err := db.RunMigrations(sqlDB, db.MigrationConfig{
		Dir:     cfg.Migrations,
		Dialect: "mysql",
	}); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = sqlDB.Close()
	}

	return database, cleanup, nil
}
