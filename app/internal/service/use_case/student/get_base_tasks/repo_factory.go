package getbasetasks

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
	accessrepo "testum-engine/app/internal/repository/access"
	resultrepo "testum-engine/app/internal/repository/result"
	studenttestrepo "testum-engine/app/internal/repository/student_test"
	taskrepo "testum-engine/app/internal/repository/task"
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
			Access:      accessrepo.NewRepository(tx, f.log),
			StudentTest: studenttestrepo.NewRepository(tx, f.log),
			Result:      resultrepo.NewRepository(tx, f.log),
			Task:        taskrepo.NewRepository(tx, f.log),
		}
		return fn(repos)
	})
}
