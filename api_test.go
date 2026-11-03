package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Defaults verify development builds. To verify linker injection, run after templ
// generation:
// go test -ldflags '-X main.version=v1.2.3 -X main.commit=0123456789abcdef0123456789abcdef01234567' . -args -version.expected=v1.2.3 -commit.expected=0123456789abcdef0123456789abcdef01234567
var (
	expectedVersion = flag.String("version.expected", "dev", "expected binary release version")
	expectedCommit  = flag.String("commit.expected", "unknown", "expected binary source commit")
)

func TestVersion(t *testing.T) {
	router := (&ApiServer{}).InitRoutes()
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/version", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if got := response.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}
	var got map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	want := map[string]string{"version": *expectedVersion, "commit": *expectedCommit}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("response = %v, want %v", got, want)
	}
}

func TestVersionUnsupportedMethods(t *testing.T) {
	router := (&ApiServer{}).InitRoutes()
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, http.MethodOptions} {
		t.Run(method, func(t *testing.T) {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(method, "/version", nil))
			if response.Code != http.StatusMethodNotAllowed {
				t.Errorf("status = %d, want %d", response.Code, http.StatusMethodNotAllowed)
			}
			if got := response.Header().Get("Allow"); got != http.MethodGet {
				t.Errorf("Allow = %q, want GET", got)
			}
		})
	}
}
