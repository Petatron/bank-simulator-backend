package api

import (
	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(store *db.Store) *Server {
	server := &Server{
		store:  store,
		router: gin.Default(),
	}

	route := gin.Default()
	server.router = route

	return server
}
