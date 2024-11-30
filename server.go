//go:build prod
// +build prod

package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

func (s *ApiServer) Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	serverLogger := slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}), slog.LevelDebug)

	router := s.InitRoutes()

	loggingMiddleware := LoggingMiddleware(logger)
	loggedRouter := loggingMiddleware(router)

	httpServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         ":http",
		Handler:      loggedRouter,
		ErrorLog:     serverLogger,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error("Failed to start HTTP server", err)
		os.Exit(1)
	}
}
