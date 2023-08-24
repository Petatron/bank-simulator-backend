package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	model "github.com/Petatron/bank-simulator-backend/model"
	"github.com/gin-gonic/gin"
)

// createAccountRequest defines the body for createAccount API request
type createAccountRequest struct {
	Owner    string             `json:"owner" binding:"required"`
	Currency model.CurrencyType `json:"currency" binding:"required,currency"`
}

// createAccount implement the API that creates a new account
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !req.Currency.IsValid() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid currency"})
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: string(req.Currency),
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// getAccountRequest defines the body for getAccount API request
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// getAccount implements the API that returns a specific account based on the given ID
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// listAccountRequest defines the body for listAccounts API request
type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// listAccount implements the API that returns a list of accounts
func (server *Server) listAccount(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

// deleteAccountRequest defines the body for deleteAccount API request
type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// deleteAccount deletes a specific account based on the given ID
func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
