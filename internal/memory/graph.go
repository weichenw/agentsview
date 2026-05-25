package memory

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.kenn.io/agentsview/internal/config"
	"go.kenn.io/agentsview/internal/parser"

	_ "github.com/mattn/go-sqlite3"
)

// NodeType categorises nodes in the memory graph.
type NodeType string

const (
	NodeSession  NodeType = "session"
	NodeMemory   NodeType = "memory"
	NodeProject  NodeType = "project"
	NodeHub      NodeType = "hub"
	NodeDomain   NodeType = "domain"
	NodeCategory NodeType = "category"
)

// Graph holds the memory graph topology and stats.
type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
	Stats Stats  `json:"stats"`
}

// Node represents a single graph vertex.
type Node struct {
	ID     string  `json:"id"`
	Type   string  `json:"type"`
	Label  string  `json:"label"`
	Count  int     `json:"count,omitempty"`
	Time   string  `json:"time,omitempty"`
	DB     string  `json:"db,omitempty"`
	Source string  `json:"source,omitempty"`
	Raw    any     `json:"raw,omitempty"`
	R      float64 `json:"r"`
}

// Link represents a graph edge.
type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}

// Stats provides high-level graph metrics.
type Stats struct {
	Nodes    int `json:"nodes"`
	Links    int `json:"links"`
	Messages int `json:"messages"`
	Memories int `json:"memories"`
}

// Meta stores runtime metadata for the memory viewer.
type Meta struct {
	DBs []DBInfo `json:"dbs"`
}

// DBInfo describes a discoverable SQLite database.
type DBInfo struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Path  string `json:"path"`
}

// AvailableDBs returns the memory databases that exist on disk,
// derived from the agents-view agent directory configuration.
func AvailableDBs(cfg *config.Config) []DBInfo {
	var dbs []DBInfo
	if cfg != nil {
		for _, d := range cfg.AgentDirs[parser.AgentPi] {
			// Pi session dirs are like ~/.pi/agent/sessions;
			// the memory DB lives in the sibling memory/ dir.
			parent := filepath.Dir(d)
			piPath := filepath.Join(parent, "memory", "sessions.db")
			if _, err := os.Stat(piPath); err == nil {
				dbs = append(dbs, DBInfo{Key: "pi", Label: "Pi sessions.db", Path: piPath})
				break // one Pi DB is enough
			}
		}
	}
	if len(dbs) == 0 {
		// Fallback to $HOME when config not available.
		home := os.Getenv("HOME")
		if home == "" {
			var err error
			home, err = os.UserHomeDir()
			if err != nil {
				return nil
			}
		}
		piPath := filepath.Join(home, ".pi", "agent", "memory", "sessions.db")
		if _, err := os.Stat(piPath); err == nil {
			dbs = append(dbs, DBInfo{Key: "pi", Label: "Pi sessions.db", Path: piPath})
		}
	}
	return dbs
}

