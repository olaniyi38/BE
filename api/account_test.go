package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/olaniyi38/BE/db/mock"
	db "github.com/olaniyi38/BE/db/sqlc"
	"github.com/olaniyi38/BE/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	userName := "sodiq"
	account := randomAccount(userName)

	testCases := []struct {
		name           string
		accountID      int64
		buildStubs     func(store *mockdb.MockStore)
		checkResponses func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				//build stub
				//this is a stub that says that, i want to call the fn on store, with arg1 as any and arg 2 as account.id
				//it confirms that the handler is passing the correct params to this value
				//and it should only be called once

				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)

			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder, account)
			},
		},
		{
			name:      "NOT FOUND",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				//build stub
				//this is a stub that says that, i want to call the fn on store, with arg1 as any and arg 2 as account.id
				//it confirms that the handler is passing the correct params to this value
				//and it should only be called once

				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				requireBodyMatchAccount(t, recorder, db.Account{})
			},
		},
		{
			name: "BAD REQUEST",
			buildStubs: func(store *mockdb.MockStore) {
				//build stub
				//this is a stub that says that, i want to call the fn on store, with arg1 as any and arg 2 as account.id
				//it confirms that the handler is passing the correct params to this value
				//and it should only be called once

				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireBodyMatchAccount(t, recorder, db.Account{})
			},
		},
		{
			name:      "INTERNAL SERVER ERROR",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				//build stub
				//this is a stub that says that, i want to call the fn on store, with arg1 as any and arg 2 as account.id
				//it confirms that the handler is passing the correct params to this value
				//and it should only be called once

				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponses: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				requireBodyMatchAccount(t, recorder, db.Account{})
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Log(fmt.Sprintf("running test case %v", testCase.name))
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			//start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			path := fmt.Sprintf("/account/%v", testCase.accountID)
			request := httptest.NewRequest("GET", path, nil)
			attachAuthToRequest(t, request, server.tokenMaker, userName)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponses(t, recorder)
		})
	}
}

func requireBodyMatchAccount(t *testing.T, recorder *httptest.ResponseRecorder, account db.Account) {
	var expected db.Account
	err := json.Unmarshal(recorder.Body.Bytes(), &expected)
	require.NoError(t, err)
	require.Equal(t, expected, account)
}

func randomAccount(name string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Name:     name,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func TestListAccounts(t *testing.T) {
	accounts := []db.Account{}
	var lastAccount db.Account
	userName := "sodiq"
	for _ = range 5 {
		lastAccount = randomAccount(userName)
		accounts = append(accounts, lastAccount)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{Limit: 10, Offset: 0, Name: lastAccount.Name})).Times(1).Return(accounts, nil)

	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/accounts?page_size=10&page_id=1", nil)
	attachAuthToRequest(t, request, server.tokenMaker, userName)
	server.router.ServeHTTP(recorder, request)
	require.Equal(t, recorder.Code, 200)
	gotAccounts := []db.Account{}

	err := json.Unmarshal(recorder.Body.Bytes(), &gotAccounts)
	require.NoError(t, err)

	require.NotEmpty(t, gotAccounts)
}
