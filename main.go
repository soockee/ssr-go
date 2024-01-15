package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {

	store, err := NewSQLiteStore()
	if err != nil {
		log.Fatal().AnErr("error", err).Str("message", "database error").Send()
	}

	fs := http.FileServer(http.Dir("./assets"))
	server := NewApiServer("0.0.0.0:80", store, fs)
	server.Run()
}
