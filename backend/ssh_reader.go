package backend

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SSHReader struct {
	connection *SSHConnection
	children   map[DocId][]DocInfo
	docById    map[DocId]DocInfo
}

func (r *SSHReader) GetDocById(id DocId) (DocInfo, bool) {
	doc, exists := r.docById[id]
	return doc, exists
}

func NewSSHReader(connection *SSHConnection) *SSHReader {
	return &SSHReader{
		connection: connection,
		children:   make(map[DocId][]DocInfo),
		docById:    make(map[DocId]DocInfo),
	}
}

func (r *SSHReader) Read() error {
	// Note: We don't have context here, so we'll use fmt.Printf for now
	// In a real implementation, we'd pass context through the call chain
	fmt.Printf("[SSH_READER] Starting to read documents...\n")
	r.children = make(map[DocId][]DocInfo)
	r.docById = make(map[DocId]DocInfo)

	fmt.Printf("[SSH_READER] Listing xochitl files...\n")
	sshFiles, err := r.connection.ListXochitlFiles()
	if err != nil {
		fmt.Printf("[SSH_READER] Failed to list SSH files: %v\n", err)
		return fmt.Errorf("failed to list SSH files: %v", err)
	}

	fmt.Printf("[SSH_READER] Found %d files, converting to DocInfo...\n", len(sshFiles))
	// Convert SSH files to DocInfo
	for _, sshFile := range sshFiles {
		docInfo := r.convertSSHFileToDocInfo(sshFile)
		r.children[docInfo.ParentId] = append(r.children[docInfo.ParentId], docInfo)
		r.docById[docInfo.Id] = docInfo
	}

	fmt.Printf("[SSH_READER] Successfully processed %d documents\n", len(sshFiles))
	return nil
}

func (r *SSHReader) convertSSHFileToDocInfo(sshFile SSHFileInfo) DocInfo {
	docInfo := DocInfo{
		Id:       sshFile.ID,
		ParentId: sshFile.Parent,
		IsFolder: sshFile.IsFolder,
		Name:     sshFile.Name,
	}

	// Parse last modified time from metadata
	if metadata, err := r.connection.ReadMetadataFile(sshFile.Path); err == nil {
		if lastModified, err := time.Parse("2006-01-02T15:04:05.000Z", metadata.LastModified); err == nil {
			docInfo.LastModified = &lastModified
		}
	}

	// Set file type for documents
	if !sshFile.IsFolder {
		fileType := "pdf" // Default to PDF, could be determined from file extension
		docInfo.FileType = &fileType
	}

	return docInfo
}

func (r *SSHReader) GetFolder(id DocId) []DocInfo {
	if items, ok := r.children[id]; ok {
		return items
	}
	return []DocInfo{}
}

func (r *SSHReader) GetCheckedFiles(selection *FileSelection) []DocInfo {
	ids := selection.GetCheckedItems()
	files := r.getElementsByIds(ids)
	r.fillPaths(files)
	return files
}

func (r *SSHReader) getElementsByIds(ids []DocId) []DocInfo {
	result := []DocInfo{}
	for _, id := range ids {
		if doc, exists := r.docById[id]; exists {
			result = append(result, doc)
		}
	}
	return result
}

func (r *SSHReader) GetChildrenMap() map[DocId][]DocInfo {
	return r.children
}

func (r *SSHReader) fillPaths(items []DocInfo) {
	for i, item := range items {
		p := r.getDisplayPath(item)
		items[i].DisplayPath = &p
		items[i].TabletPath = r.getTabletPath(item)
	}
}

func (r *SSHReader) getTabletPath(item DocInfo) []string {
	id := item.Id
	path := []string{}

	for id != "" {
		if doc, exists := r.docById[id]; exists {
			path = append(path, doc.Name)
			id = doc.ParentId
		} else {
			break
		}
	}

	// Reverse the path to get correct order
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func (r *SSHReader) getDisplayPath(item DocInfo) string {
	tabletPath := r.getTabletPath(item)
	for i, x := range tabletPath {
		if strings.Contains(x, "/") {
			tabletPath[i] = "'" + x + "'"
		}
	}
	return strings.Join(tabletPath, "/")
}

// UploadFile uploads a file using the SSH importer
func (r *SSHReader) UploadFile(localPath, fileName, parentId string) error {
	importer := NewSSHImporter(r.connection)
	return importer.UploadFile(localPath, fileName, parentId)
}

func (r *SSHReader) CreateFolder(folderName, parentId string) error {
	// Generate a new UUID for the folder
	uuidStr := uuid.New().String()

	// Create metadata file for the folder
	err := r.connection.CreateMetadataFile(uuidStr, folderName, parentId, true)
	if err != nil {
		return fmt.Errorf("failed to create folder metadata: %v", err)
	}

	// Create empty content file for folder
	err = r.connection.CreateContentFile(uuidStr, "")
	if err != nil {
		return fmt.Errorf("failed to create folder content: %v", err)
	}

	return nil
}

func (r *SSHReader) RestartXochitl() error {
	return r.connection.RestartXochitl()
}
