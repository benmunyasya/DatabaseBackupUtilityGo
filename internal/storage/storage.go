// internal/storage/storage.go
package storage

// CloudStorage defines the contract for uploading backups to cloud providers.
type CloudStorage interface {
    Upload(localPath string, remotePath string) error
}
