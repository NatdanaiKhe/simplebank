package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NatdanaiKhe/simplebank/api"
	db "github.com/NatdanaiKhe/simplebank/db/sqlc"
	"github.com/NatdanaiKhe/simplebank/service"
	"github.com/NatdanaiKhe/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	conn, err := sql.Open(config.DB_Driver, config.DB_URL)
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	svc := service.NewAccountService(store)
	server := api.NewServer(svc)

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
		log.Fatalf("Failed to start server: %v", err)
	case sig := <-quit:
		log.Printf("Received signal %v — initiating graceful shutdown", sig)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
