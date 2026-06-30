package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	srv    *http.Server
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.POST("/accounts", server.createAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.PUT("/accounts/:id", server.updateAccount)

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
	log.Printf("Server listening on %s", address)
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
