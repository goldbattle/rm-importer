package main

import (
	"context"
	"remarkable-1password-sync/network"
)

// App struct
type App struct {
	ctx       context.Context
	rm_reader network.RmReader
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
	return a.rm_reader.Read(tablet_addr)
}

func (a *App) GetTabletFolder(id network.DocId) []network.DocInfo {
	return a.rm_reader.GetFolder(id)
}

func (a *App) IsIpValid(s string) bool {
	return network.IsIpValid(s)
}
