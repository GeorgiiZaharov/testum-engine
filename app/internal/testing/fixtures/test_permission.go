package fixtures

import (
	"context"
	"fmt"
)

type testPermissionSeed struct {
	ID     int
	TestID int
	Group  string
}

func (m *Manager) seedTestPermissions(ctx context.Context) error {

	permissions := []testPermissionSeed{
		{
			ID:     1,
			TestID: 1,
			Group:  "A-101",
		},
		{
			ID:     2,
			TestID: 1,
			Group:  "B-202",
		},
		{
			ID:     3,
			TestID: 2,
			Group:  "A-101",
		},
		{
			ID:     4,
			TestID: 3,
			Group:  "C-000",
		},
	}

	query := `
    INSERT INTO test_permissions
        (id, test_id, ` + "`group`" + `)
    VALUES (?, ?, ?)
	`

	for _, p := range permissions {
		_, err := m.db.ExecContext(ctx, query,
			p.ID,
			p.TestID,
			p.Group,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to seed test_permission id=%d test_id=%d group=%s: %w",
				p.ID, p.TestID, p.Group, err,
			)
		}
	}

	return nil
}
