package services

import (
	"os"
	"testing"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pioz/faker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	db "github.com/xhynever/wallet-test/repository/sqlc"
	repository "github.com/xhynever/wallet-test/repository/sqlc"
	mockdb "github.com/xhynever/wallet-test/tests/repository/sqlc/mock"
	"github.com/xhynever/wallet-test/util"
)

var (
	testQueries *db.Queries
	testDB      *sqlx.DB
	testService Service
)

type result struct {
	error
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	config, err := util.LoadConfig("..")
	if err != nil {
		logrus.Fatal("cannot load config: ", err)
	}
	testDB, err = sqlx.Open(config.DbDriver, config.PgUrl)
	if err != nil {
		logrus.Fatal("cannot connect to db", err)
	}
	testStore := db.NewStore(testDB)
	testService = *NewService(testStore)
	os.Exit(m.Run())
}

func TestRepository_CreateAccount(t *testing.T) {
	createRandomAccount(t)
}
func createRandomAccount(t *testing.T) repository.Account {
	arg := CreateAccountRequest{
		Owner:    faker.FirstName(),
		Currency: util.RandomCurrency(),
	}
	account, err := testService.CreateAccount(arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, int64(0), account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Owner, account.Owner)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func randomAccount() db.Account {

	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func TestCreateAccountService(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		body          CreateAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, body repository.Account, err error)
	}{
		{
			name: "OK",
			body: CreateAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Currency: account.Currency,
					Balance:  0,
				}
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil).AnyTimes()
			},
			checkResponse: func(t *testing.T, body repository.Account, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "InternalError",
			body: CreateAccountRequest{
				Owner:    account.Owner,
				Currency: account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone).AnyTimes()
			},
			checkResponse: func(t *testing.T, body repository.Account, err error) {
				require.Error(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := NewService(store)
			res, err := server.CreateAccount(tc.body)
			tc.checkResponse(t, res, err)
		})
	}
}

func TestGetAccountService(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		body          GetAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			body: GetAccountRequest{
				ID: account.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil).AnyTimes()
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "NotFound",
			body: GetAccountRequest{
				ID: account.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows).AnyTimes()
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := NewService(store)
			_, err := server.GetAccount(tc.body)
			tc.checkResponse(t, err)
		})
	}
}

func TestDeleteAccountService(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		body          DeleteRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			body: DeleteRequest{
				ID: account.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.TransferTxParams{
					FromAccountID: account.ID,
					ToAccountID:   account.ID,
					Amount:        -account.Balance,
				}
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil).AnyTimes()
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},

		{
			name: "TransferTxError",
			body: DeleteRequest{
				ID: account.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil).AnyTimes()
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "GetAccountError",
			body: DeleteRequest{
				ID: account.ID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, sql.ErrTxDone).AnyTimes()
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			server := NewService(store)
			err := server.DeleteAccount(tc.body)
			tc.checkResponse(t, err)
		})
	}
}

func TestListAccountsService(t *testing.T) {
	n := 5
	accounts := make([]db.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount()
	}

	testCases := []struct {
		name          string
		query         ListAccountRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(*testing.T, error)
	}{
		{
			name: "OK",
			query: ListAccountRequest{
				PageID:   1,
				PageSize: 5,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAccountsParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},

		{
			name: "InternalError",
			query: ListAccountRequest{
				PageID:   1,
				PageSize: 5,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewService(store)
			_, err := server.ListAccounts(tc.query)
			tc.checkResponse(t, err)
		})
	}
}

func randomEntry() db.Entry {
	return db.Entry{
		ID:        util.RandomInt(1, 1000),
		AccountID: util.RandomInt(1, 1000),
		Amount:    util.RandomInt(-1001, 100),
	}
}

func TestListEntriesService(t *testing.T) {
	n := 5
	Entrys := make([]db.Entry, n)
	for i := 0; i < n; i++ {
		Entrys[i] = randomEntry()
	}

	testCases := []struct {
		name          string
		query         ListEntriesRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			query: ListEntriesRequest{
				PageID:   1,
				PageSize: 5,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListEntriesParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListEntries(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(Entrys, nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "InternalError",
			query: ListEntriesRequest{
				PageID:   1,
				PageSize: 5,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListEntries(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Entry{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewService(store)
			_, err := server.ListEntries(tc.query)
			tc.checkResponse(t, err)
		})
	}
}