// BuildGraph constructs the memory graph from enabled databases.
func BuildGraph(cfg *config.Config, dbs []string, types []string) (*Graph, error) {
	enabledDBs := newSet(dbs)
	// Memory is always active — it's a memory viewer.
	activeTypes := newSet(types)
	activeTypes["memory"] = true

	nodes := make([]Node, 0)
	links := make([]Link, 0)
	nodeMap := make(map[string]*Node)
	linkSet := make(map[string]bool)
	totalMessages := 0
	totalMemories := 0

	addNode := func(n Node) *Node {
		if existing, ok := nodeMap[n.ID]; ok {
			return existing
		}
		nodeMap[n.ID] = &n
		nodes = append(nodes, n)
		return &nodes[len(nodes)-1]
	}

	addLink := func(s, t string, w int) {
		key := s + "|" + t
		if !linkSet[key] {
			linkSet[key] = true
			links = append(links, Link{Source: s, Target: t, Value: w})
		}
	}

	dbInfos := AvailableDBs(cfg)
	for _, dbInfo := range dbInfos {
		if !enabledDBs[dbInfo.Key] {
			continue
		}
		err := processPiDB(dbInfo.Path, activeTypes, addNode, addLink, &totalMessages, &totalMemories)
		if err != nil {
			return nil, fmt.Errorf("process %s: %w", dbInfo.Path, err)
		}
	}

	// Build taxonomy: hub -> domain/category -> memory.
	// Memory always active, so hub always created.
	// Domain/category are toggleable.
	memoryNodes := filterByType(nodes, NodeMemory)
	sessionNodes := filterByType(nodes, NodeSession)

	// Count domains and categories from memory nodes.
	domains := make(map[string]int)
	categories := make(map[string]int)
	for _, mn := range memoryNodes {
		if mn.Raw == nil {
			continue
		}
		d := ""
		c := ""
		switch raw := mn.Raw.(type) {
		case MdMemory:
			d = raw.Domain
			c = raw.Category
			if c == "" {
				c = raw.Target
			}
		case MemoryRaw:
			c = strOrEmpty(raw.Category)
			if c == "" {
				c = strOrEmpty(raw.Target)
			}
		}
		if d != "" {
			domains[d]++
		}
		if c != "" {
			categories[c]++
		}
	}

	// Hub always created.
	hub := addNode(Node{ID: "hub:memories", Type: "hub", Label: "Memories", Count: len(memoryNodes), R: 20})
	_ = hub

	// Domain nodes: only if domain type is active.
	if activeTypes["domain"] {
		for dname, count := range domains {
			addNode(Node{ID: "domain:" + dname, Type: "domain", Label: dname, Count: count})
			addLink("hub:memories", "domain:"+dname, 3)
		}
	}

	// Category nodes: only if category type is active.
	if activeTypes["category"] {
		for cname, count := range categories {
			addNode(Node{ID: "category:" + cname, Type: "category", Label: cname, Count: count})
			addLink("hub:memories", "category:"+cname, 2)
		}
	}

	// Memory → taxonomy links: only if those types are active.
	for _, mn := range memoryNodes {
		if mn.Raw == nil {
			continue
		}
		d := ""
		c := ""
		switch raw := mn.Raw.(type) {
		case MdMemory:
			d = raw.Domain
			c = raw.Category
			if c == "" {
				c = raw.Target
			}
		case MemoryRaw:
			c = strOrEmpty(raw.Category)
			if c == "" {
				c = strOrEmpty(raw.Target)
			}
		}
		if d != "" && activeTypes["domain"] {
			if _, ok := nodeMap["domain:"+d]; ok {
				addLink(mn.ID, "domain:"+d, 2)
			}
		}
		if c != "" && activeTypes["category"] {
			if _, ok := nodeMap["category:"+c]; ok {
				addLink(mn.ID, "category:"+c, 2)
			}
		}
	}

	// Temporal cross-links: memory → session, only if sessions are visible.
	if activeTypes["session"] {
		memDateSessLinks := make(map[string][]string)
		for _, sn := range sessionNodes {
			if sn.Time == "" || len(sn.Time) < 10 {
				continue
			}
			date := sn.Time[:10]
			memDateSessLinks[date] = append(memDateSessLinks[date], sn.ID)
		}
		for _, mn := range memoryNodes {
			if mn.Time == "" || len(mn.Time) < 10 {
				continue
			}
			date := mn.Time[:10]
			for _, sid := range memDateSessLinks[date] {
				addLink(mn.ID, sid, 2)
			}
		}
	}

	// Filter to only the requested node types.
	// Memory is always on, not toggleable.
	allTypes := newSet([]string{"session", "project", "tool", "domain", "category"})
	isFullSet := true
	for t := range allTypes {
		if !activeTypes[t] {
			isFullSet = false
			break
		}
	}

	var filteredNodes []Node
	var filteredLinks []Link

	if isFullSet {
		filteredNodes = nodes
		filteredLinks = links
	} else {
		keepTypes := newSet(nil)
		for t := range activeTypes {
			keepTypes[t] = true
		}
		// Memory and hub always visible. Project only if toggled.
		keepTypes["memory"] = true
		keepTypes["hub"] = true

		filteredNodes = filterSlice(nodes, func(n Node) bool {
			return keepTypes[n.Type]
		})

		validIDs := newSet(nil)
		for _, n := range filteredNodes {
			validIDs[n.ID] = true
		}
		filteredLinks = filterSlice(links, func(l Link) bool {
			return validIDs[l.Source] && validIDs[l.Target]
		})

		// Orphan cleanup: memory + hub always survive; project only if toggled.
		connected := newSet(nil)
		for _, l := range filteredLinks {
			connected[l.Source] = true
			connected[l.Target] = true
		}
		orphanKeep := newSet([]string{"memory", "hub"})
		if activeTypes["project"] {
			orphanKeep["project"] = true
		}
		filteredNodes = filterSlice(filteredNodes, func(n Node) bool {
			return connected[n.ID] || orphanKeep[n.Type]
		})
	}

	// Final safety net: ensure every link references an existing node.
	finalNodeIDs := newSet(nil)
	for _, n := range filteredNodes {
		finalNodeIDs[n.ID] = true
	}
	filteredLinks = filterSlice(filteredLinks, func(l Link) bool {
		return finalNodeIDs[l.Source] && finalNodeIDs[l.Target]
	})

	// Assign radii.
	maxSess := 1
	for _, n := range filteredNodes {
		if n.Type == string(NodeSession) && n.Count > maxSess {
			maxSess = n.Count
		}
	}
	for i := range filteredNodes {
		n := &filteredNodes[i]
		switch NodeType(n.Type) {
		case NodeSession:
			n.R = 8 + float64(n.Count)/float64(maxSess)*24
		case NodeProject:
			n.R = 16
		case NodeMemory:
			n.R = 11
		case NodeHub:
			n.R = 24
		case NodeDomain:
			n.R = 14 + float64(n.Count)*2
		case NodeCategory:
			n.R = 10 + float64(n.Count)*2
		default:
			n.R = 5 + float64(n.Count)*2
		}
	}

	return &Graph{
		Nodes: filteredNodes,
		Links: filteredLinks,
		Stats: Stats{
			Nodes:    len(filteredNodes),
			Links:    len(filteredLinks),
			Messages: totalMessages,
			Memories: totalMemories,
		},
	}, nil
}

