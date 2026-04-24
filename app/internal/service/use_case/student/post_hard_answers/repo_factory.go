package posthardanswers

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"

	accessrepo "testum-engine/app/internal/repository/access"
	answerrepo "testum-engine/app/internal/repository/answer"
	resultrepo "testum-engine/app/internal/repository/result"
	studenttestrepo "testum-engine/app/internal/repository/student_test"
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
			User:        userrepo.NewRepository(tx, f.log),
			Access:      accessrepo.NewRepository(tx, f.log),
			Answer:      answerrepo.NewRepository(tx, f.log),
			Result:      resultrepo.NewRepository(tx, f.log),
			StudentTest: studenttestrepo.NewRepository(tx, f.log),
		}
		return fn(repos)
	})
}
