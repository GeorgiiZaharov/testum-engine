package file

import (
	"context"
	"testum-engine/app/internal/adapter/db"

	"go.uber.org/zap"
)

type FileRepository interface {
	GetAllTestFiles(ctx context.Context, lecturerID int) ([]string, error)
}

type repository struct {
	db  db.Executor
	log *zap.Logger
}

func NewRepository(db db.Executor, log *zap.Logger) FileRepository {
	return &repository{
		db:  db,
		log: log,
	}
}

// ================= TEST FILES =================

func (r *repository) GetAllTestFiles(ctx context.Context, lecturerID int) ([]string, error) {
	query := `
		SELECT file_name
		FROM tests
		WHERE owner_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, lecturerID)
	if err != nil {
		r.log.Error("GetAllTestFiles query failed",
			zap.Error(err),
			zap.Int("lecturer_id", lecturerID),
		)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var files []string

	for rows.Next() {
		var file string

		if err := rows.Scan(&file); err != nil {
			r.log.Error("GetAllTestFiles scan failed",
				zap.Error(err),
				zap.Int("lecturer_id", lecturerID),
			)
			return nil, err
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("GetAllTestFiles rows error",
			zap.Error(err),
			zap.Int("lecturer_id", lecturerID),
		)
		return nil, err
	}

	if len(files) == 0 {
		return files, nil
	}

	return files, nil
}
