package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"rm-importer/backend"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	rm_reader   backend.RmReader
	ssh_reader  *backend.SSHReader
	ssh_conn    *backend.SSHConnection
	tablet_addr string
	selection   backend.FileSelection

	rm_export      backend.RmExport
	ssh_export     backend.SSHExport
	export_options backend.RmExportOptions

	// SSH connection details
	ssh_host     string
	ssh_username string
	ssh_password string
	safe_mode    bool
	hybrid_mode  bool
}

//go:embed wails.json
var wailsJSON string

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		safe_mode: true, // Safe mode enabled by default
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetAppVersion() string {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(wailsJSON), &m)
	if err != nil {
		return "0.0.0"
	}
	return m["info"].(map[string]interface{})["productVersion"].(string)
}

func (a *App) IsIpValid(s string) bool {
	return backend.IsIpValid(s)
}

func (a *App) ReadDocs(tablet_addr string) error {
	a.tablet_addr = tablet_addr

	err := a.rm_reader.Read(tablet_addr)
	if err != nil {
		return err
	}

	a.selection = backend.NewFileSelection(a.rm_reader.GetChildrenMap())
	return nil
}

func (a *App) ConnectSSH(host, username, password string) error {
	runtime.LogInfof(a.ctx, "[APP] ConnectSSH called with host=%s, username=%s, password_length=%d", host, username, len(password))

	a.ssh_host = host
	a.ssh_username = username
	a.ssh_password = password

	// Create SSH connection
	runtime.LogInfo(a.ctx, "[APP] Creating SSH connection...")
	a.ssh_conn = backend.NewSSHConnection(host, username, password, a.ctx)

	runtime.LogInfo(a.ctx, "[APP] Attempting SSH connection...")
	err := a.ssh_conn.Connect()
	if err != nil {
		runtime.LogErrorf(a.ctx, "[APP] SSH connection failed: %v", err)
		return err
	}

	// Create SSH reader
	runtime.LogInfo(a.ctx, "[APP] Creating SSH reader...")
	a.ssh_reader = backend.NewSSHReader(a.ssh_conn)

	runtime.LogInfo(a.ctx, "[APP] Reading documents from SSH...")
	err = a.ssh_reader.Read()
	if err != nil {
		runtime.LogErrorf(a.ctx, "[APP] SSH read failed: %v", err)
		return err
	}

	// Create selection from SSH reader
	runtime.LogInfo(a.ctx, "[APP] Creating file selection...")
	a.selection = backend.NewFileSelection(a.ssh_reader.GetChildrenMap())

	runtime.LogInfo(a.ctx, "[APP] SSH connection and setup completed successfully!")
	return nil
}

func (a *App) ConnectSSHForUploads(host, username, password string) error {
	runtime.LogInfof(a.ctx, "[APP] ConnectSSHForUploads called with host=%s, username=%s, password_length=%d", host, username, len(password))

	a.ssh_host = host
	a.ssh_username = username
	a.ssh_password = password

	// Create SSH connection for uploads only (no file reading)
	runtime.LogInfo(a.ctx, "[APP] Creating SSH connection for uploads...")
	a.ssh_conn = backend.NewSSHConnection(host, username, password, a.ctx)

	runtime.LogInfo(a.ctx, "[APP] Attempting SSH connection...")
	err := a.ssh_conn.Connect()
	if err != nil {
		runtime.LogErrorf(a.ctx, "[APP] SSH connection failed: %v", err)
		return err
	}

	// Create SSH reader for upload functionality only
	// Note: We don't call Read() here because we're using HTTP for file listing
	runtime.LogInfo(a.ctx, "[APP] Creating SSH reader for uploads...")
	a.ssh_reader = backend.NewSSHReader(a.ssh_conn)

	runtime.LogInfo(a.ctx, "[APP] SSH connection for uploads completed successfully!")
	return nil
}

func (a *App) DisconnectSSH() error {
	if a.ssh_conn != nil {
		return a.ssh_conn.Close()
	}
	return nil
}

func (a *App) SetSafeMode(enabled bool) {
	a.safe_mode = enabled
}

func (a *App) SetHybridMode(enabled bool) {
	a.hybrid_mode = enabled
}

func (a *App) GetSafeMode() bool {
	return a.safe_mode
}

func (a *App) IsSSHMode() bool {
	return a.ssh_reader != nil
}

func (a *App) GetFolder(id backend.DocId) []backend.DocInfo {
	// In hybrid mode, always use HTTP for file listing
	if a.hybrid_mode {
		return a.rm_reader.GetFolder(id)
	}
	// In pure SSH mode, use SSH
	if a.ssh_reader != nil {
		return a.ssh_reader.GetFolder(id)
	}
	// Fallback to HTTP
	return a.rm_reader.GetFolder(id)
}

func (a *App) GetFolderSelection(id backend.DocId) []backend.SelectionInfo {
	return a.selection.GetFolderSelection(id)
}

func (a *App) GetItemSelection(id backend.DocId) backend.SelectionInfo {
	return a.selection.GetItemSelection(id)
}

