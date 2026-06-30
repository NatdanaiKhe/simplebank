package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	service service.AccountService
	router  *gin.Engine
	srv     *http.Server
	logger  *zap.Logger
}

func NewServer(svc service.AccountService, logger *zap.Logger) *Server {
	server := &Server{
		service: svc,
		logger:  logger,
	}

	router := gin.Default()
	router.Use(RequestID())
	router.Use(LoggerMiddleware(logger))

	apiRouter := router.Group("/api/v1")

	accountRouter := apiRouter.Group("/accounts")
	accountRouter.GET("/:id", server.getAccount)
	accountRouter.GET("", server.listAccounts)
	accountRouter.POST("", server.createAccount)
	accountRouter.DELETE("/:id", server.deleteAccount)
	accountRouter.PUT("/:id", server.updateAccount)

	server.router = router

	return server
}

// Start begins listening on the given address. It blocks until the server
// is stopped via Shutdown or a fatal error occurs.
func (server *Server) Start(address string) error {
	server.srv = &http.Server{
		Addr:    address,
		Handler: server.router,
	}
	server.logger.Info("server starting", zap.String("address", address))
	if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server stopped unexpectedly: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server, waiting for in-flight requests to
// complete or until the context is cancelled.
func (server *Server) Shutdown(ctx context.Context) error {
	if server.srv == nil {
		return nil
	}
	return server.srv.Shutdown(ctx)
}
