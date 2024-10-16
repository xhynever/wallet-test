package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	existed = make(map[int]bool)
)

func TestQuery_TransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n conccurent transfer transactions
	n := 3
	amount := int64(10)

	results := make(chan TransferTxResult)
	errs := make(chan error)

	for i := 0; i < n; i++ {

		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfers
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		//check accounts balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	//check the final updated balance
	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	// 存款
	account3 := createRandomAccount(t)

	amount3 := int64(10)
	result3, err := store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: account3.ID,
		ToAccountID:   account3.ID,
		Amount:        amount3,
	})
	require.NoError(t, err)
	checkTransferTxResult(t, result3, account3, amount3)
	// 取款
	account4 := createRandomAccount(t)
	amount4 := int64(-10)
	result4, err := store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: account4.ID,
		ToAccountID:   account4.ID,
		Amount:        amount4,
	})
	require.NoError(t, err)
	checkTransferTxResult(t, result4, account4, amount4)

}
func checkTransferTxResult(t *testing.T, res TransferTxResult, account Account, amount int64) {
	store := NewStore(testDB)
	//check transfers
	require.NotEmpty(t, res)
	transfer := res.Transfer
	require.Equal(t, account.ID, transfer.FromAccountID)
	require.Equal(t, account.ID, transfer.ToAccountID)
	require.Equal(t, amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	_, err := store.GetTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	//check entries
	fromEntry := res.FromEntry
	require.Empty(t, fromEntry)
	require.Equal(t, int64(0), fromEntry.AccountID)
	require.Equal(t, int64(0), fromEntry.Amount)
	require.Equal(t, int64(0), fromEntry.ID)

	//
	toEntry := res.ToEntry
	require.NotEmpty(t, toEntry)
	require.Equal(t, account.ID, toEntry.AccountID)
	require.Equal(t, amount, toEntry.Amount)
	require.NotZero(t, toEntry.ID)
	require.NotZero(t, toEntry.CreatedAt)

	_, err = store.GetEntry(context.Background(), toEntry.ID)
	require.NoError(t, err)

	//check accounts
	fromAccount := res.FromAccount
	require.Equal(t, int64(0), fromAccount.ID)
	require.Equal(t, int64(0), fromAccount.Balance)

	toAccount := res.ToAccount
	require.NotEmpty(t, toAccount)
	require.Equal(t, account.ID, toAccount.ID)

	//check accounts balance
	diff := toAccount.Balance - account.Balance

	require.Equal(t, diff, amount)
	// diff2 := toAccount.Balance - account.Balance
	require.True(t, diff%amount == 0)
}

func TestQuery_TransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n conccurent transfer transactions
	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID
		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">>after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

}
