package main

import (
	"context"
)

// App struct
type App struct {
	ctx    context.Context
	dryRun bool
}

// NewApp creates a new App application struct
func NewApp(dryRun bool) *App {
	return &App{
		dryRun: dryRun,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) DryRun() bool {
	return a.dryRun
}
