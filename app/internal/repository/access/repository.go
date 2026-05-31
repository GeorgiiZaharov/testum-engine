package access

import (
	"context"
	"testum-engine/app/internal/adapter/db"

	"go.uber.org/zap"
)

// ================= INTERFACE =================

type Repository interface {
	HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error)
	HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error)
	GiveAccess(ctx context.Context, testID int, group string) (bool, error)
	TakeAccess(ctx context.Context, testID int, group string) (bool, error)
}

// ================= IMPLEMENTATION =================

type repository struct {
	db  db.Executor
	log *zap.Logger
}

func NewRepository(db db.Executor, log *zap.Logger) Repository {
	return &repository{
		db:  db,
		log: log,
	}
}

// ================= HAS LECTURER ACCESS =================

func (r *repository) HasLecturerAccess(ctx context.Context, userID int, testID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM tests
			WHERE id = ?
			  AND owner_id = ?
		)
	`

	var exists bool

	err := r.db.QueryRowContext(ctx, query, testID, userID).Scan(&exists)
	if err != nil {
		r.log.Error("HasLecturerAccess query failed",
			zap.Error(err),
			zap.Int("user_id", userID),
			zap.Int("test_id", testID),
		)
		return false, err
	}

	return exists, nil
}

// ================= HAS STUDENT ACCESS =================

func (r *repository) HasStudentAccess(ctx context.Context, userID int, testID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users u
			JOIN test_permissions tp ON u."group" = tp."group"
			WHERE u.id = ?
			  AND tp.test_id = ?
		)
	`

	var exists bool

	err := r.db.QueryRowContext(ctx, query, userID, testID).Scan(&exists)
	if err != nil {
		r.log.Error("HasStudentAccess query failed",
			zap.Error(err),
			zap.Int("user_id", userID),
			zap.Int("test_id", testID),
		)
		return false, err
	}

	return exists, nil
}

// ================= GIVE ACCESS =================

func (r *repository) GiveAccess(ctx context.Context, testID int, group string) (bool, error) {
	query := `
		INSERT INTO test_permissions (test_id, "group")
		VALUES (?, ?)
	`

	_, err := r.db.ExecContext(ctx, query, testID, group)
	if err != nil {
		r.log.Error("GiveAccess failed",
			zap.Error(err),
			zap.Int("test_id", testID),
			zap.String("group", group),
		)
		return false, err
	}

	return true, nil
}

// ================= TAKE ACCESS =================

func (r *repository) TakeAccess(ctx context.Context, testID int, group string) (bool, error) {
	query := `
		DELETE FROM test_permissions
		WHERE test_id = ?
		  AND "group" = ?
	`

	res, err := r.db.ExecContext(ctx, query, testID, group)
	if err != nil {
		r.log.Error("TakeAccess failed",
			zap.Error(err),
			zap.Int("test_id", testID),
			zap.String("group", group),
		)
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		r.log.Error("TakeAccess rows affected failed",
			zap.Error(err),
		)
		return false, err
	}

	if affected == 0 {
		return false, ErrAccessNotFound
	}

	return true, nil
}
