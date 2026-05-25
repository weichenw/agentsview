package memory

import (
	"encoding/json"
	"net/http"
	"strings"

	"go.kenn.io/agentsview/internal/config"
)

// Handler serves the memory graph REST API.
type Handler struct {
	cfg *config.Config
}

// NewHandler creates a Handler with the given agents-view config.
func NewHandler(cfg *config.Config) *Handler { return &Handler{cfg: cfg} }

// Graph handles GET /api/v1/memory/graph.
func (h *Handler) Graph(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	dbs := splitCSV(q.Get("dbs"), []string{"pi"})
	types := splitCSV(q.Get("types"), []string{"session", "project", "tool", "domain", "category"})

	graph, err := BuildGraph(h.cfg, dbs, types)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(graph)
}

// Meta handles GET /api/v1/memory/meta.
func (h *Handler) Meta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Meta{DBs: AvailableDBs(h.cfg)})
}

// Debug handles GET /api/v1/memory/debug — returns raw counts per source.
func (h *Handler) Debug(w http.ResponseWriter, r *http.Request) {
	graph, err := BuildGraph(h.cfg, []string{"pi"}, []string{"session", "project", "tool", "domain", "category"})
	dbInfos := AvailableDBs(h.cfg)
	paths := make([]string, 0, len(dbInfos))
	for _, d := range dbInfos {
		paths = append(paths, d.Path)
	}

	resp := map[string]interface{}{
		"db_paths": paths,
		"error":    nil,
	}
	if err != nil {
		resp["error"] = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	sqlCount, mdCount := 0, 0
	for _, n := range graph.Nodes {
		if n.Type != "memory" {
			continue
		}
		if n.Source == "extended" {
			sqlCount++
		} else if n.Source == "core" {
			mdCount++
		}
	}
	resp["total_nodes"] = len(graph.Nodes)
	resp["total_links"] = len(graph.Links)
	resp["memory_nodes"] = sqlCount + mdCount
	resp["sqlite_memories"] = sqlCount
	resp["markdown_memories"] = mdCount
	resp["stats"] = graph.Stats
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func splitCSV(s string, fallback []string) []string {
	if s == "" {
		return fallback
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
