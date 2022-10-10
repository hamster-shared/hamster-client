package app

import (
	"context"
	"hamster-client/module/account"
)

type Account struct {
	accountService account.Service
	ctx            context.Context
}

func NewAccountApp(service account.Service) Account {
	return Account{
		accountService: service,
	}
}

func (s *Account) WailsInit(ctx context.Context) error {
	s.ctx = ctx
	return nil
}

// GetAccountInfo get user information
func (s *Account) GetAccountInfo() (account.Account, error) {
	info, err := s.accountService.GetAccount()
	return info, err
}

// IsAccount determine if the user exists
func (s *Account) IsAccount() bool {
	_, err := s.accountService.GetAccount()
	return err == nil
}

// IsAccountSetting determine whether to configure user information
func (s *Account) IsAccountSetting() bool {
	info, err := s.accountService.GetAccount()
	if err != nil {
		return false
	}
	return info.PublicKey != ""
}
