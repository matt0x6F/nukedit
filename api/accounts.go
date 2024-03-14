package api

import (
	"context"
	"fmt"

	"github.com/matt0x6f/nukedit/api/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Accounts struct {
	ctx      context.Context // Wails context - read only
	database *Database
	cache    []models.Account
	active   int
}

func NewAccountsService(db *Database) *Accounts {
	return &Accounts{
		database: db,
	}
}

func (s *Accounts) Startup(ctx context.Context) {
	s.ctx = ctx

	accounts, err := s.database.GetAccounts()
	if err != nil {
		runtime.LogErrorf(s.ctx, "[Accounts Service] Error fetching accounts: %s", err.Error())
		return
	}

	s.cache = accounts

	runtime.LogDebugf(s.ctx, "[Accounts Service] Found %d accounts", len(accounts))
}

func (s *Accounts) HasCredentials() bool {
	accounts, err := s.database.GetAccounts()
	if err != nil {
		runtime.LogErrorf(s.ctx, "[Accounts Service] Error fetching accounts: %s", err.Error())
		return false
	}

	if len(accounts) == 0 {
		runtime.LogDebug(s.ctx, "[Accounts Service] No accounts found")
		return false
	}

	runtime.LogDebugf(s.ctx, "[Accounts Service] Found %d accounts", len(accounts))

	return true
}

func (s *Accounts) GetAccounts() ([]models.Account, error) {
	return s.database.GetAccounts()
}

func (s *Accounts) GetAccount(key string) (*models.Account, error) {
	return s.database.GetAccount(key)
}

func (s *Accounts) SaveAccount(clientID, clientSecret, username, password string, requirepw bool) error {
	var storedpw string

	if !requirepw {
		storedpw = password
	}

	_, err := s.database.SaveAccount(clientID, clientSecret, username, storedpw, requirepw)
	if err != nil {
		return err
	}

	// refresh the accounts list
	accts, err := s.GetAccounts()
	if err != nil {
		return err
	}

	s.cache = accts

	return nil
}

func (s *Accounts) GetActive() *models.Account {
	if s.active < 0 || s.active >= len(s.cache) {
		return nil
	}

	return &s.cache[s.active]
}

func (s *Accounts) SetActive(username string) error {
	for i, acct := range s.cache {
		if acct.Username == username {
			s.active = i

			return nil
		}
	}

	return fmt.Errorf("[Accounts Service] account not found: %s", username)
}

func (s *Accounts) DeleteAccount(key string) error {
	return s.database.DeleteAccount(key)
}