func processPiDB(
	path string,
	activeTypes map[string]bool,
	addNode func(Node) *Node,
	addLink func(string, string, int),
	totalMessages, totalMemories *int,
) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	sRows, err := db.Query(`
		SELECT id, project, cwd, started_at, ended_at, message_count
		FROM sessions
	`)
	if err != nil {
		return fmt.Errorf("query sessions: %w", err)
	}
	defer sRows.Close()

	var sessionsList []struct {
		ID       string
		Project  *string
		Cwd      *string
		Started  *string
		Ended    *string
		MsgCount *int
	}
	for sRows.Next() {
		var s struct {
			ID       string
			Project  *string
			Cwd      *string
			Started  *string
			Ended    *string
			MsgCount *int
		}
		if err := sRows.Scan(&s.ID, &s.Project, &s.Cwd, &s.Started, &s.Ended, &s.MsgCount); err != nil {
			continue
		}
		sessionsList = append(sessionsList, s)
	}

	var msgCount int
	_ = db.QueryRow(`SELECT COUNT(*) FROM messages`).Scan(&msgCount)
	*totalMessages += msgCount

	// Build session + project nodes.
	for _, s := range sessionsList {
		pname := "default"
		if s.Project != nil && *s.Project != "" {
			pname = *s.Project
		}
		pNode := addNode(Node{ID: "proj:" + pname, Type: "project", Label: pname, Count: 0, DB: "pi"})
		if pNode != nil {
			pNode.Count++
		}

		lbl := s.ID
		if len(lbl) > 34 {
			lbl = lbl[:34]
		}
		count := 0
		if s.MsgCount != nil {
			count = *s.MsgCount
		}
		addNode(Node{
			ID:    "s:" + s.ID,
			Type:  "session",
			Label: lbl,
			Count: count,
			Time:  strPtr(s.Started),
			DB:    "pi",
		})
		addLink("proj:"+pname, "s:"+s.ID, 2)
	}

	// Extended (SQLite) memories.
	mRows, err := db.Query(`
		SELECT CAST(id AS TEXT) as id, project, target, category, content,
		       failure_reason, tool_state, corrected_to, created
		FROM memories
	`)
	if err != nil {
		if strings.Contains(err.Error(), "no such table") {
			// memories table may not exist in older DBs — skip gracefully.
		} else {
			return fmt.Errorf("query memories: %w", err)
		}
	} else {
		defer mRows.Close()
		for mRows.Next() {
			var m MemoryRaw
			if err := mRows.Scan(
				&m.ID, &m.Project, &m.Target, &m.Category, &m.Content,
				&m.FailureReason, &m.ToolState, &m.CorrectedTo, &m.Created); err != nil {
				continue
			}
			(*totalMemories)++
			mid := "m:" + m.ID
			label := previewLabel(m.Content, 6)
			if label == "" {
				if m.Category != nil {
					label = *m.Category
				}
				if label == "" {
					label = "memory"
				}
			}
			addNode(Node{
				ID:     mid,
				Type:   "memory",
				Label:  truncate(label, 48),
				Count:  1,
				Time:   strOrEmpty(m.Created),
				DB:     "pi",
				Raw:    m,
				Source: "extended",
			})
			if m.Project != nil && *m.Project != "" {
				// Auto-create project node for memory's project.
				addNode(Node{ID: "proj:" + *m.Project, Type: "project", Label: *m.Project, Count: 0})
				addLink("proj:"+*m.Project, mid, 1)
			}
			for _, sess := range sessionsList {
				if sess.Project != nil && m.Project != nil && *sess.Project == *m.Project {
					addLink("s:"+sess.ID, mid, 1)
				}
			}
		}
	}

	// Core (markdown) memories.
	mdDir := filepath.Dir(path)
	mdMemories := parseMdMemories(mdDir)
	for _, m := range mdMemories {
		mid := "md:" + m.ID
		label := previewLabel(m.Content, 6)
		if label == "" {
			label = m.Category
			if label == "" {
				label = m.Target
			}
		}
		addNode(Node{
			ID:     mid,
			Type:   "memory",
			Label:  truncate(label, 48),
			Count:  m.Refs,
			Time:   m.Created,
			DB:     "pi",
			Raw:    m,
			Source: "core",
		})
		addNode(Node{ID: "proj:default", Type: "project", Label: "default", Count: 0})
		addLink("proj:default", mid, 1)
	}

	return nil
}

