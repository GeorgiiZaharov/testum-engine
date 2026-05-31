package studenttest

import (
	"context"
	"testum-engine/app/internal/adapter/db"

	"go.uber.org/zap"
)

type Repository interface {
	StartTest(ctx context.Context, params StartTestParams) (bool, error)
	FinishTest(ctx context.Context, params FinishTestParams) (bool, error)
	GetActiveTests(ctx context.Context, userID int) ([]StudentActiveTestInfo, error)
	GetFinishedTests(ctx context.Context, userID int) ([]StudentFinishedTestInfo, error)
}

type repository struct {
	db  db.Executor
	log *zap.Logger
}

func NewRepository(db db.Executor, log *zap.Logger) Repository {
	return &repository{db: db, log: log}
}

// ================= START TEST =================

func (r *repository) StartTest(ctx context.Context, params StartTestParams) (bool, error) {
	query := `
		INSERT INTO student_tests (student_id, test_id, "group", date_start)
		VALUES (
			?,
			?,
			(SELECT "group" FROM users WHERE id = ?),
			CURRENT_TIMESTAMP
		)
	`

	_, err := r.db.ExecContext(ctx, query,
		params.UserID,
		params.TestID,
		params.UserID,
	)

	if err != nil {
		r.log.Error("failed to start test", zap.Error(err))
		return false, err
	}

	return true, nil
}

// ================= FINISH TEST =================

func (r *repository) FinishTest(ctx context.Context, params FinishTestParams) (bool, error) {
	query := `
		UPDATE student_tests
		SET mark = ?,
		    success_rate = ?,
		    date_end = CURRENT_TIMESTAMP
		WHERE student_id = ? AND test_id = ?
	`

	res, err := r.db.ExecContext(ctx, query,
		params.Result.Mark,
		params.Result.SuccessRate,
		params.UserID,
		params.TestID,
	)
	if err != nil {
		return false, ErrFinishFailed
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, ErrStudentTestNotFound
	}

	return true, nil
}

// ================= GET ACTIVE TESTS =================

func (r *repository) GetActiveTests(ctx context.Context, userID int) ([]StudentActiveTestInfo, error) {
	query := `
		SELECT
			t.id,
			t.name,
			u.name,
			(SELECT COUNT(*) FROM tasks WHERE test_id = t.id),
			(SELECT COUNT(*) FROM tasks WHERE test_id = t.id AND is_hard = 1),
			st.date_start
		FROM users us
		JOIN test_permissions tp ON tp."group" = us."group"
		JOIN tests t ON t.id = tp.test_id
		JOIN users u ON u.id = t.owner_id
		LEFT JOIN student_tests st ON st.test_id = t.id AND st.student_id = us.id
		WHERE us.id = ?
		  AND st.mark IS NULL
		ORDER BY st.date_start DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []StudentActiveTestInfo

	for rows.Next() {
		var item StudentActiveTestInfo

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.LecturerName,
			&item.CntQuestions,
			&item.CntHardQuestions,
			&item.DateStart,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, rows.Err()
}

// ================= GET FINISHED TESTS =================

func (r *repository) GetFinishedTests(ctx context.Context, userID int) ([]StudentFinishedTestInfo, error) {
	query := `
		SELECT
			t.id,
			t.name,
			u.name,
			(SELECT COUNT(*) FROM tasks WHERE test_id = t.id),
			(SELECT COUNT(*) FROM tasks WHERE test_id = t.id AND is_hard = 1),
			st.mark,
			st.success_rate,
			st.date_start,
			st.date_end
		FROM student_tests st
		JOIN tests t ON t.id = st.test_id
		JOIN users u ON u.id = t.owner_id
		WHERE st.student_id = ?
		  AND st.mark IS NOT NULL
		ORDER BY st.date_end DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []StudentFinishedTestInfo

	for rows.Next() {
		var item StudentFinishedTestInfo

		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.LecturerName,
			&item.CntQuestions,
			&item.CntHardQuestions,
			&item.Mark,
			&item.SuccessRate,
			&item.DateStart,
			&item.DateEnd,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, rows.Err()
}
