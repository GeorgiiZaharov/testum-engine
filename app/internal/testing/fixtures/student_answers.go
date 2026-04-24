package fixtures

import (
	"context"
	"fmt"
	"time"
)

type studentAnswerSeed struct {
	ID        int
	StudentID int
	AnswerID  int
	CreatedAt time.Time
}

func (m *Manager) seedStudentAnswers(ctx context.Context) error {
	now := time.Now()

	answers := []studentAnswerSeed{
		{
			ID:        1,
			StudentID: 3,
			AnswerID:  2,
			CreatedAt: now,
		},

		{
			ID:        2,
			StudentID: 3,
			AnswerID:  3,
			CreatedAt: now,
		},

		{
			ID:        3,
			StudentID: 4,
			AnswerID:  1,
			CreatedAt: now.AddDate(-1, 0, 0),
		},
		{
			ID:        4,
			StudentID: 2,
			AnswerID:  3,
			CreatedAt: now,
		},
	}

	query := `
		INSERT INTO student_answers
			(id, student_id, answer_id, date_created)
		VALUES (?, ?, ?, ?)
	`

	for _, a := range answers {
		_, err := m.db.ExecContext(ctx, query,
			a.ID,
			a.StudentID,
			a.AnswerID,
			a.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to seed student_answer id=%d student_id=%d: %w",
				a.ID, a.StudentID, err,
			)
		}
	}

	return nil
}
