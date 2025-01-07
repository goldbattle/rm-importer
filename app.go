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

func (a *App) GetElementsByIds(ids []backend.DocId) []backend.DocInfo {
	return a.rm_reader.GetElementsByIds(ids)
}

func (a *App) ExportPdfs(ids []backend.DocId) {
	// possible states of export: downloading, finished, error
	for _, id := range ids {
		runtime.EventsEmit(a.ctx, "downloading", id)

		res, err := backend.ExportPdf(a.tablet_addr, a.rm_reader.GetElementById(id))

		if err == nil {
			runtime.EventsEmit(a.ctx, "finished", id)
		} else {
			runtime.EventsEmit(a.ctx, "error", id, err.Error())
			break
		}

		runtime.LogDebug(a.ctx, res)
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
