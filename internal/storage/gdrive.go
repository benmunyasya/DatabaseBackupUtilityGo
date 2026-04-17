// internal/storage/gdrive.go
package storage

import (
    "context"
    "fmt"
    "os"

    "google.golang.org/api/drive/v3"
    "google.golang.org/api/option"
)

type GoogleDriveAdapter struct {
    service *drive.Service
}

func NewGoogleDriveAdapter() *GoogleDriveAdapter {
    ctx := context.Background()
    credFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
    srv, err := drive.NewService(ctx, option.WithCredentialsFile(credFile))
    if err != nil {
        panic(fmt.Errorf("failed to init Google Drive client: %w", err))
    }
    return &GoogleDriveAdapter{service: srv}
}

func (g *GoogleDriveAdapter) Upload(localPath, remotePath string) error {
    f, err := os.Open(localPath)
    if err != nil {
        return fmt.Errorf("could not open local file: %w", err)
    }
    defer f.Close()

    folderID := os.Getenv("GDRIVE_FOLDER_ID")
    if folderID == "" {
        // Graceful fallback: just log a warning and skip upload
        fmt.Println("[WARN] GDRIVE_FOLDER_ID not set, skipping cloud upload.")
        return nil
    }

	fmt.Println("[INFO] Uploading backup to Google Drive folder ID:", folderID)

    fileMetadata := &drive.File{
        Name:    remotePath,
        Parents: []string{folderID},
    }

    uploadedFile, err := g.service.Files.Create(fileMetadata).Media(f).Do()
    if err != nil {
        // Graceful fallback: log error but don’t break backup flow
        fmt.Println("[WARN] Cloud upload failed:", err)
        return nil
    }

    fmt.Println("[SUCCESS] Backup uploaded to Google Drive. File ID:", uploadedFile.Id)
    return nil
}
