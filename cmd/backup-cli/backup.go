package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/benmunyasya/dbbackuputility/internal/db"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
    Use:   "backup",
    Short: "Run a database backup",
    Long:  "Perform a backup operation for the specified database type (Postgres, MySQL, MongoDB, SQLite).",
    Run: func(cmd *cobra.Command, args []string) {
        dbType, _ := cmd.Flags().GetString("db")
        backupType, _ := cmd.Flags().GetString("type")

        fmt.Printf("Running %s backup for %s database...\n", backupType, dbType)

        adapter, err := db.GetAdapter(dbType)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        // Build env prefix (e.g., POSTGRES, MYSQL)
        prefix := strings.ToUpper(dbType)

        port, err := strconv.Atoi(os.Getenv(prefix + "_DB_PORT"))
        if err != nil {
            fmt.Println("Invalid port:", err)
            return
        }

        params := db.ConnectionParams{
            Host:     os.Getenv(prefix + "_DB_HOST"),
            Port:     port,
            User:     os.Getenv(prefix + "_DB_USER"),
            Password: os.Getenv(prefix + "_DB_PASSWORD"),
            Database: os.Getenv(prefix + "_DB_NAME"),
        }

        if err := adapter.Connect(params); err != nil {
            fmt.Println("Connection failed:", err)
            return
        }

        if err := adapter.TestConnection(); err != nil {
            fmt.Println("Connection test failed:", err)
            return
        }

        data, err := adapter.Dump(backupType)
        if err != nil {
            fmt.Println("Backup failed:", err)
            return
        }

        fmt.Printf("Backup completed successfully. Data size: %d bytes\n", len(data))
    },
}

func init() {
    backupCmd.Flags().String("db", "postgres", "Database type (postgres, mysql, mongo, sqlite)")
    backupCmd.Flags().String("type", "full", "Backup type (full, incremental, differential)")
}
