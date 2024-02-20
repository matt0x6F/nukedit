// Package db is the database service
package db

import (
	"github.com/glebarez/sqlite"
	"github.com/matt0x6f/nukedit/api/models"
	"gorm.io/gorm"
)

var Service *DBService

type DBService struct {
	db *gorm.DB
}

func NewConnection(path string) (*DBService, error) {
	db, error := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if error != nil {
		return nil, error
	}

	return &DBService{db: db}, nil
}

func Run(svc *DBService) {
	svc.db.AutoMigrate(&models.Account{}, &models.Config{})

	Service = svc
}

func (s *DBService) Close() {
	sqlDB, _ := s.db.DB()
	sqlDB.Close()
}

func (s *DBService) GetAccount(key string) (*models.Account, error) {
	account := &models.Account{}
	result := s.db.First(account, "id = ?", key)
	return account, result.Error
}

func (s *DBService) SaveAccount(clientID, clientSecret, username, password string, requirepw bool) (*models.Account, error) {
	acct := &models.Account{
		ClientID:        clientID,
		ClientSecret:    clientSecret,
		Username:        username,
		Password:        password,
		RequirePassword: requirepw,
	}

	result := s.db.Save(&acct)

	return acct, result.Error
}

func (s *DBService) DeleteAccount(key string) error {
	result := s.db.Delete(&models.Account{}, "id = ?", key)
	return result.Error
}

func (s *DBService) GetConfig(key string) (*models.Config, error) {
	config := &models.Config{}
	result := s.db.First(config, "key = ?", key)
	return config, result.Error
}

func (s *DBService) SaveConfig(config *models.Config) error {
	result := s.db.Save(config)
	return result.Error
}

func (s *DBService) DeleteConfig(key string) error {
	result := s.db.Delete(&models.Config{}, "key = ?", key)
	return result.Error
}

func (s *DBService) GetAccounts() ([]models.Account, error) {
	var accounts []models.Account
	result := s.db.Find(&accounts)
	return accounts, result.Error
}

func (s *DBService) GetConfigs() ([]models.Config, error) {
	var configs []models.Config
	result := s.db.Find(&configs)
	return configs, result.Error
}
