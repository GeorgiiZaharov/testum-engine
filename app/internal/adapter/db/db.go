package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// ================= CONFIG =================

type DBOptions struct {
	Host string
	User string
	Pass string
	Name string
	Port string
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
	dsn := opts.User + ":" + opts.Pass +
		"@tcp(" + opts.Host + ":" + opts.Port + ")/" + opts.Name +
		"?parseTime=true"

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &DB{db}, nil
}

// ================= TRANSACTIONS =================

// Начать транзакцию
func (db *DB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, nil)
}

// Helper для безопасной транзакции
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
