package api

import (
	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
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

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
