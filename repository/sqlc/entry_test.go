package repository

import (
	"context"
	"testing"
	"time"

	"github.com/pioz/faker"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreateEntry(t *testing.T) {

	createRandomEntry(t)
}

func TestRepository_GetEntry(t *testing.T) {
	Entry1 := createRandomEntry(t)
	Entry2, err := testQueries.GetEntry(context.Background(), Entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, Entry2)
	require.Equal(t, Entry1.Amount, Entry2.Amount)
	require.Equal(t, Entry1.ID, Entry2.ID)
	require.Equal(t, Entry1.CreatedAt, Entry2.CreatedAt)
	require.Equal(t, Entry1.AccountID, Entry2.AccountID)
	require.WithinDuration(t, Entry1.CreatedAt, Entry2.CreatedAt, time.Second)
}

func TestQuery_ListEntry(t *testing.T) {
	for i := 0; i < 10; i++ {
		createFixEntry(t)
	}
	arg1 := ListEntriesParams{
		AccountID: 2,
		Limit:     5,
		Offset:    5,
	}
	Entrys1, err := testQueries.ListEntries(context.Background(), arg1)
	require.NoError(t, err)
	for _, Entry := range Entrys1 {
		require.NotEmpty(t, Entry)
	}

	// 无AccountID查询
	for i := 0; i < 10; i++ {
		createFixEntry(t)
	}
	arg2 := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}
	Entrys2, err := testQueries.ListEntries(context.Background(), arg2)
	require.NoError(t, err)
	for _, Entry := range Entrys2 {
		require.NotEmpty(t, Entry)
	}

}

func createRandomEntry(t *testing.T) Entry {
	account1 := createRandomAccount(t)
	accountID := account1.ID
	arg := CreateEntryParams{
		Amount:    faker.Int64InRange(-100, 100),
		AccountID: accountID,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}
func createFixEntry(t *testing.T) Entry {
	arg := CreateEntryParams{
		Amount:    faker.Int64InRange(-100, 100),
		AccountID: 2,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.Amount, entry.Amount)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}
