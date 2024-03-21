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

	server := NewApiServer(store, fs)
	server.Run()
}
