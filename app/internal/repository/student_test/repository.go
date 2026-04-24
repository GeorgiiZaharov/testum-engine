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
	return &repository{
		db:  db,
		log: log,
	}
}

// ==================== START TEST ====================
func (r *repository) StartTest(ctx context.Context, params StartTestParams) (bool, error) {
	query := `
		insert into student_tests (student_id, test_id, ` + "`group`" + `, date_start)
		values (?, ?, ?, now())
		on duplicate key update student_id = student_id
	`

	_, err := r.db.ExecContext(ctx, query,
		params.UserID,
		params.TestID,
		params.Group,
	)

	if err != nil {
		r.log.Error("failed to start test",
			zap.Error(err),
			zap.Int("user_id", params.UserID),
			zap.Int("test_id", params.TestID),
		)
		return false, err
	}

	return true, nil
}

// ==================== FINISH TEST ====================
func (r *repository) FinishTest(ctx context.Context, params FinishTestParams) (bool, error) {
	query := `
		update student_tests
		set mark = ?,
			success_rate = ?,
			date_end = now()
		where student_id = ? and test_id = ?
	`

	res, err := r.db.ExecContext(ctx, query,
		params.Result.Mark,
		params.Result.SuccessRate,
		params.UserID,
		params.TestID,
	)

	if err != nil {
		r.log.Error("failed to finish test",
			zap.Error(err),
			zap.Int("user_id", params.UserID),
			zap.Int("test_id", params.TestID),
		)
		return false, ErrFinishFailed
	}

	rows, err := res.RowsAffected()
	if err != nil {
		r.log.Error("failed to get rows affected",
			zap.Error(err),
		)
		return false, err
	}

	if rows == 0 {
		return false, ErrStudentTestNotFound
	}

	return true, nil
}

// ==================== GET ACTIVE TESTS ====================
// тесты доступны по группе, но НЕ завершены (mark IS NULL)
func (r *repository) GetActiveTests(ctx context.Context, userID int) ([]StudentActiveTestInfo, error) {
	query := `
		select
			t.id,
			t.name,
			u.name as lecturer_name,
			(select count(*) from tasks where test_id = t.id),
			(select count(*) from tasks where test_id = t.id and is_hard = true),
			st.date_start

		from users us

		join test_permissions tp 
			on tp.group = us.` + "`group`" + `

		join tests t 
			on t.id = tp.test_id

		join users u 
			on u.id = t.owner_id

		left join student_tests st
			on st.test_id = t.id
			and st.student_id = us.id

		where us.id = ?
		  and (st.mark is null)

		order by st.date_start desc
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.log.Error("failed to get active tests",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []StudentActiveTestInfo

	for rows.Next() {
		var item StudentActiveTestInfo

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.LecturerName,
			&item.CntQuestions,
			&item.CntHardQuestions,
			&item.DateStart,
		)
		if err != nil {
			r.log.Error("scan failed in GetActiveTests",
				zap.Error(err),
			)
			return nil, err
		}

		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		r.log.Error("rows iteration error",
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}

// ==================== GET FINISHED TESTS ====================
// только завершённые (есть оценка)
func (r *repository) GetFinishedTests(ctx context.Context, userID int) ([]StudentFinishedTestInfo, error) {
	query := `
		select
			t.id,
			t.name,
			u.name as lecturer_name,
			(select count(*) from tasks where test_id = t.id),
			(select count(*) from tasks where test_id = t.id and is_hard = true),
			st.mark,
			st.success_rate,
			st.date_start,
			st.date_end

		from student_tests st

		join tests t 
			on t.id = st.test_id

		join users u 
			on u.id = t.owner_id

		where st.student_id = ?
		  and st.mark is not null

		order by st.date_end desc
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.log.Error("failed to get finished tests",
			zap.Error(err),
			zap.Int("user_id", userID),
		)
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []StudentFinishedTestInfo

	for rows.Next() {
		var item StudentFinishedTestInfo

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.LecturerName,
			&item.CntQuestions,
			&item.CntHardQuestions,
			&item.Mark,
			&item.SuccessRate,
			&item.DateStart,
			&item.DateEnd,
		)
		if err != nil {
			r.log.Error("scan failed in GetFinishedTests",
				zap.Error(err),
			)
			return nil, err
		}

		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		r.log.Error("rows iteration error",
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}
