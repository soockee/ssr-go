package main

import (
	"flag"
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

	var isProd bool
	flag.BoolVar(&isProd, "isProd", false, "Set to true if running in production environment")
	flag.Parse()

	server := NewApiServer(store, fs, isProd)
	server.Run()
}
