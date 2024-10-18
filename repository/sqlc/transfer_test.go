package repository

import (
	"context"
	"testing"
	"time"

	"github.com/pioz/faker"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestRepository_GetTransfer(t *testing.T) {
	Transfer1 := createRandomTransfer(t)
	Transfer2, err := testQueries.GetTransfer(context.Background(), Transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, Transfer2)
	require.Equal(t, Transfer1.Amount, Transfer2.Amount)
	require.Equal(t, Transfer1.ID, Transfer2.ID)
	require.Equal(t, Transfer1.CreatedAt, Transfer2.CreatedAt)
	require.Equal(t, Transfer1.FromAccountID, Transfer2.FromAccountID)
	require.Equal(t, Transfer1.ToAccountID, Transfer2.ToAccountID)
	require.WithinDuration(t, Transfer1.CreatedAt, Transfer2.CreatedAt, time.Second)
}

// ListEntries
// func TestQuery_ListTransfer(t *testing.T) {
// 	for i := 0; i < 10; i++ {
// 		createFixTransfer(t)
// 	}
// 	arg2 := ListEntriesParams{
// 		AccountID: 2,
// 		Limit:     5,
// 		Offset:    5,
// 	}
// 	Transfers, err := testQueries.ListEntries(context.Background(), arg2)
// 	require.NoError(t, err)
// 	for _, Transfer := range Transfers {
// 		require.NotEmpty(t, Transfer)
// 	}
// }

func createRandomTransfer(t *testing.T) Transfer {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTransferParams{
		Amount:        faker.Int64InRange(1, 10),
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
	}

	Transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Transfer)
	require.Equal(t, account1.ID, Transfer.FromAccountID)
	require.Equal(t, account2.ID, Transfer.ToAccountID)
	require.Equal(t, arg.Amount, Transfer.Amount)
	require.NotZero(t, Transfer.ID)
	require.NotZero(t, Transfer.CreatedAt)

	return Transfer
}
func createFixTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTransferParams{
		Amount:        faker.Int64InRange(1, 10),
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
	}

	Transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, Transfer)
	require.Equal(t, account1.ID, Transfer.FromAccountID)
	require.Equal(t, account2.ID, Transfer.ToAccountID)
	require.Equal(t, arg.Amount, Transfer.Amount)
	require.NotZero(t, Transfer.ID)
	require.NotZero(t, Transfer.CreatedAt)

	return Transfer
}
