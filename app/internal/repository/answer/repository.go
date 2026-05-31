package answer

import (
	"context"

	"go.uber.org/zap"

	"testum-engine/app/internal/adapter/db"
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

func (r *repository) getAnswers(
	ctx context.Context,
	userID,
	testID int,
	isHard bool,
) ([]TaskAnswer, error) {
	query := `
		SELECT a.task_id, sa.answer_id
		FROM student_answers sa
		JOIN answers a ON a.id = sa.answer_id
		JOIN tasks t ON t.id = a.task_id
		WHERE sa.student_id = ?
		  AND t.test_id = ?
		  AND t.is_hard = ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, testID, isHard)
	if err != nil {
		r.log.Error("get answers failed", zap.Error(err))
		return nil, ErrGetAnswers
	}
	defer func() { _ = rows.Close() }()

	resultMap := make(map[int][]int)

	for rows.Next() {
		var taskID int
		var answerID int

		if err := rows.Scan(&taskID, &answerID); err != nil {
			return nil, err
		}

		resultMap[taskID] = append(resultMap[taskID], answerID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]TaskAnswer, 0, len(resultMap))

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

func (r *repository) GetHardAnswersByTest(
	ctx context.Context,
	testID int,
) ([]TaskAnswer, error) {
	return r.getCorrectAnswers(ctx, testID, true)
}

func (r *repository) GetBaseAnswersByTest(
	ctx context.Context,
	testID int,
) ([]TaskAnswer, error) {
	return r.getCorrectAnswers(ctx, testID, false)
}

func (r *repository) getCorrectAnswers(
	ctx context.Context,
	testID int,
	isHard bool,
) ([]TaskAnswer, error) {
	query := `
		SELECT t.id, a.id
		FROM tasks t
		JOIN answers a ON a.task_id = t.id
		WHERE t.test_id = ?
		  AND t.is_hard = ?
		  AND a.is_correct = 1
	`

	rows, err := r.db.QueryContext(ctx, query, testID, isHard)
	if err != nil {
		r.log.Error("get correct answers failed", zap.Error(err))
		return nil, ErrGetAnswers
	}
	defer func() { _ = rows.Close() }()

	resultMap := make(map[int][]int)

	for rows.Next() {
		var taskID int
		var answerID int

		if err := rows.Scan(&taskID, &answerID); err != nil {
			return nil, err
		}

		resultMap[taskID] = append(resultMap[taskID], answerID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]TaskAnswer, 0, len(resultMap))

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

func (r *repository) SaveHardAnswers(
	ctx context.Context,
	userID int,
	answers []int,
) (bool, error) {
	return r.saveAnswers(ctx, userID, answers)
}

func (r *repository) SaveBaseAnswers(
	ctx context.Context,
	userID int,
	answers []int,
) (bool, error) {
	return r.saveAnswers(ctx, userID, answers)
}

func (r *repository) saveAnswers(
	ctx context.Context,
	userID int,
	answers []int,
) (bool, error) {
	if len(answers) == 0 {
		return true, nil
	}

	query := `
		INSERT INTO student_answers (
			student_id,
			answer_id,
			date_created
		)
		VALUES
	`

	args := make([]any, 0, len(answers)*2)

	values := ""

	validCount := 0

	for _, answerID := range answers {
		if answerID <= 0 {
			continue
		}

		if validCount > 0 {
			values += ","
		}

		values += "(?, ?, CURRENT_TIMESTAMP)"

		args = append(args, userID, answerID)

		validCount++
	}

	if validCount == 0 {
		return true, nil
	}

	query += values

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.log.Error("batch save answers failed", zap.Error(err))
		return false, ErrSaveFailed
	}

	return true, nil
}

// =========================
// DELETE ATTEMPT
// =========================

func (r *repository) DeleteAttempt(
	ctx context.Context,
	userID,
	testID int,
) (bool, error) {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM student_answers
		WHERE student_id = ?
		  AND answer_id IN (
			  SELECT a.id
			  FROM answers a
			  JOIN tasks t ON t.id = a.task_id
			  WHERE t.test_id = ?
		  )
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
