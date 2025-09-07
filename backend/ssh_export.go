package backend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SSHExport struct {
	ctx         context.Context
	connection  *SSHConnection
	options     RmExportOptions
	items       []DocInfo
	export_from int
}

type SSHExportOptions struct {
	Location string
	Pdf      bool
	Rmdoc    bool
}

func InitSSHExport(ctx context.Context, options RmExportOptions, items []DocInfo, connection *SSHConnection) SSHExport {
	return SSHExport{
		ctx:         ctx,
		connection:  connection,
		options:     options,
		items:       items,
		export_from: 0,
	}
}

/*
Exports all items passed in Init() method using SSH.
Calls the callbacks when:
* item started downloading;
* item download has finished;
* item download has failed.

Supports retries.
In case the last export succeeded on all items, it starts the export again from the first item;
otherwise, the export starts from the first failed item.
*/
func (s *SSHExport) Export(started, finished func(item DocInfo), failed func(item DocInfo, err error)) {
	formats := []string{}
	if s.options.Rmdoc {
		formats = append(formats, "rmdoc")
	}
	if s.options.Pdf {
		formats = append(formats, "pdf")
	}

	runtime.LogInfof(s.ctx, "[%v] SSH Export formats: %v", time.Now().UTC(), formats)
	runtime.LogInfof(s.ctx, "[%v] SSH Export location: %v", time.Now().UTC(), s.options.Location)

	for i := s.export_from; i < len(s.items); i++ {
		item := s.items[i]
		started(item)

		err := s.exportOne(item, formats)
		if err != nil {
			s.export_from = i
			failed(item, err)
			return
		}

		finished(item)
	}
}

func (s *SSHExport) exportOne(item DocInfo, formats []string) error {
	if item.IsFolder {
		return nil
	}

	runtime.LogInfof(s.ctx, "[%v] SSH Exporting item, id=%v", time.Now().UTC(), item.Id)

	// Create the local directory structure
	localDir, err := s.createLocalDirectory(item)
	if err != nil {
		return fmt.Errorf("failed to create local directory: %v", err)
	}

	// Download the document via SSH
	err = s.connection.DownloadDocument(item.Id, localDir, formats)
	if err != nil {
		return fmt.Errorf("failed to download document: %v", err)
	}

	return nil
}

func (s *SSHExport) createLocalDirectory(item DocInfo) (string, error) {
	// Create a path similar to the HTTP export but for SSH
	// Use the item's tablet path to create the directory structure
	pathParts := []string{s.options.Location}

	// Add the tablet path components
	for _, part := range item.TabletPath {
		pathParts = append(pathParts, part)
	}

	// Join all parts
	localDir := filepath.Join(pathParts...)

	// Create the directory
	err := os.MkdirAll(localDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", localDir, err)
	}

	runtime.LogDebugf(s.ctx, "[%v] Created local directory: %v, id=%v", time.Now().UTC(), localDir, item.Id)

	return localDir, nil
}
