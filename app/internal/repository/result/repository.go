package result

import (
	"context"
	"database/sql"
	"errors"
	"testum-engine/app/internal/adapter/db"
	"time"

	"go.uber.org/zap"
)

// ================= TYPES =================

type Repository interface {
	GetGroupResult(ctx context.Context, testID int, group string, year int) ([]StudentResult, error)
	GetStudentResult(ctx context.Context, userID int, testID int) (TestResult, error)
}

type repository struct {
	db  db.Executor
	log *zap.Logger

	now func() time.Time
}

func NewRepository(db db.Executor, log *zap.Logger) Repository {
	return &repository{
		db:  db,
		log: log,
		now: time.Now,
	}
}

// ================= ACADEMIC RANGE =================

func (r *repository) getAcademicRange(year int) (time.Time, time.Time) {
	now := r.now()

	currentYear := now.Year()
	august := time.Date(currentYear, time.August, 1, 0, 0, 0, 0, now.Location())

	if now.Before(august) {
		start := august.AddDate(-1-year, 0, 0)
		end := august.AddDate(-year, 0, 0)
		return start, end
	}

	start := august.AddDate(-year, 0, 0)
	end := august.AddDate(1-year, 0, 0)
	return start, end
}

// ================= GetGroupResult =================

func (r *repository) GetGroupResult(ctx context.Context, testID int, group string, year int) ([]StudentResult, error) {
	start, end := r.getAcademicRange(year)

	query := `
		SELECT 
			u.id,
			u.name,
			u.login,
			u.mail,
			st.mark,
			st.success_rate,
			st.date_start,
			st.date_end
		FROM users u
		LEFT JOIN student_tests st 
			ON st.student_id = u.id
			AND st.test_id = ?
			AND st.date_start >= ?
			AND st.date_start < ?
		WHERE u.` + "`group`" + ` = ?
	  	AND u.date_modified >= ?
	`

	rows, err := r.db.QueryContext(ctx, query,
		testID,
		start,
		end,
		group,
		start,
	)
	if err != nil {
		r.log.Error("GetGroupResult query failed",
			zap.Error(err),
			zap.Int("test_id", testID),
			zap.String("group", group),
			zap.Int("year", year),
		)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	results := make([]StudentResult, 0)

	for rows.Next() {
		var sr StudentResult

		if err := rows.Scan(
			&sr.UserID,
			&sr.Name,
			&sr.Login,
			&sr.Mail,
			&sr.Mark,
			&sr.SuccessRate,
			&sr.DateStart,
			&sr.DateEnd,
		); err != nil {
			r.log.Error("GetGroupResult scan failed", zap.Error(err))
			return nil, err
		}

		results = append(results, sr)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("GetGroupResult rows error", zap.Error(err))
		return nil, err
	}

	return results, nil
}

// ================= GetStudentResult =================

func (r *repository) GetStudentResult(ctx context.Context, userID int, testID int) (TestResult, error) {
	query := `
		SELECT 
			mark,
			success_rate,
			date_start,
			date_end
		FROM student_tests
		WHERE student_id = ? AND test_id = ?
		LIMIT 1
	`

	var result TestResult

	err := r.db.QueryRowContext(ctx, query, userID, testID).Scan(
		&result.Mark,
		&result.SuccessRate,
		&result.DateStart,
		&result.DateEnd,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Warn("GetStudentResult not found",
				zap.Int("user_id", userID),
				zap.Int("test_id", testID),
			)
			return TestResult{}, ErrResultNotFound
		}

		r.log.Error("GetStudentResult failed",
			zap.Error(err),
			zap.Int("user_id", userID),
			zap.Int("test_id", testID),
		)
		return TestResult{}, err
	}

	return result, nil
}
