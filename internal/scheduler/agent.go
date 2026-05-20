package scheduler

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// loadAgentSystemPrompt reads an agent .md file, extracts the body
// after frontmatter as the system prompt, and returns any extension
// names from the frontmatter.
func loadAgentSystemPrompt(agent string) (systemPrompt string, extensions []string) {
	agentPath := filepath.Join(os.Getenv("HOME"), ".pi", "agent", "agents", agent+".md")
	data, err := os.ReadFile(agentPath)
	if err != nil {
		log.Printf("scheduler: agent file %s: %v", agentPath, err)
		return "", nil
	}
	content := string(data)

	// Parse frontmatter (--- delimited block at start).
	body := content
	if strings.HasPrefix(strings.TrimSpace(content), "---") {
		rest := strings.TrimSpace(content)[3:]
		end := strings.Index(rest, "\n---")
		if end >= 0 {
			frontmatter := rest[:end]
			body = strings.TrimSpace(rest[end+4:])

			// Extract extensions from frontmatter (simple key: value scan).
			for _, line := range strings.Split(frontmatter, "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "extensions:") {
					val := strings.TrimSpace(strings.TrimPrefix(line, "extensions:"))
					for _, ext := range strings.Split(val, ",") {
						e := strings.TrimSpace(ext)
						if e != "" {
							extensions = append(extensions, e)
						}
					}
				}
			}
		}
	}

	return body, extensions
}
