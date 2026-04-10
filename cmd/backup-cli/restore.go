package main

import (
    "fmt"
    "os"
	"strconv"
    "strings"

    "github.com/spf13/cobra"
    "github.com/benmunyasya/dbbackuputility/internal/db"
)

var restoreCmd = &cobra.Command{
    Use:   "restore",
    Short: "Restore a database from backup",
    Long:  "Restore a database from a backup file. Supports full restore and selective restore (tables/collections).",
    Run: func(cmd *cobra.Command, args []string) {
        dbType, _ := cmd.Flags().GetString("db")
        filePath, _ := cmd.Flags().GetString("file")

        fmt.Printf("Restoring %s database from file: %s...\n", dbType, filePath)

        adapter, err := db.GetAdapter(dbType)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
		port, err := strconv.Atoi(os.Getenv(strings.ToUpper(dbType) + "_DB_PORT"))
		if err != nil {
			fmt.Println("Invalid port:", err)
			return
		}

        // Build env prefix (e.g., POSTGRES, MYSQL)
        prefix := strings.ToUpper(dbType)

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

        // Placeholder: later we’ll read file contents
        fakeData := []byte("FAKE_BACKUP_DATA")

        if err := adapter.Restore(fakeData); err != nil {
            fmt.Println("Restore failed:", err)
            return
        }

        fmt.Println("Restore completed successfully.")
    },
}

func init() {
    restoreCmd.Flags().String("db", "postgres", "Database type (postgres, mysql, mongo, sqlite)")
    restoreCmd.Flags().String("file", "", "Path to backup file")
    restoreCmd.MarkFlagRequired("file")
}
