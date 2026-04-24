package fixtures

import (
	"context"
	"fmt"
	"testum-engine/app/internal/adapter/db"

	"github.com/jmoiron/sqlx"
)

type Manager struct {
	db *sqlx.DB
}

func New(db *db.DB) *Manager {
	return &Manager{db: db.DB}
}

// =========================
// PUBLIC API
// =========================
func (m *Manager) Reset(ctx context.Context) error {
	// отключаем FK проверки
	if _, err := m.db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		return err
	}

	tables := []string{
		"student_answers",
		"student_tests",
		"test_permissions",
		"answers",
		"tasks",
		"tests",
		"users",
	}

	for _, table := range tables {
		if _, err := m.db.ExecContext(ctx, "TRUNCATE TABLE "+table); err != nil {
			return fmt.Errorf("truncate %s failed: %w", table, err)
		}
	}

	if _, err := m.db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		return err
	}

	return nil
}

func (m *Manager) SeedAll(ctx context.Context) error {
	// 1. USERS
	if err := m.seedUsers(ctx); err != nil {
		return err
	}

	// 2. TESTS
	if err := m.seedTests(ctx); err != nil {
		return err
	}

	// 3. TASKS
	if err := m.seedTasks(ctx); err != nil {
		return err
	}

	// 4. ANSWERS
	if err := m.seedAnswers(ctx); err != nil {
		return err
	}

	// 5. TEST PERMISSIONS
	if err := m.seedTestPermissions(ctx); err != nil {
		return err
	}

	// 6. STUDENT TESTS
	if err := m.seedStudentTests(ctx); err != nil {
		return err
	}

	// 7. STUDENT ANSWERS
	if err := m.seedStudentAnswers(ctx); err != nil {
		return err
	}

	return nil
}
