package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/olaniyi38/BE/db/sqlc"
	"github.com/olaniyi38/BE/token"
)

type TransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=0"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=0"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

// api endpoint to create a transfer
func (server *Server) createTransfer(c *gin.Context) {
	var params TransferRequest

	if err := c.BindJSON(&params); err != nil {
		respondWithErr(http.StatusBadRequest, err, "invalid request made", c)
		return
	}

	fromAccount, ok := server.validateAccount(c, params.FromAccountID, params.Currency)
	if !ok {
		return
	}

	payload := c.MustGet(auth_payload_key).(*token.Payload)
	if payload.Username != fromAccount.Name {
		respondWithErr(http.StatusForbidden, fmt.Errorf("forbidden"), "from account does not belong to authenticated user", c)
		return
	}

	_, ok = server.validateAccount(c, params.ToAccountID, params.Currency)
	if !ok {
		return
	}

	store := server.store

	transferResult, err := store.TransferTX(c, db.TransferTxParams{
		FromAccountID: params.FromAccountID,
		ToAccountID:   params.ToAccountID,
		Amount:        params.Amount,
	})

	if err != nil {
		respondWithErr(http.StatusInternalServerError, err, "an error occurred while making transfer", c)
		return
	}

	c.JSON(http.StatusOK, transferResult)
}

// checks that the account exists and the currency matches the passed in currency
func (server *Server) validateAccount(c *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithErr(http.StatusNotFound, fmt.Errorf("could not complete request"), fmt.Sprintf("account of ID %v not found", accountID), c)
			return account, false
		}
		respondWithErr(http.StatusInternalServerError, err, fmt.Sprintf("error getting account of ID %v ", accountID), c)
		return account, false
	}

	if account.Currency == currency {
		return account, true
	}

	respondWithErr(http.StatusBadRequest, fmt.Errorf("could not complete transfer"), fmt.Sprintf("account %v currency mismatch with currency to send: %v vs %v", accountID, account.Currency, currency), c)
	return account, false

}
