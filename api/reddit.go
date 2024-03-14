package api

import (
	"context"
	"fmt"

	"github.com/matt0x6f/nukedit/api/reddit"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Reddit struct {
	ctx      context.Context // Wails context - read only
	logger   zerolog.Logger
	database *Database
	debug    bool
	dryRun   bool

	// services
	accounts *Accounts
}

func NewRedditService(db *Database, accounts *Accounts, dryRun bool) *Reddit {
	return &Reddit{
		logger:   log.With().Str("module", "scheduler").Logger(),
		database: db,
		accounts: accounts,
		dryRun:   dryRun,
		debug:    dryRun,
	}
}

func (r *Reddit) Startup(ctx context.Context) {
	r.ctx = ctx

	if r.dryRun {
		runtime.LogDebug(r.ctx, "[Reddit Service] Dry run mode enabled")
	}

	runtime.LogDebug(r.ctx, "[Reddit Service] Starting scheduler service")
}

func (r *Reddit) Nuke(request reddit.NukeRequest) (*reddit.NukeResult, error) {
	var result reddit.NukeResult

	r.logger.Info().Msg("Nuking Reddit posts and comments")

	acct := r.accounts.GetActive()
	if acct == nil {
		runtime.LogError(r.ctx, "[Reddit Service] No active account found")
		return &result, fmt.Errorf("no active account found")
	}

	client, err := reddit.NewClient(r.ctx, acct.ClientID, acct.ClientSecret, acct.Username, acct.Password, r.dryRun)
	if err != nil {
		runtime.LogErrorf(r.ctx, "[Reddit Service] Failed to create Reddit client: %s", err.Error())
		return &result, err
	}

	// create a logger channel to stream logs from the client
	logger := make(chan reddit.Log)

	go func() {
		res, err := client.Nuke(request, logger)
		if err != nil {
			logger <- reddit.Log{Message: err.Error(), Error: err, Done: true}
		}

		result = res
	}()

	for log := range logger {
		runtime.EventsEmit(r.ctx, "nuke:log", log)
		runtime.LogDebug(r.ctx, fmt.Sprintf("[Reddit Service] "+log.Message))

		if log.Done {
			return &result, log.Error
		}
	}

	return &result, nil
}
