//go:build prod
// +build prod

package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func (s *ApiServer) Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	serverLogger := slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}), slog.LevelDebug)

	router := s.InitRoutes()

	loggingMiddleware := LoggingMiddleware(logger)
	loggedRouter := loggingMiddleware(router)

	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(s.domainName),
		Cache:      autocert.DirCache("/certs"),
	}

	httpsServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         ":https",
		TLSConfig:    certManager.TLSConfig(),
		Handler:      loggedRouter,
		ErrorLog:     serverLogger,
	}

	logger.Info("Starting HTTPS sever")
	if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
		logger.Error("Failed to start HTTPS server", err)
		os.Exit(1)
	}
}
