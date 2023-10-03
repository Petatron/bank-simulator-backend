package api

import (
	db "github.com/Petatron/bank-simulator-backend/db/sqlc"
	"github.com/Petatron/bank-simulator-backend/db/util"
	"github.com/gin-gonic/gin"
	"os"
	"testing"
	"time"
)

func newTestServer(store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey: util.GetRandomStringWithLength(32),
		AccessToken:       time.Minute,
	}
	server, err := NewServer(config, store)
	if err != nil {
		panic(err)
	}

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
