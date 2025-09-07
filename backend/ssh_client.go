package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SSHConnection struct {
	host     string
	username string
	password string
	ctx      context.Context
}

type SSHMetadata struct {
	Deleted          bool   `json:"deleted"`
	LastModified     string `json:"lastModified"`
	MetadataModified bool   `json:"metadatamodified"`
	Modified         bool   `json:"modified"`
	Parent           string `json:"parent"`
	Pinned           bool   `json:"pinned"`
	Synced           bool   `json:"synced"`
	Type             string `json:"type"`
	Version          int    `json:"version"`
	VisibleName      string `json:"visibleName"`
	LastOpenedPage   *int   `json:"lastOpenedPage,omitempty"`
}

type SSHContent struct {
	ExtraMetadata  map[string]interface{} `json:"extraMetadata"`
	FileType       string                 `json:"fileType"`
	FontName       string                 `json:"fontName"`
	LastOpenedPage int                    `json:"lastOpenedPage"`
	LineHeight     int                    `json:"lineHeight"`
	Margins        int                    `json:"margins"`
	PageCount      int                    `json:"pageCount"`
	TextScale      int                    `json:"textScale"`
	Transform      map[string]int         `json:"transform"`
}

type SSHFileInfo struct {
	ID       string
	Name     string
	IsFolder bool
	Parent   string
	Type     string
	Path     string
	Size     int64
}

func NewSSHConnection(host, username, password string, ctx context.Context) *SSHConnection {
	return &SSHConnection{
		host:     host,
		username: username,
		password: password,
		ctx:      ctx,
	}
}

// GetContext returns the runtime context
func (s *SSHConnection) GetContext() context.Context {
	return s.ctx
}

func (s *SSHConnection) Connect() error {
	runtime.LogInfof(s.ctx, "[SSH] Testing connection to %s:22 with user '%s'", s.host, s.username)

	// Test the connection with a simple command
	output, err := s.executeSSHCommand("echo 'SSH connection test successful'")
	if err != nil {
		runtime.LogErrorf(s.ctx, "[SSH] Connection test failed: %v", err)
		return fmt.Errorf("SSH connection test failed: %v", err)
	}

	runtime.LogInfof(s.ctx, "[SSH] Connection test successful: %s", strings.TrimSpace(output))
	return nil
}

func (s *SSHConnection) Close() error {
	// No persistent connection to close
	return nil
}

// executeSSHCommand executes a command on the remote server using plink
func (s *SSHConnection) executeSSHCommand(command string) (string, error) {
	runtime.LogInfof(s.ctx, "[SSH] Executing command: %s", command)

	// Use plink with -batch to prevent GUI popups
	cmd := exec.Command("plink",
		"-ssh",
		"-batch", // Prevent GUI popups and interactive prompts
		"-pw", s.password,
		fmt.Sprintf("%s@%s", s.username, s.host),
		command)

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("plink command failed: %v", err)
	}

	return string(output), nil
}

// ExecuteCommand executes a command on the remote server
func (s *SSHConnection) ExecuteCommand(command string) (string, error) {
	return s.executeSSHCommand(command)
}

// ListXochitlFiles lists all files in the xochitl directory
func (s *SSHConnection) ListXochitlFiles() ([]SSHFileInfo, error) {
	runtime.LogInfo(s.ctx, "[SSH] Listing xochitl files...")

	// Use find command to get all metadata files
	command := "find ~/.local/share/remarkable/xochitl/ -name '*.metadata' -type f"
	output, err := s.executeSSHCommand(command)
	if err != nil {
		runtime.LogErrorf(s.ctx, "[SSH] Command failed: %v", err)
		return nil, fmt.Errorf("failed to list metadata files: %v", err)
	}

	runtime.LogInfof(s.ctx, "[SSH] Command output length: %d characters", len(output))

	var files []SSHFileInfo
	lines := strings.Split(strings.TrimSpace(output), "\n")
	runtime.LogInfof(s.ctx, "[SSH] Found %d metadata files", len(lines))

	for i, line := range lines {
		if line == "" {
			continue
		}

		runtime.LogInfof(s.ctx, "[SSH] Processing file %d: %s", i+1, line)

		// Get the directory and filename
		dir := filepath.Dir(line)
		baseName := filepath.Base(line)
		// Remove .metadata extension to get the ID
		id := strings.TrimSuffix(baseName, ".metadata")

		// Read metadata file
		metadata, err := s.ReadMetadataFile(line)
		if err != nil {
			runtime.LogErrorf(s.ctx, "[SSH] Failed to read metadata for %s: %v", line, err)
			continue // Skip files with invalid metadata
		}

		// Check if it's a document or collection
		isFolder := metadata.Type == "CollectionType"

		// Get file size for documents
		var size int64
		if !isFolder {
			// Use stat command to get file size
			sizeCmd := fmt.Sprintf("stat -c%%s %s 2>/dev/null || echo 0", filepath.Join(dir, id+".pdf"))
			sizeOutput, _ := s.executeSSHCommand(sizeCmd)
			fmt.Sscanf(sizeOutput, "%d", &size)
		}

		files = append(files, SSHFileInfo{
			ID:       id,
			Name:     metadata.VisibleName,
			IsFolder: isFolder,
			Parent:   metadata.Parent,
			Type:     metadata.Type,
			Path:     line,
			Size:     size,
		})

		runtime.LogInfof(s.ctx, "[SSH] Added %s: %s (%s)", id, metadata.VisibleName, metadata.Type)
	}

	runtime.LogInfof(s.ctx, "[SSH] Successfully processed %d files", len(files))
	return files, nil
}

