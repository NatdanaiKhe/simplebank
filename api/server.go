package api

import (
	"context"
	"fmt"
	"net/http"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Server struct {
	services *service.ServiceContainer
	store    db.Store
	router   *gin.Engine
	srv      *http.Server
	logger   *zap.Logger
}

func NewServer(svc *service.ServiceContainer, store db.Store, logger *zap.Logger) (*Server, error) {
	server := &Server{
		services: svc,
		store:    store,
		logger:   logger,
	}

	router := gin.Default()
	router.Use(RequestID())
	router.Use(LoggerMiddleware(logger))
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(RequestID())
	router.Use(LoggerMiddleware(server.logger))
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	apiRouter := router.Group("/api/v1")

	accountRouter := apiRouter.Group("/accounts")
	accountRouter.GET("/:id", server.getAccount)
	accountRouter.GET("", server.listAccounts)
	accountRouter.POST("", server.createAccount)
	accountRouter.DELETE("/:id", server.deleteAccount)
	accountRouter.PUT("/:id", server.updateAccount)

	userRouter := apiRouter.Group("/users")
	userRouter.POST("", server.createUser)

	transferRouter := apiRouter.Group("/transfers")
	transferRouter.POST("", server.createTransfer)

	server.router = router
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
