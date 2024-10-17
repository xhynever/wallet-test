package services

import (
	"context"

	repository "github.com/xhynever/wallet-test/repository/sqlc"
)

type AccountsService struct {
	store repository.Store
}

func NewAccountsService(store repository.Store) *AccountsService {
	return &AccountsService{store: store}
}

type CreateAccountRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency" `
}

func (service *AccountsService) CreateAccount(req CreateAccountRequest) (repository.Account, error) {

	arg := repository.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := service.store.CreateAccount(context.Background(), arg)
	if err != nil {
		return account, err
	}
	return account, nil
}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (service *AccountsService) GetAccount(req GetAccountRequest) (repository.Account, error) {
	account, err := service.store.GetAccount(context.Background(), req.ID)
	if err != nil {
		return account, err
	}
	return account, nil
}

type UpdateAccountRequest struct {
	Owner   string `json:"owner" binding:"required"`
	Balance int64  `json:"balance" binding:"required,min=0"`
	ID      int64  `uri:"id" binding:"required,min=1"`
}

func (service *AccountsService) UpdateAccount(req UpdateAccountRequest) (repository.Account, error) {
	arg := repository.UpdateAccountParams{
		Owner:   req.Owner,
		Balance: req.Balance,
		ID:      req.ID,
	}
	account, err := service.store.UpdateAccount(context.Background(), arg)
	if err != nil {
		return account, err
	}
	return account, nil
}

type DeleteRequest struct {
	ID int64 `uri:"id" `
}

func (service *AccountsService) DeleteAccount(req DeleteRequest) error {
	// 若账户有余额，执行取款逻辑
	account, err := service.store.GetAccount(context.Background(), req.ID)
	if err != nil {
		return err
	}
	if account.Balance > 0 {
		req := TxRequest{
			FromAccountID: account.ID,
			ToAccountID:   account.ID,
			Amount:        -account.Balance,
			Currency:      account.Currency,
		}
		_, err := service.CreateTransfer(req)
		if err != nil {
			return err
		}
	}

	return nil
}

type ListAccountRequest struct {
	ID       int64 `form:"id"`
	PageID   int32 `form:"page_id"`
	PageSize int32 `form:"page_size" `
}

func (service *AccountsService) ListAccounts(req ListAccountRequest) ([]repository.Account, error) {

	arg := repository.ListAccountsParams{
		ID:     req.ID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	account, err := service.store.ListAccounts(context.Background(), arg)
	if err != nil {
		return account, err
	}
	return account, nil
}

type ListEntriesRequest struct {
	ID       int64 `form:"id"`
	PageID   int32 `form:"page_id" `
	PageSize int32 `form:"page_size"`
}

func (service *AccountsService) ListEntries(req ListAccountRequest) ([]repository.Entry, error) {

	arg := repository.ListEntriesParams{
		AccountID: req.ID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
	}
	account, err := service.store.ListEntries(context.Background(), arg)
	if err != nil {
		return account, err
	}
	return account, nil
}
