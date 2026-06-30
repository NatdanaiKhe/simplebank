package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NatdanaiKhe/simplebank/api"
	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/NatdanaiKhe/simplebank/util"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	config, err := util.LoadConfig(".")
	if err != nil {
		logger.Fatal("Cannot load config", zap.Error(err))
	}

	conn, err := sql.Open(config.DB_Driver, config.DB_URL)
	if err != nil {
		logger.Fatal("Cannot connect to database", zap.Error(err))
	}
	defer conn.Close()

	store := db.NewStore(conn)
	svc := service.NewAccountService(store)
	server := api.NewServer(svc, logger)

	errChan := make(chan error, 1)
	go func() {
		if err := server.Start(config.SERVER_ADDRESS); err != nil {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Fatal("Failed to start server", zap.Error(err))
	case sig := <-quit:
		logger.Info("Received signal — initiating graceful shutdown",
			zap.String("signal", sig.String()),
		)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server stopped gracefully")
}