// ReadMetadataFile reads a metadata file from the remote server
func (s *SSHConnection) ReadMetadataFile(path string) (*SSHMetadata, error) {
	command := fmt.Sprintf("cat %s", path)
	output, err := s.executeSSHCommand(command)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file %s: %v", path, err)
	}

	var metadata SSHMetadata
	err = json.Unmarshal([]byte(output), &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %v", err)
	}

	return &metadata, nil
}

// ReadContentFile reads a content file from the remote server
func (s *SSHConnection) ReadContentFile(path string) (*SSHContent, error) {
	command := fmt.Sprintf("cat %s", path)
	output, err := s.executeSSHCommand(command)
	if err != nil {
		return nil, fmt.Errorf("failed to read content file %s: %v", path, err)
	}

	var content SSHContent
	err = json.Unmarshal([]byte(output), &content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content: %v", err)
	}

	return &content, nil
}

// DownloadFile downloads a file from the remote server to local path
func (s *SSHConnection) DownloadFile(remotePath, localPath string) error {
	runtime.LogInfof(s.ctx, "[SSH] Downloading file: %s -> %s", remotePath, localPath)

	// Use plink for file download
	cmd := exec.Command("plink",
		"-ssh",
		"-batch", // Prevent GUI popups
		"-pw", s.password,
		fmt.Sprintf("%s@%s", s.username, s.host),
		fmt.Sprintf("cat %s", remotePath))

	// Set up stdout to write to local file
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	// Set up stderr capture
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start plink command: %v", err)
	}

	// Create local file
	localFile, err := os.Create(localPath)
	if err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer localFile.Close()

	// Copy stdout to local file
	_, err = io.Copy(localFile, stdout)
	if err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("failed to copy data to local file: %v", err)
	}

	// Read stderr
	errorBytes, _ := io.ReadAll(stderr)

	// Wait for the command to complete
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("plink command failed: %v, stderr: %s", err, string(errorBytes))
	}

	runtime.LogInfo(s.ctx, "[SSH] File download successful")
	return nil
}

// UploadFile uploads a local file to the remote server
func (s *SSHConnection) UploadFile(localPath, remotePath string) error {
	runtime.LogInfof(s.ctx, "[SSH] Uploading file: %s -> %s", localPath, remotePath)

	// Create remote directory if it doesn't exist
	remoteDir := filepath.Dir(remotePath)
	dirCommand := fmt.Sprintf("mkdir -p %s", remoteDir)
	_, err := s.executeSSHCommand(dirCommand)
	if err != nil {
		return fmt.Errorf("failed to create remote directory: %v", err)
	}

	// Read local file
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer localFile.Close()

	// Use plink for file upload
	cmd := exec.Command("plink",
		"-ssh",
		"-batch", // Prevent GUI popups
		"-pw", s.password,
		fmt.Sprintf("%s@%s", s.username, s.host),
		fmt.Sprintf("cat > %s", remotePath))

	// Set up stdin to provide file content
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	// Set up stderr capture
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start plink command: %v", err)
	}

	// Send file content
	go func() {
		defer stdin.Close()
		io.Copy(stdin, localFile)
	}()

	// Read stderr
	errorBytes, _ := io.ReadAll(stderr)

	// Wait for the command to complete
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("plink command failed: %v, stderr: %s", err, string(errorBytes))
	}

	runtime.LogInfo(s.ctx, "[SSH] File upload successful")
	return nil
}

// CreateMetadataFile creates a metadata file on the remote server
func (s *SSHConnection) CreateMetadataFile(id, name, parent string, isFolder bool) error {
	metadata := SSHMetadata{
		Deleted:          false,
		LastModified:     fmt.Sprintf("%d", time.Now().UnixMilli()),
		MetadataModified: false,
		Modified:         false,
		Parent:           parent,
		Pinned:           false,
		Synced:           false,
		Type:             "CollectionType",
		Version:          1,
		VisibleName:      name,
	}

	if !isFolder {
		metadata.Type = "DocumentType"
		lastOpenedPage := 0
		metadata.LastOpenedPage = &lastOpenedPage
	}

	metadataJSON, err := json.MarshalIndent(metadata, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %v", err)
	}

	// Write metadata to temporary file and upload
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("%s.metadata", id))
	err = s.WriteRemoteFile(tempFile, string(metadataJSON))
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %v", err)
	}

	// Move to final location
	finalPath := fmt.Sprintf("~/.local/share/remarkable/xochitl/%s.metadata", id)
	moveCommand := fmt.Sprintf("mv %s %s", tempFile, finalPath)
	_, err = s.executeSSHCommand(moveCommand)
	if err != nil {
		return fmt.Errorf("failed to move metadata file: %v", err)
	}

	return nil
}

