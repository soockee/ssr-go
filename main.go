package main

import (
	"net/http"
	"os"

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

	isProd := false
	if os.Getenv("DEPLOYMENT_ENVIRONMENT") == "prod" {
		isProd = true
	}
	log.Info().Msgf("IsProd: %t", isProd)

	log.Info().Msg("Starting Server")
	server := NewApiServer("0.0.0.0:80", store, fs, isProd)
	server.Run()
}
