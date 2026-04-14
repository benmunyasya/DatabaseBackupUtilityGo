// cmd/backup-cli/schedule.go
package main

import (
    "database/sql"
    "fmt"

    "github.com/benmunyasya/dbbackuputility/internal/log"
    "github.com/benmunyasya/dbbackuputility/internal/scheduler"
    "github.com/spf13/cobra"
    _ "github.com/mattn/go-sqlite3"
)

var scheduleCmd = &cobra.Command{
    Use:   "schedule",
    Short: "Schedule automatic backups",
    Long:  "Configure and manage automatic backup schedules using cron-like expressions.",
    Run: func(cmd *cobra.Command, args []string) {
        dbType, _ := cmd.Flags().GetString("db")
        dbName, _ := cmd.Flags().GetString("db-name")
        backupType, _ := cmd.Flags().GetString("backup")
        cronExpr, _ := cmd.Flags().GetString("cron")

        // Open SQLite DB
        sqlite, err := sql.Open("sqlite3", "schedules.db")
        if err != nil {
            log.Error("Failed to open schedules DB: " + err.Error())
            return
        }
        defer sqlite.Close()

        // Insert schedule into DB
        _, err = sqlite.Exec(
            `INSERT INTO schedules (db_type, db_name, backup_type, cron_expr) VALUES (?, ?, ?, ?)`,
            dbType, dbName, backupType, cronExpr,
        )
        if err != nil {
            log.Error("Failed to insert schedule: " + err.Error())
            return
        }

        log.Info(fmt.Sprintf("Created schedule for %s (%s) backups [%s] with cron '%s'",
            dbName, dbType, backupType, cronExpr))

        // Register with scheduler immediately
        sm := scheduler.NewScheduleManager(sqlite)
        if err := sm.RegisterJob(dbType, dbName, backupType, cronExpr); err != nil {
            log.Error("Failed to register cron job: " + err.Error())
            return
        }

        log.Success("Schedule created and registered successfully.")
    },
}

func init() {
    scheduleCmd.Flags().String("db", "postgres", "Database type (postgres, mysql, mongo, sqlite)")
    scheduleCmd.Flags().String("db-name", "", "Database name")
    scheduleCmd.Flags().String("backup", "full", "Backup type (full, incremental)")
    scheduleCmd.Flags().String("cron", "0 2 * * *", "Cron expression for schedule")
}
