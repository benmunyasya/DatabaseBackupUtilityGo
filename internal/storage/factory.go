// internal/storage/factory.go
package storage

import (
    "fmt"
)

func GetAdapter(provider string) (CloudStorage, error) {
    switch provider {
    case "gdrive":
        return NewGoogleDriveAdapter(), nil
    case "s3":
		return nil, fmt.Errorf("S3 storage provider not implemented")
    case "azure":
        return nil, fmt.Errorf("Azure storage provider not implemented")
    default:
        return nil, fmt.Errorf("unsupported storage provider: %s", provider)
    }
}
