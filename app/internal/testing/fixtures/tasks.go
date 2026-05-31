package fixtures

import (
	"context"
	"fmt"
)

type taskSeed struct {
	ID       int
	TestID   int
	Text     string
	ImageURL *string
	IsHard   bool
}

func (m *Manager) seedTasks(ctx context.Context) error {
	tasks := []taskSeed{
		{
			ID:       1,
			TestID:   1,
			Text:     "Solve equation: 2 + 2 = ?",
			ImageURL: ptr("img/test1_easy.png"),
			IsHard:   false,
		},
		{
			ID:       2,
			TestID:   1,
			Text:     "Prove derivative of x^2",
			ImageURL: nil,
			IsHard:   true,
		},
		{
			ID:       3,
			TestID:   3,
			Text:     "What is Newton's second law?",
			ImageURL: ptr("img/test3_easy.png"),
			IsHard:   false,
		},
		{
			ID:       4,
			TestID:   3,
			Text:     "Derive kinetic energy formula from work theorem",
			ImageURL: ptr("img/test3_hard.png"),
			IsHard:   true,
		},
	}

	query := `
		INSERT INTO tasks
			(id, test_id, text, image_url, is_hard)
		VALUES (?, ?, ?, ?, ?)
	`

	for _, t := range tasks {
		_, err := m.db.ExecContext(ctx, query,
			t.ID,
			t.TestID,
			t.Text,
			t.ImageURL,
			boolToInt(t.IsHard), // SQLite-safe bool
		)
		if err != nil {
			return fmt.Errorf(
				"failed to seed task id=%d test_id=%d: %w",
				t.ID, t.TestID, err,
			)
		}
	}

	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
