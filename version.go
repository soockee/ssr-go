package main

import (
	"log/slog"
	"net/http"
)

// Release builds override these values with Go linker flags.
var (
	version = "dev"
	commit  = "unknown"
)

func handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := WriteJson(w, http.StatusOK, struct {
		Version string `json:"version"`
		Commit  string `json:"commit"`
	}{Version: version, Commit: commit}); err != nil {
		slog.Error("failed to write version response", "err", err)
	}
}
