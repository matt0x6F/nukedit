package main

import (
	"context"
	"embed"

	"github.com/matt0x6f/nukedit/api/accounts"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	acct := accounts.NewService()

	app.SetAccountService(acct)

	// Create application with options
	err := wails.Run(&options.App{
		Title:    "nukedit",
		Width:    1024,
		Height:   768,
		MinWidth: 1024,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			acct.Startup(ctx)
		},
		Bind: []interface{}{
			app,
			acct,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
