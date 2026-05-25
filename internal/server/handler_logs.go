package server

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// maxLogLines is the hard upper bound to avoid memory pressure.
const maxLogLines = 5000

// logFile describes a log file exposed by the /api/v1/logs endpoint.
type logFile struct {
	Name       string   `json:"name"`
	Path       string   `json:"path"`
	Lines      []string `json:"lines"`
	TotalLines int      `json:"total_lines"`
}

// handleGetLogs serves GET /api/v1/logs.
// Returns the last N lines of each known agentsview log file.
// Query params:
//   lines  — number of lines from the end (default 500, max 5000)
func (s *Server) handleGetLogs(w http.ResponseWriter, r *http.Request) {
	lines := 500
	if v := r.URL.Query().Get("lines"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			if n > maxLogLines {
				n = maxLogLines
			}
			lines = n
		}
	}

	dataDir := s.dataDir
	files := []logFile{
		{Name: "debug.log", Path: filepath.Join(dataDir, "debug.log")},
		{Name: "server.log", Path: filepath.Join(dataDir, "logs", "server.log")},
		{Name: "server-error.log", Path: filepath.Join(dataDir, "logs", "server-error.log")},
	}

	for i := range files {
		files[i].Lines, files[i].TotalLines = tailLines(files[i].Path, lines)
	}

	writeJSON(w, http.StatusOK, map[string]any{"files": files})
}

// tailLines reads the last n lines from a file.
// If the file does not exist or cannot be read, an empty slice and 0 are returned.
func tailLines(path string, n int) ([]string, int) {
	f, err := os.Open(path)
	if err != nil {
		return []string{}, 0
	}
	defer f.Close()

	// Two-pass: count total lines, then keep last n.
	// The file is capped at 10 MB, so this is fast.
	scanner := bufio.NewScanner(f)
	var total int
	for scanner.Scan() {
		total++
	}
	if err := scanner.Err(); err != nil {
		return []string{fmt.Sprintf("error reading log: %v", err)}, total
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return []string{fmt.Sprintf("error rewinding log: %v", err)}, total
	}

	scanner = bufio.NewScanner(f)
	skip := 0
	if total > n {
		skip = total - n
	}

	var out []string
	idx := 0
	for scanner.Scan() {
		if idx >= skip {
			out = append(out, scanner.Text())
		}
		idx++
	}
	if err := scanner.Err(); err != nil {
		return []string{fmt.Sprintf("error reading log: %v", err)}, total
	}

	return out, total
}