// CreateContentFile creates a content file on the remote server
func (s *SSHConnection) CreateContentFile(id, fileType string) error {
	content := SSHContent{
		ExtraMetadata:  make(map[string]interface{}),
		FileType:       fileType,
		FontName:       "",
		LastOpenedPage: 0,
		LineHeight:     -1,
		Margins:        100,
		PageCount:      1,
		TextScale:      1,
		Transform: map[string]int{
			"m11": 1, "m12": 1, "m13": 1,
			"m21": 1, "m22": 1, "m23": 1,
			"m31": 1, "m32": 1, "m33": 1,
		},
	}

	contentJSON, err := json.MarshalIndent(content, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal content: %v", err)
	}

	// Write content to temporary file and upload
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("%s.content", id))
	err = s.WriteRemoteFile(tempFile, string(contentJSON))
	if err != nil {
		return fmt.Errorf("failed to write content file: %v", err)
	}

	// Move to final location
	finalPath := fmt.Sprintf("~/.local/share/remarkable/xochitl/%s.content", id)
	moveCommand := fmt.Sprintf("mv %s %s", tempFile, finalPath)
	_, err = s.executeSSHCommand(moveCommand)
	if err != nil {
		return fmt.Errorf("failed to move content file: %v", err)
	}

	return nil
}

// WriteRemoteFile writes content to a file on the remote server
func (s *SSHConnection) WriteRemoteFile(path, content string) error {
	// Create a temporary local file
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("temp_%d", time.Now().UnixNano()))
	err := os.WriteFile(tempFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile)

	// Upload the temporary file to the remote location
	err = s.UploadFile(tempFile, path)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil
}

// DownloadDocument downloads a complete document (PDF, metadata, content) via SSH
func (s *SSHConnection) DownloadDocument(id, localDir string, formats []string) error {
	runtime.LogInfof(s.ctx, "[SSH] Downloading document %s to %s", id, localDir)

	// Create local directory
	err := os.MkdirAll(localDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create local directory: %v", err)
	}

	// Download each requested format
	for _, format := range formats {
		remotePath := fmt.Sprintf("~/.local/share/remarkable/xochitl/%s.%s", id, format)
		localPath := filepath.Join(localDir, fmt.Sprintf("%s.%s", id, format))

		err = s.DownloadFile(remotePath, localPath)
		if err != nil {
			runtime.LogErrorf(s.ctx, "[SSH] Failed to download %s: %v", format, err)
			// Continue with other formats even if one fails
		} else {
			runtime.LogInfof(s.ctx, "[SSH] Successfully downloaded %s", format)
		}
	}

	// Download metadata file
	metadataPath := fmt.Sprintf("~/.local/share/remarkable/xochitl/%s.metadata", id)
	localMetadataPath := filepath.Join(localDir, fmt.Sprintf("%s.metadata", id))
	err = s.DownloadFile(metadataPath, localMetadataPath)
	if err != nil {
		runtime.LogErrorf(s.ctx, "[SSH] Failed to download metadata: %v", err)
	}

	// Download content file
	contentPath := fmt.Sprintf("~/.local/share/remarkable/xochitl/%s.content", id)
	localContentPath := filepath.Join(localDir, fmt.Sprintf("%s.content", id))
	err = s.DownloadFile(contentPath, localContentPath)
	if err != nil {
		runtime.LogErrorf(s.ctx, "[SSH] Failed to download content: %v", err)
	}

	return nil
}

// CreateDirectories creates the necessary directories for a document
func (s *SSHConnection) CreateDirectories(id string) error {
	basePath := fmt.Sprintf("~/.local/share/remarkable/xochitl/%s", id)

	dirs := []string{
		fmt.Sprintf("%s.cache", basePath),
		fmt.Sprintf("%s.highlights", basePath),
		fmt.Sprintf("%s.thumbnails", basePath),
	}

	for _, dir := range dirs {
		_, err := s.executeSSHCommand(fmt.Sprintf("mkdir -p %s", dir))
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	return nil
}

// RestartXochitl restarts the xochitl service
func (s *SSHConnection) RestartXochitl() error {
	_, err := s.executeSSHCommand("systemctl restart xochitl")
	if err != nil {
		return fmt.Errorf("failed to restart xochitl: %v", err)
	}
	return nil
}