func (a *App) OnItemSelect(id backend.DocId, selection bool) {
	a.selection.Select(id, selection)
}

func (a *App) GetCheckedFilesCount() int {
	return a.selection.GetCheckedFilesCount()
}

type RmExportOptions struct {
	Format   string
	Location string
}

func (a *App) SetExportOptions(options backend.RmExportOptions) {
	a.export_options = options
}

func (a *App) GetExportOptions() backend.RmExportOptions {
	return a.export_options
}

func (a *App) InitExport() {
	if a.ssh_reader != nil {
		// Use SSH export
		runtime.LogInfo(a.ctx, "[APP] Initializing SSH export")
		a.ssh_export = backend.InitSSHExport(a.ctx, a.export_options, a.GetCheckedFiles(), a.ssh_conn)
	} else {
		// Use HTTP export
		runtime.LogInfo(a.ctx, "[APP] Initializing HTTP export")
		a.rm_export = backend.InitExport(a.ctx, a.export_options, a.GetCheckedFiles(), a.tablet_addr)
	}
}

func (a *App) Export() {
	started := func(item backend.DocInfo) {
		runtime.LogInfof(a.ctx, "[%v] Started file id=%v", time.Now().UTC(), item.Id)
		runtime.EventsEmit(a.ctx, "started", item.Id)
	}

	finished := func(item backend.DocInfo) {
		runtime.LogInfof(a.ctx, "[%v] Finished file id=%v", time.Now().UTC(), item.Id)
		runtime.EventsEmit(a.ctx, "finished", item.Id)
	}

	failed := func(item backend.DocInfo, err error) {
		runtime.LogInfof(a.ctx, "[%v] Failed file id=%v, error: %v", time.Now().UTC(), item.Id, err)
		runtime.EventsEmit(a.ctx, "failed", item.Id, err.Error())
	}

	if a.ssh_reader != nil {
		// Use SSH export
		runtime.LogInfo(a.ctx, "[APP] Starting SSH export")
		a.ssh_export.Export(started, finished, failed)
	} else {
		// Use HTTP export
		runtime.LogInfo(a.ctx, "[APP] Starting HTTP export")
		a.rm_export.Export(started, finished, failed)
	}
}

/* Includes path for every checked file */
func (a *App) GetCheckedFiles() []backend.DocInfo {
	if a.ssh_reader != nil {
		return a.ssh_reader.GetCheckedFiles(&a.selection)
	}
	return a.rm_reader.GetCheckedFiles(&a.selection)
}

func (a *App) DirectoryDialog() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return ""
	}
	return dir
}

func (a *App) FileDialog() string {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "PDF Files",
				Pattern:     "*.pdf",
			},
		},
	})
	if err != nil {
		return ""
	}
	return file
}

func (a *App) UploadFileSSH(localPath, fileName, parentId string) (string, error) {
	runtime.LogInfof(a.ctx, "[APP] UploadFileSSH called: localPath=%s, fileName=%s, parentId=%s", localPath, fileName, parentId)

	if a.ssh_reader == nil {
		runtime.LogError(a.ctx, "[APP] SSH connection not established")
		return "", fmt.Errorf("SSH connection not established")
	}

	if a.safe_mode {
		runtime.LogError(a.ctx, "[APP] Upload blocked: safe mode is enabled")
		return "", fmt.Errorf("upload blocked: safe mode is enabled")
	}

	runtime.LogInfo(a.ctx, "[APP] Starting file upload...")
	uuid, err := a.ssh_reader.UploadFile(localPath, fileName, parentId)
	if err != nil {
		runtime.LogErrorf(a.ctx, "[APP] Upload failed: %v", err)
		return "", err
	}

	runtime.LogInfof(a.ctx, "[APP] Upload completed successfully with UUID: %s", uuid)
	return uuid, nil
}

func (a *App) CreateFolderSSH(folderName, parentId string) error {
	if a.ssh_reader == nil {
		return fmt.Errorf("SSH connection not established")
	}

	if a.safe_mode {
		return fmt.Errorf("folder creation blocked: safe mode is enabled")
	}

	return a.ssh_reader.CreateFolder(folderName, parentId)
}

func (a *App) RestartXochitlSSH() error {
	if a.ssh_reader == nil {
		return fmt.Errorf("SSH connection not established")
	}

	if a.safe_mode {
		return fmt.Errorf("restart blocked: safe mode is enabled")
	}

	return a.ssh_reader.RestartXochitl()
}

func (a *App) TestSSHConnection(host, username, password string) error {
	runtime.LogInfof(a.ctx, "[APP] Testing SSH connection to %s with user %s", host, username)

	// Create a temporary SSH connection for testing
	testConn := backend.NewSSHConnection(host, username, password, a.ctx)

	// Test SSH connection
	err := testConn.Connect()
	if err != nil {
		runtime.LogErrorf(a.ctx, "[APP] SSH connection test failed: %v", err)
		return err
	}

	runtime.LogInfo(a.ctx, "[APP] SSH connection test successful!")

	// Close the test connection
	defer testConn.Close()

	return nil
}
