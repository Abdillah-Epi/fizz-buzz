package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Abdillah-Epi/fizz-buzz/internal/config"
	"github.com/Abdillah-Epi/fizz-buzz/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	app, err := server.New(cfg)
	if err != nil {
		logger.Error("failed to create server", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := app.Close(); err != nil {
			logger.Error("failed to close application resources", "error", err)
		}
	}()

	server := &http.Server{
		Addr:              ":8080",
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		fmt.Print(`
███████╗██╗███████╗███████╗     ██████╗ ██╗   ██╗███████╗███████╗
██╔════╝██║╚══███╔╝╚══███╔╝     ██╔══██╗██║   ██║╚══███╔╝╚══███╔╝
█████╗  ██║  ███╔╝   ███╔╝█████╗██████╔╝██║   ██║  ███╔╝   ███╔╝
██╔══╝  ██║ ███╔╝   ███╔╝ ╚════╝██╔══██╗██║   ██║ ███╔╝   ███╔╝
██║     ██║███████╗███████╗     ██████╔╝╚██████╔╝███████╗███████╗
╚═╝     ╚═╝╚══════╝╚══════╝     ╚═════╝  ╚═════╝ ╚══════╝ ╚══════╝
`)
		logger.Info("HTTP server started", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(
		sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	sig := <-sigChan

	logger.Info("shutdown signal received", "signal", sig.String())

	// Give active requests time to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("HTTP Server stopped")
}
