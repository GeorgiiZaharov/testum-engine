package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// ================= CONFIG =================

type DBOptions struct {
	Path string // путь к sqlite файлу
}

// ================= DB =================

type DB struct {
	*sqlx.DB
}

// ================= EXECUTOR =================

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// ================= INIT =================

func NewDB(opts DBOptions) (*DB, error) {
	db, err := sqlx.Connect("sqlite", opts.Path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	return &DB{db}, nil
}

// ================= TRANSACTIONS =================

func (db *DB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, nil)
}

func (db *DB) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
