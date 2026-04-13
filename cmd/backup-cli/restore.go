package main

import (
    "os"
    "strconv"
    "strings"

    "github.com/benmunyasya/dbbackuputility/internal/db"
    "github.com/benmunyasya/dbbackuputility/internal/log" // 👈 import logging utility
    "github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
    Use:   "restore",
    Short: "Restore a database from backup",
    Long:  "Restore a database from a backup file for the specified database type (Postgres, MySQL, MongoDB, SQLite).",
    Run: func(cmd *cobra.Command, args []string) {
        dbType, _ := cmd.Flags().GetString("db")
        dbName, _ := cmd.Flags().GetString("db-name")
        filePath, _ := cmd.Flags().GetString("file")

        log.Info("Restoring " + dbType + " database (" + dbName + ") from file: " + filePath)

        adapter, err := db.GetAdapter(dbType)
        if err != nil {
            log.Error("Adapter error: " + err.Error())
            return
        }

        // Build env prefix (e.g., POSTGRES, MYSQL)
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

        if err := adapter.Restore([]byte(filePath)); err != nil {
            log.Error("Restore failed: " + err.Error())
            return
        }

        log.Success("Restore completed successfully.")
    },
}

func init() {
    restoreCmd.Flags().String("db", "postgres", "Database type (postgres, mysql, mongo, sqlite)")
    restoreCmd.Flags().String("db-name", "", "Database name to restore into")
    restoreCmd.Flags().String("file", "", "Backup file path (.sql or .sql.gz)")
    restoreCmd.MarkFlagRequired("db-name")
    restoreCmd.MarkFlagRequired("file")
}
