package migrations

import (
    "database/sql"
    "fmt"

    _ "github.com/mattn/go-sqlite3"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/sqlite3"
)

func RunMigrations(dbPath string, migrationsPath string) error {
    sqlite, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return fmt.Errorf("failed to open sqlite db: %w", err)
    }

    driver, err := sqlite3.WithInstance(sqlite, &sqlite3.Config{})
    if err != nil {
        return fmt.Errorf("failed to init sqlite driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://" + migrationsPath,
        "sqlite3", driver,
    )
    if err != nil {
        return fmt.Errorf("failed to init migrate: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migration failed: %w", err)
    }

    return nil
}
