package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

func main() {
    var rootCmd = &cobra.Command{
        Use:   "backup-cli",
        Short: "Database Backup Utility",
        Long:  "A CLI tool to backup and restore multiple databases with compression and cloud storage support.",
    }

    // Add subcommands
    rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(scheduleCmd)
	    rootCmd.AddCommand(configCmd)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
