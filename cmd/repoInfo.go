package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// Repository represents a GitHub repository.
type Repository struct {
	ID              int64     `json:"id"`
	NodeID          string    `json:"node_id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Private         bool      `json:"private"`
	Owner           Owner     `json:"owner"`
	HTMLURL         string    `json:"html_url"`
	Description     string    `json:"description"`
	Fork            bool      `json:"fork"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Size            int       `json:"size"`
	StargazersCount int       `json:"stargazers_count"`
	WatchersCount   int       `json:"watchers_count"`
	Language        string    `json:"language"`
	ForksCount      int       `json:"forks_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	DefaultBranch   string    `json:"default_branch"`
	Visibility      string    `json:"visibility"`
}

// Owner represents the owner of a repository.
type Owner struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
}

func getRepositories(ctx context.Context, username string) ([]Repository, error) {
	url := fmt.Sprintf("%s/users/%s/repos?per_page=100&sort=updated", githubBaseURL, username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status %d: %s", resp.StatusCode, resp.Status)
	}

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return repos, nil
}

// repoInfoCmd represents the repoInfo command
var repoInfoCmd = &cobra.Command{
	Use:   "repos [username]",
	Short: "Get GitHub user repositories",
	Long: `Fetch and display information about all public repositories of a GitHub user.

Example:
  usertracker repos octocat`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		repos, err := getRepositories(ctx, username)
		if err != nil {
			return fmt.Errorf("failed to get repositories: %w", err)
		}

		if len(repos) == 0 {
			fmt.Printf("User %s has no public repositories\n", username)
			return nil
		}

		fmt.Printf("Found %d repositories for user %s:\n\n", len(repos), username)

		for i, repo := range repos {
			fmt.Printf("%d. %s\n", i+1, repo.FullName)
			if repo.Description != "" {
				fmt.Printf("   Description: %s\n", repo.Description)
			}
			if repo.Language != "" {
				fmt.Printf("   Language: %s\n", repo.Language)
			}
			fmt.Printf("   Stars: %d | Forks: %d | Open Issues: %d\n",
				repo.StargazersCount, repo.ForksCount, repo.OpenIssuesCount)
			fmt.Printf("   URL: %s\n", repo.HTMLURL)
			fmt.Printf("   Updated: %s\n", repo.UpdatedAt.Format("2006-01-02 15:04"))
			if repo.Fork {
				fmt.Printf("   [FORK]\n")
			}
			if repo.Private {
				fmt.Printf("   [PRIVATE]\n")
			}
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(repoInfoCmd)
}
