package api

import (
	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}

	// Set up currency validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil
		}
	}

	route := gin.Default()
	route.POST("/accounts", server.createAccount)
	route.GET("/accounts/:id", server.getAccount)
	route.GET("/accounts", server.listAccount)
	route.DELETE("/accounts/:id", server.deleteAccount)
	route.POST("/transfers", server.createTransfer)

	server.router = route

	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse creates a gin.H with a single "error" field containing the error message
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
