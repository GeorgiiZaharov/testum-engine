package uploadtest

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	lecturertestrepo "testum-engine/app/internal/repository/lecturer_test"
	userrepo "testum-engine/app/internal/repository/user"
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
			User: userrepo.NewRepository(tx, f.log),
			Test: lecturertestrepo.NewRepository(tx, f.log),
		}
		return fn(repos)
	})
}
