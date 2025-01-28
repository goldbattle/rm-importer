package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"rm-exporter/backend"
	"slices"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	rm_reader   backend.RmReader
	tablet_addr string
	selection   backend.FileSelection

	rm_export      backend.RmExport
	export_options backend.RmExportOptions
}

//go:embed wails.json
var wailsJSON string

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
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

func (a *App) ReadTabletDocs(tablet_addr string) error {
	a.tablet_addr = tablet_addr

	err := a.rm_reader.Read(tablet_addr)
	if err != nil {
		return err
	}

	a.selection = backend.NewFileSelection(a.rm_reader.GetChildren())
	return nil
}

func (a *App) GetTabletFolder(id backend.DocId) []backend.DocInfo {
	return a.rm_reader.GetFolder(id)
}

func (a *App) GetTabletFolderSelection(id backend.DocId) []backend.SelectionInfo {
	return a.selection.GetFolderSelection(id)
}

func (a *App) IsIpValid(s string) bool {
	return backend.IsIpValid(s)
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
	a.rm_export = backend.InitExport(a.ctx, a.export_options, a.GetCheckedFiles(), a.tablet_addr)
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

	a.rm_export.Export(a.ctx, started, finished, failed)
}

func (a *App) OnItemSelect(id backend.DocId, selection bool) {
	a.selection.Select(id, selection)
}

/* Includes path for every checked file */
func (a *App) GetCheckedFiles() []backend.DocInfo {
	ids := a.selection.GetCheckedFiles()
	files := a.rm_reader.GetElementsByIds(ids)
	paths := a.rm_reader.GetPaths(ids)
	for i := 0; i < len(files); i += 1 {
		files[i].Path = &paths[i]
	}
	slices.SortFunc(files, func(i, j backend.DocInfo) int {
		if *i.Path < *j.Path {
			return -1
		}
		if *i.Path == *j.Path {
			return 0
		}
		return 1
	})
	return files
}

func (a *App) GetCheckedFilesCount() int {
	return a.selection.GetCheckedFilesCount()
}

func (a *App) GetPaths(ids []backend.DocId) []string {
	return a.rm_reader.GetPaths(ids)
}

func (a *App) DirectoryDialog() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return ""
	}
	return dir
}
