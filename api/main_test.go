package api

import (
	"os"
	"testing"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func newTestServer(t *testing.T, services *service.ServiceContainer, store db.Store, logger *zap.Logger) *Server {
	server, err := NewServer(services, store, logger)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
