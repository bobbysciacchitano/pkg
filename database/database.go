package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func NewDatabase(ctx context.Context, path string) (*sql.DB, error) {
	instance, err := sql.Open("sqlite", path)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err = instance.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if _, err = instance.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	if err = instance.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return instance, nil
}
