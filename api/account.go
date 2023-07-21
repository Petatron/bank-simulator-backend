package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// createAccountRequest defines the body for createAccount API request
type createAccount struct {
	owner    string `json:"owner" binding:"required"`
	currency string `json:"currency" binding:"required, oneof=USD EUR CAD"`
}

// createAccount creates a new account
func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccount
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

// errorResponse creates a gin.H with a single "error" field containing the error message
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
