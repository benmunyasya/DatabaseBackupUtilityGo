package db

import (
    "compress/gzip"
    "database/sql"
    "fmt"
    "io"
    "os"
    "os/exec"
    "path/filepath"
     "strings"
    "time"

    _ "github.com/lib/pq" // PostgreSQL driver
)

type PostgresAdapter struct {
    conn     *sql.DB
    host     string
    port     int
    user     string
    password string
    database string
}

func (p *PostgresAdapter) Connect(params ConnectionParams) error {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        params.Host, params.Port, params.User, params.Password, params.Database)

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return err
    }

    p.conn = db
    p.host = params.Host
    p.port = params.Port
    p.user = params.User
    p.password = params.Password
    p.database = params.Database

    return nil
}

func (p *PostgresAdapter) TestConnection() error {
    if p.conn == nil {
        return fmt.Errorf("no active connection")
    }
    return p.conn.Ping()
}

func (p *PostgresAdapter) Dump(backupType string) ([]byte, error) {
    // Determine backup directory
    backupDir := os.Getenv("BACKUP_DIR")
    if backupDir == "" {
        backupDir = "./backups"
    }
    if err := os.MkdirAll(backupDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create backup directory: %w", err)
    }

    // Timestamped file name
    timestamp := time.Now().Format("20060102_150405")
    fileName := fmt.Sprintf("%s_%s_%s_backup.sql", p.database, backupType, timestamp)
    fullPath := filepath.Join(backupDir, fileName)

    // Run pg_dump
    cmd := exec.Command("pg_dump",
        "-h", p.host,
        "-p", fmt.Sprintf("%d", p.port),
        "-U", p.user,
        "-d", p.database,
        "-f", fullPath,
    )
    cmd.Env = append(os.Environ(), "PGPASSWORD="+p.password)

    if err := cmd.Run(); err != nil {
        return nil, err
    }

    // Compress to .gz
    compressedFile := fullPath + ".gz"
    if err := compressFile(fullPath, compressedFile); err != nil {
        return nil, err
    }

    // Optionally remove raw .sql
    _ = os.Remove(fullPath)

    return []byte(compressedFile), nil
}

func (p *PostgresAdapter) Restore(data []byte) error {
    fileName := string(data)

    // If file is gzipped, decompress first
    if strings.HasSuffix(fileName, ".gz") {
        decompressedFile := strings.TrimSuffix(fileName, ".gz")
        if err := decompressFile(fileName, decompressedFile); err != nil {
            return err
        }
        fileName = decompressedFile
    }

    cmd := exec.Command("psql",
        "-h", p.host,
        "-p", fmt.Sprintf("%d", p.port),
        "-U", p.user,
        "-d", p.database,
        "-f", fileName,
    )
    cmd.Env = append(os.Environ(), "PGPASSWORD="+p.password)

    if err := cmd.Run(); err != nil {
        return err
    }

    return nil
}

// Helper: compress a file into gzip
func compressFile(src, dst string) error {
    inFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer inFile.Close()

    outFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer outFile.Close()

    writer := gzip.NewWriter(outFile)
    defer writer.Close()

    _, err = io.Copy(writer, inFile)
    return err
}

// Helper: decompress a gzip file
func decompressFile(src, dst string) error {
    inFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer inFile.Close()

    reader, err := gzip.NewReader(inFile)
    if err != nil {
        return err
    }
    defer reader.Close()

    outFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer outFile.Close()

    _, err = io.Copy(outFile, reader)
    return err
}
