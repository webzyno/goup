package goup

import (
	"github.com/stretchr/testify/require"
	"path"
	"runtime"
	"testing"
)

//goland:noinspection GoBoolExpressions
func TestApply(t *testing.T) {
	var name string
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		name = "goup-test_0.1.0_darwin_arm64"
	}
	update := &Update{
		GetFile: FromFile(path.Join("test", name)),
	}
	err := Apply(update, &ApplyConfig{})
	require.NoError(t, err)
}
