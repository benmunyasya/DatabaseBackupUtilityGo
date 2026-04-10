package db

// ConnectionParams holds common DB connection parameters
type ConnectionParams struct {
    Host     string
    Port     int
    User     string
    Password string
    Database string
}

// DatabaseAdapter defines the contract for all DB adapters
type DatabaseAdapter interface {
    Connect(params ConnectionParams) error
    TestConnection() error
    Dump(backupType string) ([]byte, error)   // Perform backup
    Restore(data []byte) error                // Perform restore
}
