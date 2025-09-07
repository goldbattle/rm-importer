package backend

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SSHImporter handles file uploads to the reMarkable device via SSH
type SSHImporter struct {
	connection *SSHConnection
}

// NewSSHImporter creates a new SSH importer
func NewSSHImporter(connection *SSHConnection) *SSHImporter {
	return &SSHImporter{
		connection: connection,
	}
}

// UploadFile uploads a file to the reMarkable device and returns the generated UUID
func (s *SSHImporter) UploadFile(localPath, fileName, parentId string) (string, error) {
	ctx := s.connection.GetContext()
	runtime.LogInfof(ctx, "[SSH_IMPORT] Starting upload: %s -> %s (parent: %s)", localPath, fileName, parentId)

	// Validate file extension - only allow PDF files for now
	ext := strings.ToLower(filepath.Ext(fileName))
	if ext != ".pdf" {
		return "", fmt.Errorf("only PDF files are supported, got: %s", ext)
	}

	// Generate a new UUID for the document
	uuidStr := uuid.New().String()
	runtime.LogInfof(ctx, "[SSH_IMPORT] Generated UUID: %s", uuidStr)

	// Remove extension from filename for visible name
	visibleName := strings.TrimSuffix(fileName, ext)

	// Upload the main file
	remoteFilePath := fmt.Sprintf("~/.local/share/remarkable/xochitl/%s%s", uuidStr, ext)
	err := s.connection.UploadFile(localPath, remoteFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	runtime.LogInfo(ctx, "[SSH_IMPORT] File upload successful")

	// Create metadata file
	err = s.connection.CreateMetadataFile(uuidStr, visibleName, parentId, false)
	if err != nil {
		return "", fmt.Errorf("failed to create metadata: %v", err)
	}
	runtime.LogInfo(ctx, "[SSH_IMPORT] Metadata file created successfully")

	// Create content file
	fileType := strings.TrimPrefix(ext, ".")
	err = s.connection.CreateContentFile(uuidStr, fileType)
	if err != nil {
		return "", fmt.Errorf("failed to create content file: %v", err)
	}
	runtime.LogInfo(ctx, "[SSH_IMPORT] Content file created successfully")

	// Create necessary directories
	err = s.connection.CreateDirectories(uuidStr)
	if err != nil {
		return "", fmt.Errorf("failed to create directories: %v", err)
	}
	runtime.LogInfo(ctx, "[SSH_IMPORT] Directories created successfully")

	// Restart xochitl to make the new document visible
	err = s.connection.RestartXochitl()
	if err != nil {
		return "", fmt.Errorf("failed to restart xochitl: %v", err)
	}
	runtime.LogInfo(ctx, "[SSH_IMPORT] Xochitl restarted successfully")

	runtime.LogInfo(ctx, "[SSH_IMPORT] Upload process completed successfully!")
	return uuidStr, nil
}
