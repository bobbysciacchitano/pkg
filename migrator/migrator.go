package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
)

// Migration represents a single SQL migration
type Migration struct {
	Filename string
	Content  string
}

// Run applies the provided migrations to the database
func Run(ctx context.Context, db *sql.DB, migrations []Migration) error {
	log.Println("📦 ensuring schema_migrations table exists")

	if err := ensureMigrationTable(ctx, db); err != nil {
		return err
	}

	applied, err := getAppliedMigrations(ctx, db)
	if err != nil {
		return err
	}

	// Sort by version prefix (e.g. 001_)
	sort.Slice(migrations, func(i, j int) bool {
		vi, _ := migrations[i].getVersion()
		vj, _ := migrations[j].getVersion()
		return vi < vj
	})

	for _, m := range migrations {
		if applied[m.Filename] {
			log.Printf("⏭️  skipping (already applied): %s", m.Filename)
			continue
		}

		log.Printf("⚙️  applying migration: %s", m.Filename)

		if err := m.Execute(ctx, db); err != nil {
			return fmt.Errorf("❌ failed migration %s: %w", m.Filename, err)
		}

		_, err := db.ExecContext(ctx, `
			INSERT INTO schema_migrations (name, applied_at)
			VALUES ($1, CURRENT_TIMESTAMP)`,
			m.Filename,
		)

		if err != nil {
			return fmt.Errorf("❌ failed to record migration %s: %w", m.Filename, err)
		}
	}

	log.Println("✅ migrations completed successfully")

	return nil
}

// Execute applies the migration to the DB
func (m Migration) Execute(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, m.Content)

	return err
}

// getVersion extracts the numeric prefix from the filename (e.g. 001_)
func (m Migration) getVersion() (int, error) {
	pattern := regexp.MustCompile(`^(\d+)_`)
	matches := pattern.FindStringSubmatch(m.Filename)

	if len(matches) != 2 {
		return 0, fmt.Errorf("invalid migration filename format: %s", m.Filename)
	}

	return strconv.Atoi(matches[1])
}

func ensureMigrationTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`)

	return err
}

func getAppliedMigrations(ctx context.Context, db *sql.DB) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, "SELECT name FROM schema_migrations")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	applied := make(map[string]bool)

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}

	return applied, nil
}