// MemoryRaw describes a memory entry as stored in SQLite.
type MemoryRaw struct {
	ID            string         `json:"id"`
	Project       sql.NullString `json:"project,omitempty"`
	Target        sql.NullString `json:"target,omitempty"`
	Category      sql.NullString `json:"category,omitempty"`
	Content       string         `json:"content,omitempty"`
	FailureReason sql.NullString `json:"failure_reason,omitempty"`
	ToolState     sql.NullString `json:"tool_state,omitempty"`
	CorrectedTo   sql.NullString `json:"corrected_to,omitempty"`
	Created       sql.NullString `json:"created,omitempty"`
}

// MdMemory describes a markdown core memory entry.
type MdMemory struct {
	ID             string `json:"id"`
	Target         string `json:"target"`
	Source         string `json:"source"`
	Content        string `json:"content"`
	Category       string `json:"category,omitempty"`
	Created        string `json:"created,omitempty"`
	LastReferenced string `json:"last_referenced,omitempty"`
	Refs           int    `json:"refs,omitempty"`
	Domain         string `json:"domain,omitempty"`
}

var metaRegex = regexp.MustCompile(`(?s)<!--\s*(.*?)\s*-->`)
var createdRegex = regexp.MustCompile(`created=([^,\s]+)`)
var lastRegex = regexp.MustCompile(`last=([^,\s]+)`)
var refsRegex = regexp.MustCompile(`refs=(\d+)`)
var domainRegex = regexp.MustCompile(`domain=([^,\s]+)`)
var categoryRegex = regexp.MustCompile(`^\[?(failure|correction|insight|preference|convention|tool-quirk)]?`)

func parseMdMemories(dir string) []MdMemory {
	var results []MdMemory
	files := []struct {
		path   string
		target string
		source string
	}{
		{filepath.Join(dir, "MEMORY.md"), "memory", "core"},
		{filepath.Join(dir, "USER.md"), "user", "core"},
		{filepath.Join(dir, "failures.md"), "failure", "core"},
	}
	for _, f := range files {
		data, err := os.ReadFile(f.path)
		if err != nil {
			continue
		}
		entries := strings.Split(string(data), "§")
		for idx, entry := range entries {
			entry = strings.TrimSpace(entry)
			if entry == "" {
				continue
			}
			metaMatch := metaRegex.FindStringSubmatch(entry)
			var meta, content string
			if metaMatch != nil {
				meta = metaMatch[1]
				content = strings.TrimSpace(entry[:metaRegex.FindStringIndex(entry)[0]])
			} else {
				content = entry
			}
			if content == "" {
				continue
			}
			m := MdMemory{
				ID:      f.target + "-" + fmt.Sprintf("%d", idx),
				Target:  f.target,
				Source:  f.source,
				Content: content,
				Refs:    1,
			}
			if m2 := createdRegex.FindStringSubmatch(meta); m2 != nil {
				m.Created = m2[1]
			}
			if m2 := lastRegex.FindStringSubmatch(meta); m2 != nil {
				m.LastReferenced = m2[1]
			}
			if m2 := refsRegex.FindStringSubmatch(meta); m2 != nil {
				fmt.Sscanf(m2[1], "%d", &m.Refs)
			}
			if m2 := domainRegex.FindStringSubmatch(meta); m2 != nil {
				m.Domain = m2[1]
			}
			if m2 := categoryRegex.FindStringSubmatch(content); m2 != nil {
				m.Category = strings.ToLower(m2[1])
			}
			results = append(results, m)
		}
	}
	return results
}

// Helpers -------------------------------------------------------------------

func newSet(items []string) map[string]bool {
	s := make(map[string]bool)
	for _, i := range items {
		s[i] = true
	}
	return s
}

func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func filterByType(nodes []Node, typ NodeType) []Node {
	return filterSlice(nodes, func(n Node) bool { return n.Type == string(typ) })
}

func filterSlice[T any](items []T, keep func(T) bool) []T {
	out := make([]T, 0, len(items))
	for _, x := range items {
		if keep(x) {
			out = append(out, x)
		}
	}
	return out
}

func previewLabel(content string, words int) string {
	content = strings.Join(strings.Fields(content), " ")
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return ""
	}
	if len(parts) <= words {
		return content
	}
	return strings.Join(parts[:words], " ") + "…"
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}
