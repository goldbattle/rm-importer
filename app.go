package main

import (
	"context"
	"rm-exporter/backend"
	"slices"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	rm_reader   backend.RmReader
	tablet_addr string
	selection   backend.FileSelection
	rm_export   backend.RmExport
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

func (a *App) SetExportOptions(export backend.RmExport) {
	runtime.LogDebugf(a.ctx, "%v", export)
	a.rm_export = export
}

func (a *App) Export() {
	files := a.GetCheckedFiles()

	// possible states of export: downloading, finished, error
	for _, item := range files {
		runtime.EventsEmit(a.ctx, "downloading", item.Id)

		res, err := a.rm_export.Export(a.tablet_addr, item)

		if err == nil {
			runtime.EventsEmit(a.ctx, "finished", item.Id)
		} else {
			runtime.EventsEmit(a.ctx, "error", item.Id, err.Error())
			break
		}
	}
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
