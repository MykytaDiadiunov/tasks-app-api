package main

import (
	"context"
	"go-rest-api/config"
	"go-rest-api/config/container"
	"go-rest-api/internal/infra/database"
	"go-rest-api/internal/infra/http"
	"go-rest-api/internal/infra/logger"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

func main() {
	exitCode := 0
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.GetConfiguration()
	logger.Init(cfg)

	// Recover
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Panic("The system panicked!: %v\n", r)
			logger.Logger.Panic("Stack trace form panic: %s\n", string(debug.Stack()))
			exitCode = 1
		}
		os.Exit(exitCode)
	}()

	// Signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-c
		logger.Logger.Info("Received signal '%s', stopping... \n", sig.String())
		cancel()
		logger.Logger.Info("Sent cancel to all threads...")
	}()

	err := database.Migrate(cfg)
	if err != nil {
		logger.Logger.Error("Unable to apply migrations: %q\n", err)
	}

	cont := container.New()

	err = http.Server(
		ctx,
		http.CreateRouter(cont),
	)

	if err != nil {
		logger.Logger.Error("http server error: %s", err)
		exitCode = 2
		return
	}
}
