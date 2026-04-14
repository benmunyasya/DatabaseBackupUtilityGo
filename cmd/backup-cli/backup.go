// cmd/backup-cli/backup.go
package main

import (
    "github.com/benmunyasya/dbbackuputility/internal/backup"
    "github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
    Use:   "backup",
    Short: "Run a database backup",
    Long:  "Perform a backup operation for the specified database type (Postgres, MySQL, MongoDB, SQLite).",
    Run: func(cmd *cobra.Command, args []string) {
        dbType, _ := cmd.Flags().GetString("db")
        backupType, _ := cmd.Flags().GetString("type")
        dbName, _ := cmd.Flags().GetString("db-name")

        backup.RunBackup(dbType, dbName, backupType)
    },
}

func init() {
    backupCmd.Flags().String("db", "postgres", "Database type (postgres, mysql, mongo, sqlite)")
    backupCmd.Flags().String("type", "full", "Backup type (full, incremental, differential)")
    backupCmd.Flags().String("db-name", "", "Database name to back up")
    backupCmd.MarkFlagRequired("db-name")
}
