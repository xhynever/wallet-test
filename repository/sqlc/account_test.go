package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pioz/faker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/xhynever/wallet-test/util"
)

var (
	testQueries *Queries
	testDB      *sqlx.DB
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	config, err := util.LoadConfig("../..")
	if err != nil {
		logrus.Fatal("cannot load config: ", err)
	}
	testDB, err = sqlx.Open(config.DbDriver, config.PgUrl)
	if err != nil {
		logrus.Fatal("cannot connect to db", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}

func TestRepository_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestRepository_UpdateAccount(t *testing.T) {

	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Owner:   account1.Owner,
		Balance: faker.Int64InRange(555, 888),
	}
	updated, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updated)
	require.Equal(t, account1.ID, updated.ID)
	require.Equal(t, account1.Owner, updated.Owner)
	require.Equal(t, arg.Balance, updated.Balance)
	require.WithinDuration(t, account1.CreatedAt, updated.CreatedAt, time.Second)
}
func TestRepository_GetAccount(t *testing.T) {

	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestQuery_DeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	deleted, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, deleted)
}

func TestQuery_ListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createFixAccount(t)
	}
	arg := ListAccountsParams{
		Owner:  "xhy",
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

	for i := 0; i < 10; i++ {
		createFixAccount(t)
	}
	arg2 := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts2, err := testQueries.ListAccounts(context.Background(), arg2)
	require.NoError(t, err)
	require.Len(t, accounts2, 5)
	for _, account := range accounts2 {
		require.NotEmpty(t, account)
	}
}
func TestQuery_AddAccountBalance(t *testing.T) {

	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	//超额扣款，交易失败
	arg := AddAccountBalanceParams{
		Amount: -(account2.Balance + 1),
		ID:     account2.ID,
	}
	account, err := testQueries.AddAccountBalance(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, account)
	account3, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.Equal(t, account3.Balance, account2.Balance)
	require.Equal(t, account3.Owner, account2.Owner)
	require.Equal(t, account3.Currency, account2.Currency)
	require.WithinDuration(t, account3.CreatedAt, account2.CreatedAt, time.Second)

	// 账户余额变更
	arg2 := AddAccountBalanceParams{
		Amount: -(account2.Balance - 1),
		ID:     account2.ID,
	}
	account4, err := testQueries.AddAccountBalance(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, account4)
	require.Equal(t, account4.Balance, int64(1))
	require.Equal(t, account4.Owner, account2.Owner)
	require.Equal(t, account4.Currency, account2.Currency)
	require.WithinDuration(t, account4.CreatedAt, account2.CreatedAt, time.Second)

}

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    faker.FirstName(),
		Balance:  faker.Int64InRange(100, 1000),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Owner, account.Owner)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func createFixAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    "xhy",
		Balance:  0,
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Owner, account.Owner)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
