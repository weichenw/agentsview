package scheduler

import (
	"os"
	"path/filepath"
	"strings"
)

func loadAgentSystemPrompt(agentName string) (string, error) {
	agentPath := filepath.Join(os.Getenv("HOME"), ".pi", "agent", "agents", agentName+".md")
	data, err := os.ReadFile(agentPath)
	if err != nil {
		return "", err
	}

	content := string(data)

	// Strip YAML frontmatter (--- block at top)
	if strings.HasPrefix(content, "---") {
		end := strings.Index(content[3:], "---")
		if end != -1 {
			content = strings.TrimSpace(content[end+6:])
		}
	}

	return content, nil
}
