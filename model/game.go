package model

import "sort"

// GameEntry represents a game discovered from GitHub releases.
type GameEntry struct {
	Slug     string // route slug, e.g. "pong"
	Name     string // display name, e.g. "PONG"
	WasmFile string // cached filename, e.g. "pong.wasm"
}

// NavLinks builds the navigation link list from the discovered games.
func NavLinks(games []GameEntry) []struct{ Name, URL string } {
	sorted := make([]GameEntry, len(games))
	copy(sorted, games)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Slug < sorted[j].Slug
	})

	links := []struct{ Name, URL string }{
		{"Home", "/"},
	}
	for _, g := range sorted {
		links = append(links, struct{ Name, URL string }{g.Name, "/games/" + g.Slug})
	}
	links = append(links, struct{ Name, URL string }{"Eberstadt", "/eberstadt/event"})
	return links
}
