package main

import (
    "os"
    "strconv"
    "strings"

    "github.com/benmunyasya/dbbackuputility/internal/db"
    "github.com/benmunyasya/dbbackuputility/internal/log" //  import logging utility
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

        log.Info("Running " + backupType + " backup for " + dbType + " database (" + dbName + ")...")

        adapter, err := db.GetAdapter(dbType)
        if err != nil {
            log.Error("Adapter error: " + err.Error())
            return
        }

        prefix := strings.ToUpper(dbType)

        portStr := os.Getenv(prefix + "_DB_PORT")
        port := 5432 // default for Postgres
        if portStr != "" {
            p, err := strconv.Atoi(portStr)
            if err != nil {
                log.Error("Invalid port: " + err.Error())
                return
            }
            port = p
        }

        params := db.ConnectionParams{
            Host:     os.Getenv(prefix + "_DB_HOST"),
            Port:     port,
            User:     os.Getenv(prefix + "_DB_USER"),
            Password: os.Getenv(prefix + "_DB_PASSWORD"),
            Database: dbName,
        }

        if err := adapter.Connect(params); err != nil {
            log.Error("Connection failed: " + err.Error())
            return
        }

        if err := adapter.TestConnection(); err != nil {
            log.Warn("Connection test failed: " + err.Error())
            return
        }

        data, err := adapter.Dump(backupType)
        if err != nil {
            log.Error("Backup failed: " + err.Error())
            return
        }

        log.Success("Backup completed successfully. File path: " + string(data))
    },
}

func init() {
    backupCmd.Flags().String("db", "postgres", "Database type (postgres, mysql, mongo, sqlite)")
    backupCmd.Flags().String("type", "full", "Backup type (full, incremental, differential)")
    backupCmd.Flags().String("db-name", "", "Database name to back up")
    backupCmd.MarkFlagRequired("db-name")
}
