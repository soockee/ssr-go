package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/soockee/ssr-go/model"
)

// GameProxy discovers games from GitHub Releases (one release per game,
// tagged as v{semver}-{game}), downloads the .tar.gz asset, extracts
// the .wasm binary, and serves it from a local cache.
type GameProxy struct {
	owner    string
	repo     string
	cacheDir string

	mu    sync.RWMutex
	games map[string]model.GameEntry // slug → entry

	client *http.Client
}

// tag pattern: v{major}.{minor}.{patch}-{game}
var tagPattern = regexp.MustCompile(`^v\d+\.\d+\.\d+-(.+)$`)

func NewGameProxy(owner, repo, cacheDir string) *GameProxy {
	return &GameProxy{
		owner:    owner,
		repo:     repo,
		cacheDir: cacheDir,
		games:    make(map[string]model.GameEntry),
		client:   &http.Client{},
	}
}

// --- GitHub Release types (only what we need) -----------------------

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// --- Public API -----------------------------------------------------

// Games returns the list of discovered games.
func (gp *GameProxy) Games() []model.GameEntry {
	gp.mu.RLock()
	defer gp.mu.RUnlock()

	entries := make([]model.GameEntry, 0, len(gp.games))
	for _, e := range gp.games {
		entries = append(entries, e)
	}
	return entries
}

// HasGame reports whether a game with the given slug exists.
func (gp *GameProxy) HasGame(slug string) (model.GameEntry, bool) {
	gp.mu.RLock()
	defer gp.mu.RUnlock()
	e, ok := gp.games[slug]
	return e, ok
}

// FetchLatestReleases lists all releases, picks the latest per game
// (releases are returned newest-first by GitHub), downloads the
// .tar.gz asset, and extracts the .wasm binary into cacheDir.
func (gp *GameProxy) FetchLatestReleases() error {
	releases, err := gp.listReleases()
	if err != nil {
		return err
	}

	// First-seen per slug wins — GitHub returns newest first.
	latestPerGame := make(map[string]githubRelease)
	for _, r := range releases {
		slug := gameSlugFromTag(r.TagName)
		if slug == "" {
			continue
		}
		if _, exists := latestPerGame[slug]; !exists {
			latestPerGame[slug] = r
		}
	}

	if err := os.MkdirAll(gp.cacheDir, 0o750); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}

	discovered := make(map[string]model.GameEntry)

	for slug, release := range latestPerGame {
		for _, asset := range release.Assets {
			if !strings.HasSuffix(asset.Name, ".tar.gz") {
				continue
			}
			// Skip checksum files
			if strings.HasSuffix(asset.Name, ".md5") {
				continue
			}

			wasmFile, err := gp.downloadAndExtract(asset)
			if err != nil {
				slog.Error("Failed to download/extract asset",
					"name", asset.Name, "err", err)
				continue
			}

			discovered[slug] = model.GameEntry{
				Slug:     slug,
				Name:     strings.ToUpper(slug),
				WasmFile: wasmFile,
			}
			slog.Info("Cached game",
				"slug", slug, "wasm", wasmFile, "tag", release.TagName)
		}
	}

	gp.mu.Lock()
	gp.games = discovered
	gp.mu.Unlock()

	slog.Info("Game discovery complete", "games", len(discovered))
	return nil
}

// ServeHTTP serves a cached .wasm file.
func (gp *GameProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := filepath.Base(r.URL.Path)

	if !strings.HasSuffix(name, ".wasm") {
		http.NotFound(w, r)
		return
	}

	name = filepath.Base(filepath.Clean(name))
	path := filepath.Join(gp.cacheDir, name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/wasm")
	http.ServeFile(w, r, path)
}

// --- internals ------------------------------------------------------

func (gp *GameProxy) listReleases() ([]githubRelease, error) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/releases?per_page=100",
		gp.owner, gp.repo,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := gp.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var releases []githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("decode releases: %w", err)
	}

	slog.Info("Fetched releases", "count", len(releases))
	return releases, nil
}

func gameSlugFromTag(tag string) string {
	m := tagPattern.FindStringSubmatch(tag)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

const maxWasmSize = 100 * 1024 * 1024 // 100 MB

func (gp *GameProxy) downloadAndExtract(asset githubAsset) (string, error) {
	resp, err := gp.client.Get(asset.BrowserDownloadURL)
	if err != nil {
		return "", fmt.Errorf("download %s: %w", asset.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download %s: HTTP %d", asset.Name, resp.StatusCode)
	}

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gzip reader: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("tar read: %w", err)
		}

		if !strings.HasSuffix(hdr.Name, ".wasm") {
			continue
		}

		wasmName := filepath.Base(hdr.Name)
		dest := filepath.Join(gp.cacheDir, wasmName)

		// Atomic write via temp file.
		tmp, err := os.CreateTemp(gp.cacheDir, ".extract-*")
		if err != nil {
			return "", err
		}
		tmpName := tmp.Name()

		limited := io.LimitReader(tr, maxWasmSize)
		if _, err := io.Copy(tmp, limited); err != nil {
			tmp.Close()
			os.Remove(tmpName)
			return "", err
		}
		if err := tmp.Close(); err != nil {
			os.Remove(tmpName)
			return "", err
		}

		if err := os.Rename(tmpName, dest); err != nil {
			os.Remove(tmpName)
			return "", err
		}

		return wasmName, nil
	}

	return "", fmt.Errorf("no .wasm file found in %s", asset.Name)
}
