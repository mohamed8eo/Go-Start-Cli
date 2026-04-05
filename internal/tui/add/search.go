package add

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
)

type searchResultMsg struct {
	packages []Package
	err      error
}

type triggerSearchMsg struct{ query string }

func searchPackages(query string) tea.Cmd {
	return func() tea.Msg {
		if query == "" {
			return searchResultMsg{packages: []Package{}}
		}

		apiURL := fmt.Sprintf(
			"https://api.github.com/search/repositories?q=%s+language:go&per_page=20",
			url.QueryEscape(query),
		)

		client := &http.Client{Timeout: 8 * time.Second}
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return searchResultMsg{err: err}
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "gostart-cli/1.0")

		if token := os.Getenv("GITHUB_TOKEN"); token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		resp, err := client.Do(req)
		if err != nil {
			return searchResultMsg{err: err}
		}
		defer resp.Body.Close()

		var result struct {
			Items []struct {
				FullName    string `json:"full_name"`
				Description string `json:"description"`
				Stars       int    `json:"stargazers_count"`
			} `json:"items"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return searchResultMsg{err: fmt.Errorf("failed to parse response: %w", err)}
		}

		packages := make([]Package, 0, len(result.Items))
		for _, item := range result.Items {
			packages = append(packages, Package{
				Path:        "github.com/" + item.FullName,
				Description: item.Description,
				Version:     fmt.Sprintf("★ %d", item.Stars),
			})
		}

		return searchResultMsg{packages: packages}
	}
}
