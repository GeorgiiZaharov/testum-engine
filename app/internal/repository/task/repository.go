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
	return &repository{
		db:  db,
		log: log,
	}
}

// ================= INTERNAL MAPPER =================

func scanTasks(rows *sql.Rows) ([]Task, error) {
	taskMap := make(map[int]*Task)

	for rows.Next() {
		var (
			taskID   int
			text     string
			imageURL *string
			isHard   bool

			answerText  string
			answerImage *string
			isCorrect   bool
		)

		if err := rows.Scan(
			&taskID,
			&text,
			&imageURL,
			&isHard,
			&answerText,
			&answerImage,
			&isCorrect,
		); err != nil {
			return nil, err
		}

		task, exists := taskMap[taskID]
		if !exists {
			task = &Task{
				Text:     text,
				ImageURL: imageURL,
				IsHard:   isHard,
				Answers:  []Answer{},
			}
			taskMap[taskID] = task
		}

		task.Answers = append(task.Answers, Answer{
			Text:      answerText,
			ImageURL:  answerImage,
			IsCorrect: isCorrect,
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

// ================= BASE QUERY =================

func (r *repository) fetchTasks(ctx context.Context, testID int, isHard bool) ([]Task, error) {
	query := `
		SELECT 
			t.id,
			t.text,
			t.image_url,
			t.is_hard,

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
	defer func() { _ = rows.Close() }()

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

// ================= PUBLIC METHODS =================

func (r *repository) GetHardTasks(ctx context.Context, testID int) ([]Task, error) {
	return r.fetchTasks(ctx, testID, true)
}

func (r *repository) GetBaseTasks(ctx context.Context, testID int) ([]Task, error) {
	return r.fetchTasks(ctx, testID, false)
}
