package task

import (
	"context"
	"database/sql"
	"testum-engine/app/internal/adapter/db"

	"go.uber.org/zap"
)

// ================= INTERFACE =================

type Repository interface {
	GetHardTasks(ctx context.Context, testID int) ([]Task, error)
	GetBaseTasks(ctx context.Context, testID int) ([]Task, error)
}

// ================= IMPLEMENTATION =================

type repository struct {
	db  db.Executor
	log *zap.Logger
}

func NewRepository(db db.Executor, log *zap.Logger) Repository {
	return &repository{db: db, log: log}
}

// ================= SCAN =================

func scanTasks(rows *sql.Rows) ([]Task, error) {
	taskMap := make(map[int]*Task)

	for rows.Next() {
		var (
			taskID   int
			text     string
			imageURL *string
			isHard   int // SQLite: 0/1

			answerID    int
			answerText  string
			answerImage *string
			isCorrect   int // SQLite: 0/1
		)

		if err := rows.Scan(
			&taskID,
			&text,
			&imageURL,
			&isHard,
			&answerID,
			&answerText,
			&answerImage,
			&isCorrect,
		); err != nil {
			return nil, err
		}

		task, exists := taskMap[taskID]
		if !exists {
			task = &Task{
				ID:       taskID,
				Text:     text,
				ImageURL: imageURL,
				IsHard:   isHard == 1,
				Answers:  make([]Answer, 0),
			}
			taskMap[taskID] = task
		}

		task.Answers = append(task.Answers, Answer{
			ID:        answerID,
			Text:      answerText,
			ImageURL:  answerImage,
			IsCorrect: isCorrect == 1,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make([]Task, 0, len(taskMap))
	for _, t := range taskMap {
		result = append(result, *t)
	}

	return result, nil
}

// ================= QUERY =================

func (r *repository) fetchTasks(ctx context.Context, testID int, isHard bool) ([]Task, error) {
	query := `
		SELECT 
			t.id,
			t.text,
			t.image_url,
			t.is_hard,
			a.id,
			a.text,
			a.image_url,
			a.is_correct
		FROM tasks t
		LEFT JOIN answers a ON a.task_id = t.id
		WHERE t.test_id = ?
		  AND t.is_hard = ?
		ORDER BY t.id
	`

	rows, err := r.db.QueryContext(ctx, query, testID, isHard)
	if err != nil {
		r.log.Error("fetchTasks query failed",
			zap.Error(err),
			zap.Int("test_id", testID),
			zap.Bool("is_hard", isHard),
		)
		return nil, err
	}
	defer rows.Close()

	tasks, err := scanTasks(rows)
	if err != nil {
		r.log.Error("fetchTasks scan failed", zap.Error(err))
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, ErrTaskNotFound
	}

	return tasks, nil
}

// ================= PUBLIC =================

func (r *repository) GetHardTasks(ctx context.Context, testID int) ([]Task, error) {
	return r.fetchTasks(ctx, testID, true)
}

func (r *repository) GetBaseTasks(ctx context.Context, testID int) ([]Task, error) {
	return r.fetchTasks(ctx, testID, false)
}
