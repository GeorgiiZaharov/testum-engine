package fixtures

import (
	"context"
	"fmt"
)

type answerSeed struct {
	ID        int
	TaskID    int
	Text      string
	ImageURL  *string
	IsCorrect bool
}

func (m *Manager) seedAnswers(ctx context.Context) error {
	answers := []answerSeed{
		// TASK 1
		{ID: 1, TaskID: 1, Text: "4", ImageURL: nil, IsCorrect: true},
		{ID: 2, TaskID: 1, Text: "5", ImageURL: nil, IsCorrect: false},

		// TASK 2
		{ID: 3, TaskID: 2, Text: "2x", ImageURL: ptr("img/answer_math_hard_1.png"), IsCorrect: true},
		{ID: 4, TaskID: 2, Text: "x^2", ImageURL: ptr("img/answer_math_hard_2.png"), IsCorrect: false},

		// TASK 3
		{ID: 5, TaskID: 3, Text: "Force = mass * acceleration", ImageURL: nil, IsCorrect: true},
		{ID: 6, TaskID: 3, Text: "Force = mass / acceleration", ImageURL: nil, IsCorrect: false},

		// TASK 4
		{ID: 7, TaskID: 4, Text: "E = (mv^2)/2 derived from work integral", ImageURL: ptr("img/answer_physics_hard_1.png"), IsCorrect: true},
		{ID: 8, TaskID: 4, Text: "E = (mv^2)/2 derived from work integral", ImageURL: ptr("img/answer_physics_hard_2.png"), IsCorrect: true},
	}

	query := `
		INSERT INTO answers
			(id, task_id, text, image_url, is_correct)
		VALUES (?, ?, ?, ?, ?)
	`

	for _, a := range answers {
		var img any
		if a.ImageURL != nil {
			img = *a.ImageURL
		} else {
			img = nil
		}

		isCorrect := 0
		if a.IsCorrect {
			isCorrect = 1
		}

		_, err := m.db.ExecContext(ctx, query,
			a.ID,
			a.TaskID,
			a.Text,
			img,
			isCorrect,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to seed answer id=%d task_id=%d: %w",
				a.ID, a.TaskID, err,
			)
		}
	}

	return nil
}
