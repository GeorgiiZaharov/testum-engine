package user

import (
	"context"
	"database/sql"
	"errors"
	"testum-engine/app/internal/adapter/db"

	"go.uber.org/zap"
)

type Repository interface {
	Upsert(ctx context.Context, params CreateUserParams) (int, error)
	GetByID(ctx context.Context, userID int) (User, error)
	GetLecturers(ctx context.Context) ([]User, error)
	CreateLecturer(ctx context.Context, userID int) error
	DeleteLecturer(ctx context.Context, userID int) error
}

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

// ================= UPSERT (SQLite version) =================
// В SQLite используем ON CONFLICT(login)
func (r *repository) Upsert(ctx context.Context, params CreateUserParams) (int, error) {
	query := `
		INSERT INTO users (login, mail, name, ` + "`group`" + `, is_lecturer, date_created, date_modified)
		VALUES (?, ?, ?, ?, 0, datetime('now'), datetime('now'))
		ON CONFLICT(login) DO UPDATE SET
			mail = excluded.mail,
			name = excluded.name,
			` + "`group`" + ` = excluded.` + "`group`" + `,
			date_modified = datetime('now')
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		params.Login,
		params.Mail,
		params.Name,
		params.Group,
	).Scan(&id)

	if err != nil {
		r.log.Error("Upsert failed",
			zap.Error(err),
			zap.String("login", params.Login),
		)
		return 0, err
	}

	return id, nil
}

// ================= GET BY ID =================
func (r *repository) GetByID(ctx context.Context, userID int) (User, error) {
	query := `
		SELECT id, login, mail, name, ` + "`group`" + `, is_lecturer, date_created, date_modified
		FROM users
		WHERE id = ?
	`

	var u User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&u.ID,
		&u.Login,
		&u.Mail,
		&u.Name,
		&u.Group,
		&u.IsLecturer,
		&u.DateCreated,
		&u.DateModified,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}

		r.log.Error("GetByID failed",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		return User{}, err
	}

	return u, nil
}

// ================= GET LECTURERS =================
func (r *repository) GetLecturers(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, login, mail, name, ` + "`group`" + `, is_lecturer, date_created, date_modified
		FROM users
		WHERE is_lecturer = 1
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.log.Error("GetLecturers query failed", zap.Error(err))
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []User

	for rows.Next() {
		var u User

		err := rows.Scan(
			&u.ID,
			&u.Login,
			&u.Mail,
			&u.Name,
			&u.Group,
			&u.IsLecturer,
			&u.DateCreated,
			&u.DateModified,
		)
		if err != nil {
			r.log.Error("GetLecturers scan failed", zap.Error(err))
			return nil, err
		}

		result = append(result, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// ================= CREATE LECTURER =================
func (r *repository) CreateLecturer(ctx context.Context, userID int) error {
	query := `
		UPDATE users
		SET is_lecturer = 1
		WHERE id = ? AND is_lecturer = 0
	`

	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		r.log.Error("CreateLecturer failed",
			zap.Error(err),
			zap.Int("id", userID),
		)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows > 0 {
		return nil
	}

	var exists bool
	err = r.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)",
		userID,
	).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrUserNotFound
	}

	return nil
}

// ================= DELETE LECTURER =================
func (r *repository) DeleteLecturer(ctx context.Context, userID int) error {
	query := `
		UPDATE users
		SET is_lecturer = 0, date_modified = datetime('now')
		WHERE id = ?
	`

	res, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		r.log.Error("DeleteLecturer failed",
			zap.Error(err),
			zap.Int("id", userID),
		)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrUserNotFound
	}

	return nil
}
