package answer

import (
	"context"
	"testum-engine/app/internal/adapter/db"

	"go.uber.org/zap"
)

type Repository interface {
	GetHardAnswers(ctx context.Context, userID, testID int) ([]TaskAnswer, error)
	GetBaseAnswers(ctx context.Context, userID, testID int) ([]TaskAnswer, error)

	GetHardAnswersByTest(ctx context.Context, testID int) ([]TaskAnswer, error)
	GetBaseAnswersByTest(ctx context.Context, testID int) ([]TaskAnswer, error)

	SaveHardAnswers(ctx context.Context, userID int, answers []int) (bool, error)
	SaveBaseAnswers(ctx context.Context, userID int, answers []int) (bool, error)

	DeleteAttempt(ctx context.Context, userID, testID int) (bool, error)
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

// =========================
// GET STUDENT ANSWERS
// =========================

func (r *repository) GetHardAnswers(ctx context.Context, userID, testID int) ([]TaskAnswer, error) {
	return r.getAnswers(ctx, userID, testID, true)
}

func (r *repository) GetBaseAnswers(ctx context.Context, userID, testID int) ([]TaskAnswer, error) {
	return r.getAnswers(ctx, userID, testID, false)
}

func (r *repository) getAnswers(ctx context.Context, userID, testID int, isHard bool) ([]TaskAnswer, error) {
	query := `
		select a.task_id, sa.answer_id
		from student_answers sa
		join answers a on a.id = sa.answer_id
		join tasks t on t.id = a.task_id
		where sa.student_id = ? and t.test_id = ? and t.is_hard = ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, testID, isHard)
	if err != nil {
		r.log.Error("get answers failed", zap.Error(err))
		return nil, ErrGetAnswers
	}
	defer func() { _ = rows.Close() }()

	resultMap := make(map[int][]int)

	for rows.Next() {
		var taskID, answerID int

		if err := rows.Scan(&taskID, &answerID); err != nil {
			return nil, err
		}

		resultMap[taskID] = append(resultMap[taskID], answerID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var result []TaskAnswer
	for taskID, opts := range resultMap {
		result = append(result, TaskAnswer{
			TaskID:  taskID,
			Options: opts,
		})
	}

	return result, nil
}

// =========================
// GET CORRECT ANSWERS BY TEST
// =========================

func (r *repository) GetHardAnswersByTest(ctx context.Context, testID int) ([]TaskAnswer, error) {
	return r.getCorrectAnswers(ctx, testID, true)
}

func (r *repository) GetBaseAnswersByTest(ctx context.Context, testID int) ([]TaskAnswer, error) {
	return r.getCorrectAnswers(ctx, testID, false)
}

func (r *repository) getCorrectAnswers(ctx context.Context, testID int, isHard bool) ([]TaskAnswer, error) {
	query := `
		select t.id, a.id
		from tasks t
		join answers a on a.task_id = t.id
		where t.test_id = ? and t.is_hard = ? and a.is_correct = true
	`

	rows, err := r.db.QueryContext(ctx, query, testID, isHard)
	if err != nil {
		r.log.Error("get correct answers failed", zap.Error(err))
		return nil, ErrGetAnswers
	}
	defer func() { _ = rows.Close() }()

	resultMap := make(map[int][]int)

	for rows.Next() {
		var taskID, answerID int

		if err := rows.Scan(&taskID, &answerID); err != nil {
			return nil, err
		}

		resultMap[taskID] = append(resultMap[taskID], answerID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var result []TaskAnswer
	for taskID, opts := range resultMap {
		result = append(result, TaskAnswer{
			TaskID:  taskID,
			Options: opts,
		})
	}

	return result, nil
}

// =========================
// SAVE ANSWERS
// =========================

func (r *repository) SaveHardAnswers(ctx context.Context, userID int, answers []int) (bool, error) {
	return r.saveAnswers(ctx, userID, answers)
}

func (r *repository) SaveBaseAnswers(ctx context.Context, userID int, answers []int) (bool, error) {
	return r.saveAnswers(ctx, userID, answers)
}

func (r *repository) saveAnswers(ctx context.Context, userID int, answers []int) (bool, error) {
	if len(answers) == 0 {
		return true, nil
	}

	query := `
		INSERT INTO student_answers (student_id, answer_id, date_created)
		VALUES 
	`

	args := make([]any, 0, len(answers)*2)

	values := ""

	for i, answerID := range answers {
		if answerID <= 0 {
			continue
		}

		if i > 0 && len(values) > 0 {
			values += ","
		}

		values += "(?, ?, NOW())"
		args = append(args, userID, answerID)
	}

	query += values

	if len(args) == 0 {
		return true, nil
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.log.Error("batch save answers failed", zap.Error(err))
		return false, ErrSaveFailed
	}

	return true, nil
}

// func (r *repository) saveAnswers(ctx context.Context, userID int, answers []TaskAnswer) (bool, error) {
// 	for _, ans := range answers {
// 		payload, err := json.Marshal(ans.Options)
// 		if err != nil {
// 			r.log.Error("marshal answers failed", zap.Error(err))
// 			return false, ErrSaveFailed
// 		}
//
// 		_, err = r.db.ExecContext(ctx, `
// 			insert into student_answers (student_id, answer_id, date_created)
// 			values (?, ?, now())
// 		`, userID, string(payload))
//
// 		if err != nil {
// 			r.log.Error("save answers failed", zap.Error(err))
// 			return false, ErrSaveFailed
// 		}
// 	}
//
// 	return true, nil
// }

// =========================
// DELETE ATTEMPT
// =========================

func (r *repository) DeleteAttempt(ctx context.Context, userID, testID int) (bool, error) {
	res, err := r.db.ExecContext(ctx, `
		delete sa
		from student_answers sa
		join answers a on a.id = sa.answer_id
		join tasks t on t.id = a.task_id
		where sa.student_id = ? and t.test_id = ?
	`, userID, testID)

	if err != nil {
		r.log.Error("delete attempt failed", zap.Error(err))
		return false, ErrDeleteAttempt
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, ErrNotFound
	}

	return true, nil
}
