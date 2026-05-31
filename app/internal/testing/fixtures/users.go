package fixtures

import (
	"context"
	"fmt"
	"time"
)

type userSeed struct {
	ID           int
	Login        string
	Mail         string
	Name         string
	Group        *string
	IsLecturer   bool
	DateCreated  time.Time
	DateModified time.Time
}

func ptr(s string) *string {
	return &s
}

func (m *Manager) seedUsers(ctx context.Context) error {
	now := time.Now()
	yearAgo := now.AddDate(-1, 0, 0)

	users := []userSeed{
		{
			ID:           1,
			Login:        "student1",
			Mail:         "student1@mail.com",
			Name:         "Student One",
			Group:        ptr("A-101"),
			IsLecturer:   false,
			DateCreated:  now,
			DateModified: now,
		},
		{
			ID:           2,
			Login:        "student2",
			Mail:         "student2@mail.com",
			Name:         "Student Two",
			Group:        ptr("A-101"),
			IsLecturer:   false,
			DateCreated:  now,
			DateModified: now,
		},
		{
			ID:           3,
			Login:        "student3",
			Mail:         "student3@mail.com",
			Name:         "Student Three",
			Group:        ptr("B-202"),
			IsLecturer:   false,
			DateCreated:  now,
			DateModified: now,
		},
		{
			ID:           4,
			Login:        "old_student",
			Mail:         "old_student@mail.com",
			Name:         "Old Student",
			Group:        ptr("A-101"),
			IsLecturer:   false,
			DateCreated:  yearAgo,
			DateModified: yearAgo,
		},
		{
			ID:           5,
			Login:        "magistr",
			Mail:         "magistr@mail.com",
			Name:         "Magistr Magistr",
			Group:        ptr("C-000"),
			IsLecturer:   true,
			DateCreated:  now,
			DateModified: now,
		},
		{
			ID:           6,
			Login:        "lecturer1",
			Mail:         "lecturer1@mail.com",
			Name:         "Lecturer One",
			Group:        nil,
			IsLecturer:   true,
			DateCreated:  now,
			DateModified: now,
		},
		{
			ID:           7,
			Login:        "lecturer2",
			Mail:         "lecturer2@mail.com",
			Name:         "Lecturer Two",
			Group:        nil,
			IsLecturer:   true,
			DateCreated:  now,
			DateModified: now,
		},
		{
			ID:           8,
			Login:        "olbgvl",
			Mail:         "admin@mail.com",
			Name:         "Admin Admin",
			Group:        nil,
			IsLecturer:   true,
			DateCreated:  now,
			DateModified: now,
		},
	}

	query := `
		INSERT INTO users
			(id, login, mail, name, ` + "`group`" + `, is_lecturer, date_created, date_modified)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	for _, u := range users {
		_, err := m.db.ExecContext(ctx, query,
			u.ID,
			u.Login,
			u.Mail,
			u.Name,
			u.Group,
			boolToInt(u.IsLecturer), // SQLite-safe
			u.DateCreated,
			u.DateModified,
		)
		if err != nil {
			return fmt.Errorf("failed to seed user %s (id=%d): %w", u.Login, u.ID, err)
		}
	}

	return nil
}
