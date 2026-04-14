// internal/scheduler/manager.go
package scheduler

import (
    "database/sql"
    "fmt"

    "github.com/robfig/cron/v3"
    "github.com/benmunyasya/dbbackuputility/internal/log"
    "github.com/benmunyasya/dbbackuputility/internal/backup"
)

type ScheduleManager struct {
    db   *sql.DB
    cron *cron.Cron
}

func NewScheduleManager(db *sql.DB) *ScheduleManager {
    return &ScheduleManager{
        db:   db,
        cron: cron.New(),
    }
}

// StartScheduler loads persisted schedules from DB and registers them with cron
func (sm *ScheduleManager) StartScheduler() error {
    rows, err := sm.db.Query(`SELECT id, db_type, db_name, backup_type, cron_expr FROM schedules`)
    if err != nil {
        return fmt.Errorf("failed to load schedules: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var id int
        var dbType, dbName, backupType, cronExpr string
        if err := rows.Scan(&id, &dbType, &dbName, &backupType, &cronExpr); err != nil {
            return err
        }

        // Register cron job
        _, err := sm.cron.AddFunc(cronExpr, func() {
            log.Info(fmt.Sprintf("Running scheduled %s backup for %s (%s)", backupType, dbName, dbType))
            backup.RunBackup(dbType, dbName, backupType)
        })
        if err != nil {
            log.Error("Failed to register cron job: " + err.Error())
        }
    }

    sm.cron.Start()
    log.Success("Scheduler started and jobs loaded.")
    return nil
}

// RegisterJob adds a single job to cron without reloading all schedules
func (sm *ScheduleManager) RegisterJob(dbType, dbName, backupType, cronExpr string) error {
    _, err := sm.cron.AddFunc(cronExpr, func() {
        log.Info(fmt.Sprintf("Running scheduled %s backup for %s (%s)", backupType, dbName, dbType))
        backup.RunBackup(dbType, dbName, backupType)
    })
    return err
}
