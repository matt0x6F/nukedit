package main

import (
	"context"
	"embed"
	"flag"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/matt0x6f/nukedit/api"
)

//go:embed all:frontend/build
var assets embed.FS

func main() {
	dryRun := flag.Bool("dry-run", false, "Perform reads but no writes")
	flag.Parse()

	db, err := api.NewDatabaseConnection("nukedit.db")
	if err != nil {
		println("Error:", err.Error())
		return
	}

	// Create an instance of the app structure
	app := NewApp(*dryRun)
	acct := api.NewAccountsService(db)
	scheduler := api.NewSchedulerService(db)
	reddit := api.NewRedditService(db, acct, *dryRun)

	// Create application with options
	err = wails.Run(&options.App{
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
			scheduler.Startup(ctx)
			reddit.Startup(ctx)
		},
		Bind: []interface{}{
			app,
			acct,
			scheduler,
			reddit,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
