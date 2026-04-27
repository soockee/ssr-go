package main

import (
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("Starting...")
	slog.Info("Setup Storage...")
	store, err := NewSQLiteStore()
	if err != nil {
		slog.Any("err", err)
	}

	fs := http.FileServer(http.Dir("./assets"))

	gameProxy := NewGameProxy("soockee", "terminal-games", "./cache/games")
	if err := gameProxy.FetchLatestReleases(); err != nil {
		slog.Warn("Could not fetch game binaries from GitHub, serving from cache if available", "err", err)
	}

	server := NewApiServer(store, fs, gameProxy)
	server.Run()
}
