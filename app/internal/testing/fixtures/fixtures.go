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
// PUBLIC API (SQLite VERSION)
// =========================
func (m *Manager) Reset(ctx context.Context) error {
	// отключаем FK (SQLite way)
	if _, err := m.db.ExecContext(ctx, "PRAGMA foreign_keys = OFF"); err != nil {
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
		// SQLite doesn't support TRUNCATE
		if _, err := m.db.ExecContext(ctx, "DELETE FROM "+table); err != nil {
			return fmt.Errorf("delete from %s failed: %w", table, err)
		}
	}

	// reset AUTOINCREMENT counters (optional but useful)
	if _, err := m.db.ExecContext(ctx, "DELETE FROM sqlite_sequence"); err != nil {
		// не критично — просто логически можно игнорировать
	}

	// включаем FK обратно
	if _, err := m.db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		return err
	}

	return nil
}

// =========================
// SEED (без изменений)
// =========================
func (m *Manager) SeedAll(ctx context.Context) error {
	if err := m.seedUsers(ctx); err != nil {
		return err
	}
	if err := m.seedTests(ctx); err != nil {
		return err
	}
	if err := m.seedTasks(ctx); err != nil {
		return err
	}
	if err := m.seedAnswers(ctx); err != nil {
		return err
	}
	if err := m.seedTestPermissions(ctx); err != nil {
		return err
	}
	if err := m.seedStudentTests(ctx); err != nil {
		return err
	}
	if err := m.seedStudentAnswers(ctx); err != nil {
		return err
	}
	return nil
}
