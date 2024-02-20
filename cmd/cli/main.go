package main

import (
	"os"

	"github.com/matt0x6f/nukedit/api/reddit"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	client, err := reddit.NewClient(os.Getenv("NUKEDIT_CLIENT_ID"), os.Getenv("NUKEDIT_CLIENT_SECRET"), os.Getenv("NUKEDIT_USERNAME"), os.Getenv("NUKEDIT_PASSWORD"))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create client")
		os.Exit(1)
	}

	client.SetDebug(true)

	err = client.EditAndDeleteAllUserComments(150)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete comments")
		os.Exit(1)
	}
}
