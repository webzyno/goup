package goreleaser

import (
	"bufio"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/webzyno/goup"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// GitHubConfig is used to configure the behavior of GitHubChecker.
type GitHubConfig struct {
	// Owner is the account owner of the repository. The name is not case-sensitive.
	Owner string

	// Repo is the name of the repository without the .git extension. The name is not case-sensitive.
	Repo string

	// BaseURL is the GitHub base API endpoint.
	// Providing a value is a requirement when working with GitHub Enterprise.
	BaseURL string

	// Token is the GitHub personal access token to send with the release checking request.
	Token string

	// Client is the HTTP client used to send HTTP request
	Client *http.Client
}

func (c *GitHubConfig) appendDefaults() {
	if c.Client == nil {
		c.Client = &http.Client{}
	}

	if c.BaseURL == "" {
		c.BaseURL = "https://api.github.com"
	}
}

type gitHubChecker struct {
	config *GitHubConfig

	client *resty.Client
}

// NewGitHubChecker returns a new GitHub GoReleaser checker.
// The config must be non-nil and must include at least Owner and Repo attributes.
func NewGitHubChecker(config *GitHubConfig) goup.Checker {
	config.appendDefaults()

	return &gitHubChecker{
		config: config,
		client: resty.NewWithClient(config.Client).
			SetBaseURL(config.BaseURL).
			SetHeader("Accept", "application/vnd.github+json").
			SetHeader("X-GitHub-Api-Version", "2022-11-28").
			SetAuthToken(config.Token),
	}
}

func (g *gitHubChecker) GetLatestRelease() (*goup.Release, error) {
	// Validate configuration
	if g.config.Owner == "" || g.config.Repo == "" {
		return nil, os.ErrInvalid
	}

	resp, err := g.client.R().
		SetPathParam("owner", g.config.Owner).
		SetPathParam("repo", g.config.Repo).
		SetResult(&githubRelease{}).
		Get("/repos/{owner}/{repo}/releases/latest")
	if err != nil {
		return nil, err
	}
	latestRelease := resp.Result().(*githubRelease)

	// Filter release assets that match the OS and Arch and checksum file
	matchedAssets := make([]githubReleaseAsset, 0)
	var checksumAsset githubReleaseAsset
	for i, asset := range latestRelease.Assets {
		if strings.Contains(strings.ToLower(asset.Name), runtime.GOOS) && strings.Contains(strings.ToLower(asset.Name), runtime.GOARCH) {
			matchedAssets = append(matchedAssets, latestRelease.Assets[i])
		}
		if strings.Contains(strings.ToLower(asset.Name), "checksums") {
			checksumAsset = latestRelease.Assets[i]
		}
	}
	if len(matchedAssets) > 1 {
		return nil, goup.ErrAmbiguousRelease
	} else if len(matchedAssets) == 0 {
		return nil, nil
	}
	asset := matchedAssets[0]

	// Download checksum if checksum exists
	var checksum []byte = nil
	if checksumAsset.Id != 0 {
		resp, err := g.client.R().
			SetDoNotParseResponse(true).
			SetHeader("Accept", "application/octet-stream").
			SetPathParam("owner", g.config.Owner).
			SetPathParam("repo", g.config.Repo).
			SetPathParam("assetId", strconv.Itoa(checksumAsset.Id)).
			Get("/repos/{owner}/{repo}/releases/assets/{assetId}")
		if err != nil {
			return nil, err
		}
		checksum, err = g.parseChecksumFile(resp.RawBody(), asset.Name)
		if err != nil {
			return nil, err
		}
	}

	return &goup.Release{
		Update: goup.Update{
			GetFile: goup.DownloadWithResty(
				fmt.Sprintf("/repos/%s/%s/releases/assets/%d", g.config.Owner, g.config.Repo, asset.Id),
				g.client,
			),
			Version:  latestRelease.TagName,
			Checksum: checksum,
			Time:     latestRelease.PublishedAt,
			Size:     uint64(asset.Size),
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Extras:   latestRelease,
		},
	}, nil
}

func (g *gitHubChecker) parseChecksumFile(body io.ReadCloser, assetName string) ([]byte, error) {
	defer body.Close()

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, "  ")
		if len(words) == 2 && words[1] == assetName {
			return []byte(words[0]), nil
		}
	}
	return nil, scanner.Err()
}

type githubRelease struct {
	Url             string               `json:"url"`
	HtmlUrl         string               `json:"html_url"`
	AssetsUrl       string               `json:"assets_url"`
	UploadUrl       string               `json:"upload_url"`
	TarballUrl      string               `json:"tarball_url"`
	ZipballUrl      string               `json:"zipball_url"`
	DiscussionUrl   string               `json:"discussion_url"`
	Id              int                  `json:"id"`
	NodeId          string               `json:"node_id"`
	TagName         string               `json:"tag_name"`
	TargetCommitish string               `json:"target_commitish"`
	Name            string               `json:"name"`
	Body            string               `json:"body"`
	Draft           bool                 `json:"draft"`
	Prerelease      bool                 `json:"prerelease"`
	CreatedAt       time.Time            `json:"created_at"`
	PublishedAt     time.Time            `json:"published_at"`
	Author          githubUser           `json:"author"`
	Assets          []githubReleaseAsset `json:"assets"`
}

type githubUser struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type githubReleaseAsset struct {
	Url                string     `json:"url"`
	BrowserDownloadUrl string     `json:"browser_download_url"`
	Id                 int        `json:"id"`
	NodeId             string     `json:"node_id"`
	Name               string     `json:"name"`
	Label              string     `json:"label"`
	State              string     `json:"state"`
	ContentType        string     `json:"content_type"`
	Size               int        `json:"size"`
	DownloadCount      int        `json:"download_count"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Uploader           githubUser `json:"uploader"`
}
