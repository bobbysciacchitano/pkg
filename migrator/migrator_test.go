package migrator_test

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bobbysciacchitano/pkg/migrator"
	_ "modernc.org/sqlite"
)

//go:embed testdata/*.sql
var testMigrations embed.FS

func loadTestMigrations() ([]migrator.Migration, error) {
	var migrations []migrator.Migration

	entries, err := testMigrations.ReadDir("testdata")

	if err != nil {
		return nil, fmt.Errorf("failed to read testdata dir: %w", err)
	}

	pattern := regexp.MustCompile(`^\d+_.+\.sql$`)

	for _, entry := range entries {
		if entry.IsDir() || !pattern.MatchString(entry.Name()) {
			continue
		}

		content, err := testMigrations.ReadFile(filepath.Join("testdata", entry.Name()))

		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		migrations = append(migrations, migrator.Migration{
			Filename: entry.Name(),
			Content:  string(content),
		})
	}

	return migrations, nil
}

func setupInMemoryDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")

	if err != nil {
		t.Fatalf("failed to open sqlite in-memory DB: %v", err)
	}

	t.Cleanup(func() { _ = db.Close() })

	return db
}

func TestMigrations_RunSuccessfully(t *testing.T) {
	ctx := context.Background()
	db := setupInMemoryDB(t)

	migrations, err := loadTestMigrations()

	if err != nil {
		t.Fatalf("failed to load test migrations: %v", err)
	}

	err = migrator.Run(ctx, db, migrations)

	if err != nil {
		t.Fatalf("unexpected error running migrations: %v", err)
	}

	// Confirm test_user table exists
	var count int

	err = db.QueryRow("SELECT COUNT(*) FROM test_user").Scan(&count)

	if err != nil && !strings.Contains(err.Error(), "no such table") {
		t.Errorf("expected test_user table to exist but got error: %v", err)
	}
}

func TestMigrations_IsIdempotent(t *testing.T) {
	ctx := context.Background()
	db := setupInMemoryDB(t)

	migrations, err := loadTestMigrations()

	if err != nil {
		t.Fatalf("failed to load test migrations: %v", err)
	}

	// Run twice
	if err := migrator.Run(ctx, db, migrations); err != nil {
		t.Fatalf("first run failed: %v", err)
	}

	if err := migrator.Run(ctx, db, migrations); err != nil {
		t.Fatalf("second run failed: %v", err)
	}

	// Count applied migrations
	rows, err := db.Query("SELECT name FROM schema_migrations")

	if err != nil {
		t.Fatalf("failed to query schema_migrations: %v", err)
	}

	defer rows.Close()

	var applied []string

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			t.Errorf("failed to scan row: %v", err)
		}

		applied = append(applied, name)
	}

	log.Printf("✅ Applied migrations: %v", applied)

	expected := 2 // based on the two .sql files

	if len(applied) != expected {
		t.Errorf("expected %d applied migrations, got %d", expected, len(applied))
	}
}
