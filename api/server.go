package api

import (
	"fmt"
	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	"github.com/Petatron/bank-simulator-backend/db/util"
	"github.com/Petatron/bank-simulator-backend/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w ", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	// Set up currency validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, nil
		}
	}

	server.setupRouter()

	return server, nil
}

// setupRouter sets up all the routes for the HTTP server.
func (server *Server) setupRouter() {
	route := gin.Default()

	route.POST("/users", server.createUser)
	route.POST("/users/login", server.loginUser)

	authRoutes := route.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.DELETE("/accounts/:id", server.deleteAccount)
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = route
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse creates a gin.H with a single "error" field containing the error message
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
