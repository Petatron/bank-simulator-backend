package api

import (
	"database/sql"
	"errors"
	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	m "github.com/Petatron/bank-simulator-backend/model"
	"github.com/Petatron/bank-simulator-backend/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

// transferRequest defines the body for transfer API request
type transferRequest struct {
	FromAccountID int64          `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64          `json:"to_account_id" binding:"required,min=1"`
	Amount        int64          `json:"amount" binding:"required,gt=0"`
	Currency      m.CurrencyType `json:"currency" binding:"required,currency"`
}

// createTransfer implements the API that creates a new transfer
func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	// createTransfer API rule: A logged-in user can only create a transfer for the accounts they own
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("the account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	_, valid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	// createTransfer API rule: A logged-in user can only create a transfer for the accounts they own
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// validAccount checks if the given account is valid(availability, currency validation)
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency m.CurrencyType) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "this account is not found"})
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != string(currency) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "account currency does not match or incorrect"})
		return account, false
	}

	return account, true
}
