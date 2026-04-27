package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"

	"github.com/soockee/ssr-go/components"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string
}

type ApiServer struct {
	store      Storage
	fs         http.Handler
	gameProxy  *GameProxy
	domainName string
}

func NewApiServer(store Storage, fs http.Handler, gameProxy *GameProxy) *ApiServer {
	server := &ApiServer{
		store:     store,
		fs:        fs,
		gameProxy: gameProxy,
	}

	server.domainName = "stockhause.info"
	return server
}

func (s *ApiServer) InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", makeHTTPHandleFunc(s.handleHome))
	router.HandleFunc("/games/{id}", makeHTTPHandleFunc(s.handleGames))
	router.PathPrefix("/assets/games/").Handler(http.StripPrefix("/assets/games/", s.gameProxy))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", s.fs))

	router.HandleFunc("/eberstadt/event", makeHTTPHandleFunc(s.handleEberstadtEvent))
	return router
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return cors(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			if writeErr := WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()}); writeErr != nil {
				slog.Error("failed to write JSON error response", "err", writeErr)
			}
		}
	})
}

func (s *ApiServer) handleHome(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		component := components.Home(s.gameProxy.Games())
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)
		return nil
	case "POST":
		return WriteJson(w, http.StatusNotImplemented, "")
	default:
		return errors.New("method not allowed")
	}
}

func (s *ApiServer) handleGames(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetGames(w, r)
	default:
		return errors.New("game not found")
	}
}

func (s *ApiServer) handleGetGames(w http.ResponseWriter, r *http.Request) error {
	slug := mux.Vars(r)["id"]
	game, ok := s.gameProxy.HasGame(slug)
	if !ok {
		return errors.New("game not found")
	}

	component := components.GameLayout(game.Name, game.WasmFile, s.gameProxy.Games())
	handler := templ.Handler(component)
	handler.ServeHTTP(w, r)
	return nil
}

func (s *ApiServer) handleEberstadtEvent(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		component := components.Eberstadt()
		handler := templ.Handler(component)
		handler.ServeHTTP(w, r)
		return nil
	default:
		return errors.New("method not allowed")
	}
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		next(w, r)
	}
}
