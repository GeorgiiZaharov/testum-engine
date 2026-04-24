package gethardtasks

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"testum-engine/app/internal/adapter/db"
	"testum-engine/app/internal/repository/access"
	"testum-engine/app/internal/repository/task"
)

type factory struct {
	db  *db.DB
	log *zap.Logger
}

func NewFactory(db *db.DB, log *zap.Logger) repoFactory {
	return &factory{
		db:  db,
		log: log,
	}
}

func (f *factory) WithTx(ctx context.Context, fn func(r repositories) error) error {
	return f.db.WithTx(ctx, func(tx *sql.Tx) error {
		repos := repositories{
			Access: access.NewRepository(tx, f.log),
			Task:   task.NewRepository(tx, f.log),
		}
		return fn(repos)
	})
}
