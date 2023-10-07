package goreleaser

import (
	"github.com/webzyno/goup"
	"net/http"
)

// GitHubConfig is used to configure the behavior of GitHubChecker.
type GitHubConfig struct {
	// Owner is the account owner of the repository. The name is not case-sensitive.
	Owner string

	// Repo is the name of the repository without the .git extension. The name is not case-sensitive.
	Repo string

	// Token is the GitHub personal access token to send with the release checking request.
	Token string
}

type gitHubChecker struct {
	client *http.Client

	config *GitHubConfig
}

// NewGitHubChecker returns a new GitHub GoReleaser checker.
// The config must be non-nil and must include at least Owner and Repo attributes.
func NewGitHubChecker(config *GitHubConfig) goup.Checker {
	return &gitHubChecker{client: new(http.Client), config: config}
}

func (g *gitHubChecker) GetLatestRelease() (goup.Release, error) {
	//TODO implement me
	panic("implement me")
}
