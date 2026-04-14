package main

import (
    "github.com/benmunyasya/dbbackuputility/internal/log"
    "github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
    Use:   "schedule",
    Short: "Schedule automatic backups",
    Long:  "Configure and manage automatic backup schedules using cron-like expressions or predefined intervals.",
    Run: func(cmd *cobra.Command, args []string) {
        interval, _ := cmd.Flags().GetString("interval")
        dbType, _ := cmd.Flags().GetString("db")

        log.Info("Scheduling backups for " + dbType + " database every " + interval + "...")
        
        // Placeholder: later we’ll call Scheduler service here
        log.Success("Schedule created successfully (placeholder).")
    },
}

func init() {
    scheduleCmd.Flags().String("db", "postgres", "Database type (postgres, mysql, mongo, sqlite)")
    scheduleCmd.Flags().String("interval", "daily", "Backup interval (e.g., hourly, daily, weekly, cron expression)")
}
