package services

import (
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	db "github.com/xhynever/wallet-test/repository/sqlc"
	mockdb "github.com/xhynever/wallet-test/tests/repository/sqlc/mock"
	"github.com/xhynever/wallet-test/util"
)

func TestTransferAPI(t *testing.T) {
	amount := int64(10)
	account1 := randomAccount()
	account2 := randomAccount()

	account1.Currency = util.USD
	account2.Currency = util.USD

	testCases := []struct {
		name          string
		body          TxRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK",
			body: TxRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil).AnyTimes()
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil).AnyTimes()

				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "FromAccountNotFound",
			body: TxRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, nil).AnyTimes()
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "ToAccountNotFound",
			body: TxRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil).AnyTimes()
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Account{}, nil).AnyTimes()
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "WithDraw",
			body: TxRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account1.ID,
				Amount:        -amount,
				Currency:      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(account1, nil).AnyTimes()
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(account1, nil).AnyTimes()
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).AnyTimes()
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},

		{
			name: "WithDrawErr",
			body: TxRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account1.ID,
				Amount:        -(account1.Balance + 1),
				Currency:      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, nil).AnyTimes()
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).AnyTimes()
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "DepositErr",
			body: TxRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account1.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, nil).AnyTimes()
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).AnyTimes()
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).AnyTimes()
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "Deposit",
			body: TxRequest{
				FromAccountID: account1.ID,
				ToAccountID:   account1.ID,
				Amount:        amount,
				Currency:      util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(account1, nil).AnyTimes()
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(account1, nil).AnyTimes()
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).AnyTimes()
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
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
			_, err := server.CreateTransfer(tc.body)
			tc.checkResponse(t, err)
		})
	}
}
