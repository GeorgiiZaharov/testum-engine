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

func floatPtr(f float64) *float64 { return &f }
func intPtr(f int) *int           { return &f }

func (m *Manager) seedStudentTests(ctx context.Context) error {
	// SQLite-safe deterministic time (вместо time.Now)
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	yearAgo := base.AddDate(-1, 0, 0)
	finished := base.Add(-24 * time.Hour)

	tests := []studentTestSeed{
		{
			ID:        1,
			StudentID: 2,
			TestID:    1,
			Mark:      nil,
			Group:     "A-101",
			Success:   nil,
			DateStart: base,
			DateEnd:   nil,
		},
		{
			ID:        2,
			StudentID: 4,
			TestID:    1,
			Mark:      intPtr(5),
			Group:     "B-202",
			Success:   floatPtr(100),
			DateStart: yearAgo,
			DateEnd:   &base,
		},
		{
			ID:        3,
			StudentID: 3,
			TestID:    1,
			Mark:      intPtr(4),
			Group:     "A-101",
			Success:   floatPtr(80),
			DateStart: finished,
			DateEnd:   &base,
		},
		{
			ID:        4,
			StudentID: 1,
			TestID:    1,
			Mark:      nil,
			Group:     "A-101",
			Success:   nil,
			DateStart: base,
			DateEnd:   nil,
		},
	}

	query := `
	INSERT INTO student_tests
		(id, student_id, test_id, mark, ` + "`group`" + `, success_rate, date_start, date_end)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, st := range tests {
		// SQLite-safe conversions

		var mark any
		if st.Mark != nil {
			mark = *st.Mark
		} else {
			mark = nil
		}

		var success any
		if st.Success != nil {
			success = *st.Success
		} else {
			success = nil
		}

		var dateEnd any
		if st.DateEnd != nil {
			dateEnd = st.DateEnd.UTC().Format(time.RFC3339)
		} else {
			dateEnd = nil
		}

		_, err := m.db.ExecContext(ctx, query,
			st.ID,
			st.StudentID,
			st.TestID,
			mark,
			st.Group,
			success,
			st.DateStart.UTC().Format(time.RFC3339),
			dateEnd,
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
