package api

import (
	"context"
	"fmt"

	"github.com/matt0x6f/nukedit/api/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Scheduler struct {
	ctx    context.Context // Wails context - read only
	logger zerolog.Logger

	crons []models.Schedule

	database *Database
}

func NewSchedulerService(db *Database) *Scheduler {
	return &Scheduler{
		logger:   log.With().Str("module", "scheduler").Logger(),
		database: db,
	}
}

func (s *Scheduler) Startup(ctx context.Context) {
	s.ctx = ctx

	runtime.LogDebug(s.ctx, "[Scheduler Service] Scheduler service started")
}

func (s *Scheduler) GetSchedules() ([]models.Schedule, error) {
	return s.database.GetSchedules()
}

func (s *Scheduler) SaveSchedule(schedule models.Schedule) error {
	fmt.Printf("Saving schedule: %+v\n", schedule)

	return s.database.SaveSchedule(&schedule)
}

func (s *Scheduler) DeleteSchedule(key int) error {
	return s.database.DeleteSchedule(key)
}

func (s *Scheduler) GetScheduleByUsername(username string) ([]models.Schedule, error) {
	return s.database.GetScheduleByUsername(username)
}
