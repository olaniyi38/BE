package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/olaniyi38/BE/db/sqlc"
	"github.com/olaniyi38/BE/token"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respondWithErr(http.StatusBadRequest, err, "invalid request", ctx)
		return
	}

	payload, ok := ctx.MustGet(auth_payload_key).(*token.Payload)
	if !ok {
		respondWithErr(http.StatusUnauthorized, fmt.Errorf("bad request"), "no auth session found", ctx)
		return
	}

	createParams := db.CreateAccountParams{
		Name:     payload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	newAccount, err := server.store.CreateAccount(ctx, createParams)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				respondWithErr(http.StatusForbidden, err, "account creation error", ctx)
				return
			}
		}
		respondWithErr(http.StatusInternalServerError, err, "account creation error", ctx)
		return
	}

	ctx.JSON(http.StatusOK, newAccount)

}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(c *gin.Context) {
	var request GetAccountRequest

	if err := c.BindUri(&request); err != nil {
		respondWithErr(http.StatusBadRequest, err, "error parsing uri", c)
		return
	}

	account, err := server.store.GetAccount(c, request.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithErr(http.StatusNotFound, errors.New("account not found"), fmt.Sprintf("the account with id %v was not found", request.ID), c)
			return
		}
		respondWithErr(http.StatusInternalServerError, err, "error getting account", c)
		return
	}

	payload := c.MustGet(auth_payload_key).(*token.Payload)
	if payload.Username != account.Name {
		respondWithErr(http.StatusForbidden, fmt.Errorf("forbidden"), "you are unauthorized to perform this action", c)
		return
	}

	c.JSON(http.StatusOK, account)
}

type ListAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10"`
}

func (server *Server) listAccounts(c *gin.Context) {
	var request ListAccountsRequest

	if err := c.BindQuery(&request); err != nil {
		respondWithErr(http.StatusBadRequest, err, "invalid request", c)
		return
	}

	payload := c.MustGet(auth_payload_key).(*token.Payload)

	//allows the client define how many records per page
	offset := (request.PageID - 1) * request.PageSize
	accounts, err := server.store.ListAccounts(c, db.ListAccountsParams{
		Offset: offset,
		Limit:  request.PageSize,
		Name:   payload.Username,
	})

	if err != nil {
		respondWithErr(http.StatusInternalServerError, err, "error listing accounts", c)
		return
	}

	c.JSON(200, accounts)
}
