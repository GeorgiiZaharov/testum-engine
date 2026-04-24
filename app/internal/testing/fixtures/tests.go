package fixtures

import (
	"context"
	"fmt"
	"time"
)

type testSeed struct {
	ID         int
	LecturerID int
	Name       string
	FileName   string
	CreatedAt  time.Time
}

func (m *Manager) seedTests(ctx context.Context) error {
	now := time.Now()

	tests := []testSeed{
		{
			ID:         1,
			LecturerID: 5,
			Name:       "Math Basics",
			FileName:   "math_basics.json",
			CreatedAt:  now,
		},
		{
			ID:         2,
			LecturerID: 6,
			Name:       "Linear Algebra",
			FileName:   "linear_algebra.json",
			CreatedAt:  now,
		},
		{
			ID:         3,
			LecturerID: 6,
			Name:       "Physics Intro",
			FileName:   "physics_intro.json",
			CreatedAt:  now,
		},
		{
			ID:         4,
			LecturerID: 7,
			Name:       "Programming Basics",
			FileName:   "programming_basics.json",
			CreatedAt:  now,
		},
	}

	query := `
		INSERT INTO tests
			(id, owner_id, name, file_name, date_created)
		VALUES (?, ?, ?, ?, ?)
	`

	for _, t := range tests {
		_, err := m.db.ExecContext(ctx, query,
			t.ID,
			t.LecturerID,
			t.Name,
			t.FileName,
			t.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to seed test id=%d name=%s: %w",
				t.ID, t.Name, err,
			)
		}
	}

	return nil
}
