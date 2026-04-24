package fixtures

import (
	"context"
	"fmt"
	"time"
)

type studentTestSeed struct {
	ID        int
	StudentID int
	TestID    int
	Mark      *int
	Group     string
	Success   *float64
	DateStart time.Time
	DateEnd   *time.Time
}

func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(f int) *int {
	return &f
}

func (m *Manager) seedStudentTests(ctx context.Context) error {
	now := time.Now()
	yearAgo := now.AddDate(-1, 0, 0)
	yearAndHourAgo := now.AddDate(-1, 0, 0)

	finished := now.Add(-24 * time.Hour)

	tests := []studentTestSeed{
		{
			ID:        1,
			StudentID: 2,
			TestID:    1,
			Mark:      nil,
			Group:     "A-101",
			Success:   nil,
			DateStart: now,
			DateEnd:   nil,
		},
		{
			ID:        2,
			StudentID: 4,
			TestID:    1,
			Mark:      intPtr(5),
			Group:     "B-202",
			Success:   floatPtr(100),
			DateStart: yearAndHourAgo,
			DateEnd:   &yearAgo,
		},
		{
			ID:        3,
			StudentID: 3,
			TestID:    1,
			Mark:      intPtr(4),
			Group:     "A-101",
			Success:   floatPtr(80),
			DateStart: finished,
			DateEnd:   &now,
		},
	}

	query := `
    INSERT INTO student_tests
        (id, student_id, test_id, mark, ` + "`group`" + `, success_rate, date_start, date_end)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, st := range tests {
		_, err := m.db.ExecContext(ctx, query,
			st.ID,
			st.StudentID,
			st.TestID,
			st.Mark,
			st.Group,
			st.Success,
			st.DateStart,
			st.DateEnd,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to seed student_test id=%d student=%d test=%d: %w",
				st.ID, st.StudentID, st.TestID, err,
			)
		}
	}

	return nil
}
