package api

import (
	"os"
	"testing"

	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func newTestServer(t *testing.T, services *service.ServiceContainer, logger *zap.Logger) *Server {
	server := NewServer(services, logger)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
