package main

import (
	"flag"
	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting...")
	log.Info().Msg("Setup Storage...")
	store, err := NewSQLiteStore()
	if err != nil {
		log.Fatal().AnErr("error", err).Str("message", "database error").Send()
	}

	log.Info().Msg("Setup FileServer...")
	fs := http.FileServer(http.Dir("./assets"))

	var isProd bool
	flag.BoolVar(&isProd, "isProd", false, "Set to true if running in production environment")
	flag.Parse()

	log.Info().Msg("Starting Server")
	server := NewApiServer(store, fs, isProd)
	server.Run()
}
