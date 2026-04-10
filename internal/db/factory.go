package db

import "fmt"

// GetAdapter returns the correct DatabaseAdapter based on dbType
func GetAdapter(dbType string) (DatabaseAdapter, error) {
    switch dbType {
    case "postgres":
        return &PostgresAdapter{}, nil
    // case "mysql":
    //     return &MySQLAdapter{}, nil
    // case "mongo":
    //     return &MongoAdapter{}, nil
    // case "sqlite":
    //     return &SQLiteAdapter{}, nil
    default:
        return nil, fmt.Errorf("unsupported database type: %s", dbType)
    }
}
