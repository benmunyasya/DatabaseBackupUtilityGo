// cmd/backup-cli/main.go
package main

import (
    "os"
    "os/signal"
    "syscall"
    "database/sql"

    "github.com/benmunyasya/dbbackuputility/internal/log"
    "github.com/benmunyasya/dbbackuputility/migrations"
    "github.com/benmunyasya/dbbackuputility/internal/scheduler"
    "github.com/spf13/cobra"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    var rootCmd = &cobra.Command{
        Use:   "backup-cli",
        Short: "Database Backup Utility",
        Long:  "A CLI tool to backup and restore multiple databases with compression, scheduling, and cloud storage support.",
    }

    // Add subcommands
    rootCmd.AddCommand(backupCmd)
    rootCmd.AddCommand(restoreCmd)
    rootCmd.AddCommand(scheduleCmd)
    rootCmd.AddCommand(configCmd)

    // Run migrations before executing commands
    dbPath := "schedules.db"
    migrationsPath := "migrations"

    if err := migrations.RunMigrations(dbPath, migrationsPath); err != nil {
        log.Error("Migration failed: " + err.Error())
        os.Exit(1)
    }

    // Initialize scheduler manager (so cron jobs can be loaded at startup)
    sqlite, _ := sql.Open("sqlite3", dbPath)
    sm := scheduler.NewScheduleManager(sqlite)
    _ = sm.StartScheduler() // load persisted jobs

    // Execute CLI commands
    if err := rootCmd.Execute(); err != nil {
        log.Error("Command execution failed: " + err.Error())
        os.Exit(1)
    }

    // Block until interrupt so cron jobs keep running
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
}
