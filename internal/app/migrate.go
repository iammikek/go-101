package app

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(databaseURL string) error {
	migrationsDir, err := findMigrationsDir()
	if err != nil {
		return err
	}

	m, err := migrate.New(
		"file://"+migrationsDir,
		fmt.Sprintf("sqlite3://%s", databaseURL),
	)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Join(fmt.Errorf("run migrations: %w", err), closeMigrator(m))
	}

	return closeMigrator(m)
}

func closeMigrator(m *migrate.Migrate) error {
	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		return fmt.Errorf("close migration source: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("close migration database: %w", dbErr)
	}
	return nil
}

func findMigrationsDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}

	for {
		candidate := filepath.Join(dir, "migrations")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("migrations directory not found")
		}
		dir = parent
	}
}
