package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme/autocert"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string
}

type ApiServer struct {
	listenAddr    string
	store         Storage
	fs            http.Handler
	isProd        bool
	domainName    string
	ssl_cache_dir string
	owner_email   string
}

func NewApiServer(listenAddr string, store Storage, fs http.Handler, isProd bool) *ApiServer {
	server := &ApiServer{
		listenAddr: listenAddr,
		store:      store,
		fs:         fs,
		isProd:     isProd,
	}

	server.domainName = os.Getenv("DOMAIN_NAME")
	server.ssl_cache_dir = os.Getenv("SSL_CACHE_DIR")
	server.owner_email = os.Getenv("OWNER_EMAIL")
	return server
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return CORS(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	})
}

func (s *ApiServer) Run() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	router := mux.NewRouter()
	router.HandleFunc("/", makeHTTPHandleFunc(s.handleHome))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", s.fs))

	loggingMiddleware := LoggingMiddleware(logger)
	loggedRouter := loggingMiddleware(router)

	var certManager *autocert.Manager
	if s.isProd {
		certManager = &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(s.domainName),
			Cache:      autocert.DirCache(s.ssl_cache_dir),
			Email:      s.owner_email,
		}

		httpsServer := &http.Server{
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  120 * time.Second,
			Addr:         ":https",
			TLSConfig:    certManager.TLSConfig(),
			Handler:      loggedRouter,
		}

		go func() {
			logger.Info().Msg("Starting HTTPS sever")
			if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
				logger.Error().Err(err).Send()
				os.Exit(1)
			}
		}()
	}

	httpServer := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         s.listenAddr,
		Handler:      loggedRouter,
	}

	if certManager != nil {
		// allow autocert handle Let's Encrypt auth callbacks over HTTP.
		httpServer.Handler = certManager.HTTPHandler(nil)
	}

	if err := httpServer.ListenAndServe(); err != nil {
		logger.Error().Err(err)
		os.Exit(1)
	}
}

func (s *ApiServer) handleHome(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		component := Home()
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)
		return nil
	case "POST":
		WriteJson(w, http.StatusNotImplemented, "")
	default:
		return errors.New("method not allowed")
	}
	return nil
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		next(w, r)
	}
}
