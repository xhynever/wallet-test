package services

import (
	repository "github.com/xhynever/wallet-test/repository/sqlc"
)

type Accounts interface {
	CreateAccount(req CreateAccountRequest) (repository.Account, error)
	GetAccount(req GetAccountRequest) (repository.Account, error)
	UpdateAccount(req UpdateAccountRequest) (repository.Account, error)
	DeleteAccount(req DeleteRequest) error
	ListAccounts(req ListAccountRequest) ([]repository.Account, error)
	CreateTransfer(req TxRequest) (repository.TransferTxResult, error)
	ListEntries(req ListAccountRequest) ([]repository.Entry, error)
}

type Service struct {
	Accounts
}

func NewService(store repository.Store) *Service {
	return &Service{
		Accounts: NewAccountsService(store),
	}
}
