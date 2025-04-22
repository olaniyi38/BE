package api

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http/httptest"
// 	"testing"

// 	mockdb "github.com/olaniyi38/BE/db/mock"
// 	db "github.com/olaniyi38/BE/db/sqlc"
// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/mock/gomock"
// )

// func TestMakeTransfer(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	store := mockdb.NewMockStore(ctrl)

// 	account1 := randomAccount()
// 	account2 := randomAccount()
// 	amount := int64(10)

// 	params := db.TransferTxParams{
// 		FromAccountID: account1.ID,
// 		ToAccountID:   account2.ID,
// 		Amount:        amount,
// 	}

// 	store.EXPECT().TransferTX(gomock.Any(), gomock.Eq(params)).Return()
// 	body, err := json.Marshal(params)
// 	require.NoError(t, err)

// 	server := newTestServer(t, store)
// 	recorder := httptest.NewRecorder()
// 	request := httptest.NewRequest("POST", "/makeTransfer", bytes.NewReader(body))
// 	attachAuthToRequest(t, request, server.tokenMaker)

// 	server.router.ServeHTTP(recorder, request)

// 	response := db.TransferTxResult{}

// 	err = json.Unmarshal(recorder.Body.Bytes(), &response)
// }
