package main

import (
	"context"
	"remarkable-1p-sync/network"
)

// App struct
type App struct {
	ctx context.Context
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

func (a *App) ReadFiles() ([]network.DocInfo, error) {
	return network.ReadFiles()
}
