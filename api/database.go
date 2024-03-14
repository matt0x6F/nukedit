package api

import (
	"github.com/glebarez/sqlite"
	"github.com/matt0x6f/nukedit/api/models"
	"gorm.io/gorm"
)

type Database struct {
	database *gorm.DB
}

func NewDatabaseConnection(path string) (*Database, error) {
	db, error := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if error != nil {
		return nil, error
	}

	return &Database{database: db}, nil
}

func (s *Database) Migrate() error {
	return s.database.AutoMigrate(&models.Account{}, &models.Config{}, &models.Schedule{})
}

func (s *Database) Close() {
	sqlDB, _ := s.database.DB()
	sqlDB.Close()
}

func (s *Database) GetAccount(key string) (*models.Account, error) {
	account := &models.Account{}
	result := s.database.First(account, "id = ?", key)
	return account, result.Error
}

func (s *Database) SaveAccount(clientID, clientSecret, username, password string, requirepw bool) (*models.Account, error) {
	acct := &models.Account{
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		Username:        username,
		Password:        password,
		RequirePassword: requirepw,
	}

	result := s.database.Save(&acct)

	return acct, result.Error
}

func (s *Database) DeleteAccount(key string) error {
	result := s.database.Delete(&models.Account{}, "id = ?", key)
	return result.Error
}

func (s *Database) GetConfig(key string) (*models.Config, error) {
	config := &models.Config{}
	result := s.database.First(config, "key = ?", key)
	return config, result.Error
}

func (s *Database) SaveConfig(config *models.Config) error {
	result := s.database.Save(config)
	return result.Error
}

func (s *Database) DeleteConfig(key string) error {
	result := s.database.Delete(&models.Config{}, "key = ?", key)
	return result.Error
}

func (s *Database) GetAccounts() ([]models.Account, error) {
	var accounts []models.Account
	result := s.database.Find(&accounts)
	return accounts, result.Error
}

func (s *Database) GetConfigs() ([]models.Config, error) {
	var configs []models.Config
	result := s.database.Find(&configs)
	return configs, result.Error
}

func (s *Database) GetSchedules() ([]models.Schedule, error) {
	var schedules []models.Schedule
	result := s.database.Find(&schedules)
	return schedules, result.Error
}

func (s *Database) GetSchedule(key int) (*models.Schedule, error) {
	schedule := &models.Schedule{}
	result := s.database.First(schedule, "id = ?", key)
	return schedule, result.Error
}

func (s *Database) SaveSchedule(schedule *models.Schedule) error {
	result := s.database.Save(schedule)
	return result.Error
}

func (s *Database) DeleteSchedule(key int) error {
	result := s.database.Delete(&models.Schedule{}, "id = ?", key)
	return result.Error
}

func (s *Database) GetScheduleByUsername(username string) ([]models.Schedule, error) {
	var schedules []models.Schedule
	result := s.database.Find(&schedules, "username = ?", username)
	return schedules, result.Error
}
