package services

import (
	"context"
	"database/sql"
	"fmt"

	repository "github.com/xhynever/wallet-test/repository/sqlc"
)

type TxRequest struct {
	FromAccountID int64  `json:"from_account_id"  binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required"`
	Currency      string `json:"currency" binding:"required"`
}

func (service *AccountsService) validAccount(ctx context.Context, accountID int64, currency string) (bool, error) {
	account, err := service.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, sql.ErrNoRows
		}
		return false, err
	}
	if account.Currency != currency {
		err := fmt.Errorf("currency mismatch")
		return false, err
	}
	return true, nil
}

func (service *AccountsService) CreateTransfer(req TxRequest) (repository.TransferTxResult, error) {
	if req.FromAccountID == req.ToAccountID {
		if req.Amount > 0 {
			// 存款
			toAccount, err := service.validAccount(context.Background(), req.ToAccountID, req.Currency)
			if err != nil || !toAccount {
				return repository.TransferTxResult{}, err
			}
		} else {
			//取款
			fromAccount, err := service.validAccount(context.Background(), req.FromAccountID, req.Currency)
			if err != nil || !fromAccount {
				return repository.TransferTxResult{}, err
			}
			account, err := service.store.GetAccount(context.Background(), req.FromAccountID)
			if err != nil {
				return repository.TransferTxResult{}, err
			}
			if account.Balance < -req.Amount {
				err := fmt.Errorf("账户余额不足")
				return repository.TransferTxResult{}, err
			}
		}
	} else {
		// 不同账户，转账
		fromAccount, err := service.validAccount(context.Background(), req.FromAccountID, req.Currency)
		if err != nil || !fromAccount {
			return repository.TransferTxResult{}, err
		}
		toAccount, err := service.validAccount(context.Background(), req.ToAccountID, req.Currency)
		if err != nil || !toAccount {
			return repository.TransferTxResult{}, err
		}
	}
	// fmt.Println("创建交易tx结束", req.FromAccountID)
	arg := repository.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := service.store.TransferTx(context.Background(), arg)
	if err != nil {
		return result, err
	}
	return result, nil
}
