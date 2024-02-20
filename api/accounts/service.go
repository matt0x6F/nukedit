// Package accounts is used by the frontend and backend to manage accounts and credentials
package accounts

import (
	"context"
	"fmt"

	"github.com/matt0x6f/nukedit/api/db"
	"github.com/matt0x6f/nukedit/api/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Service struct {
	ctx      context.Context
	accounts []models.Account
	active   int
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Startup(ctx context.Context) {
	s.ctx = ctx

	accounts, err := db.Service.GetAccounts()
	if err != nil {
		runtime.LogErrorf(s.ctx, "[nukedit] Error fetching accounts: %s", err.Error())
		return
	}

	s.accounts = accounts
}

func (s *Service) HasCredentials() bool {
	accounts, err := db.Service.GetAccounts()
	if err != nil {
		runtime.LogErrorf(s.ctx, "[nokedit] Error fetching accounts: %s", err.Error())
		return false
	}

	if len(accounts) == 0 {
		runtime.LogDebug(s.ctx, "[nukedit] No accounts found")
		return false
	}

	runtime.LogDebugf(s.ctx, "[nukedit] Found %d accounts", len(accounts))

	return true
}

func (s *Service) GetAccounts() ([]models.Account, error) {
	return db.Service.GetAccounts()
}

func (s *Service) GetAccount(key string) (*models.Account, error) {
	return db.Service.GetAccount(key)
}

func (s *Service) SaveAccount(clientID, clientSecret, username, password string, requirepw bool) error {
	var storedpw string

	if !requirepw {
		storedpw = password
	}

	_, err := db.Service.SaveAccount(clientID, clientSecret, username, storedpw, requirepw)
	if err != nil {
		return err
	}

	// refresh the accounts list
	accts, err := s.GetAccounts()
	if err != nil {
		return err
	}

	s.accounts = accts

	return nil
}

func (s *Service) GetActive() *models.Account {
	if s.active < 0 || s.active >= len(s.accounts) {
		return nil
	}

	return &s.accounts[s.active]
}

func (s *Service) SetActive(username string) error {
	for i, acct := range s.accounts {
		if acct.Username == username {
			s.active = i

			return nil
		}
	}

	return fmt.Errorf("account not found: %s", username)
}

func (s *Service) DeleteAccount(key string) error {
	return db.Service.DeleteAccount(key)
}
