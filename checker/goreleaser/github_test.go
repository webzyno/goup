package goreleaser

import (
	"github.com/stretchr/testify/assert"
	"github.com/webzyno/goup"
	"io"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

func TestGitHubChecker_GetLatestRelease(t *testing.T) {
	config := &GitHubConfig{
		Owner: "webzyno",
		Repo:  "goup-test",
		Token: os.Getenv("GITHUB_TOKEN"),
	}
	checker := NewGitHubChecker(config)
	release, err := checker.GetLatestUpdate()
	assert.NoError(t, err)

	// Asset release
	releaseTime, _ := time.Parse(time.RFC3339, "2023-10-07T14:09:25Z")
	assert.Equal(t, release.Version, "v0.1.0")
	assert.Equal(t, release.Time, releaseTime)
	assert.Equal(t, release.OS, runtime.GOOS)
	assert.Equal(t, release.Arch, runtime.GOARCH)

	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		assert.Equal(t, release.Size, uint64(1432994))
		assert.Equal(t, release.Checksum, []byte("3edda6fd54d12fc5c7e615100f54d3930035bb8e7b82a6d9ea1e51c335a9f675"))
		assertDownloader(t, release.GetFile, "goup-test_0.1.0_darwin_arm64")
	} else if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		assert.Equal(t, release.Size, uint64(1253376))
		assert.Equal(t, release.Checksum, []byte("9f58cb7e7c064629d0a4830206d34be1daf595570522ba477903727639d0ce90"))
		assertDownloader(t, release.GetFile, "goup-test_0.1.0_linux_amd64")
	}
}

func assertDownloader(t *testing.T, downloader goup.Downloader, name string) {
	reader, err := downloader.Download()
	assert.NoError(t, err)
	defer reader.Close()

	content, _ := os.ReadFile(path.Join("test", name))
	downloaderContent, err := io.ReadAll(reader)
	assert.Equal(t, content, downloaderContent)
}
