package main

import (
	"context"
	"os"

	"github.com/matt0x6f/nukedit/api/accounts"
	"github.com/matt0x6f/nukedit/api/db"
	"github.com/matt0x6f/nukedit/api/reddit"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx        context.Context
	accountsvc *accounts.Service
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) SetAccountService(svc *accounts.Service) {
	a.accountsvc = svc
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dbService, err := db.NewConnection("nukedit.db")
	if err != nil {
		runtime.LogErrorf(ctx, "Error creating database service: %s", err.Error())
		os.Exit(1)
	}

	runtime.LogDebug(ctx, "Running database service")

	db.Run(dbService)
}

func (a *App) Nuke(req reddit.NukeRequest) error {
	account := a.accountsvc.GetActive()

	client, err := reddit.NewClient(account.ClientID, account.ClientSecret, account.Username, account.Password)
	if err != nil {
		return err
	}

	if req.Scheduled {

		return nil
	}

	return client.Nuke(req)
}
