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

	"github.com/Abdillah-Epi/fizz-buzz/internal/router"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	router := router.New()

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
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
