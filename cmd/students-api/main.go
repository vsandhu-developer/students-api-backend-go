package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vsandhu-developer/students-api-backend-go/internal/config"
	"github.com/vsandhu-developer/students-api-backend-go/internal/http/handlers/student"
	"github.com/vsandhu-developer/students-api-backend-go/internal/storage/sqlite"
)

func main() {
	// load config

	cfg := config.MustLoad()

	// database setup

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env))

	// setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))

	// setup server

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	slog.Info("Server Started", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()

		if err != nil {
			log.Fatal("Failed to start server")

		}
	}()

	<-done

	slog.Info("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown the server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

}
