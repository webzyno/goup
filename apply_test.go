package goup

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"path"
	"runtime"
	"testing"
)

//goland:noinspection GoBoolExpressions
func TestApply(t *testing.T) {
	var name string
	var checksum []byte
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		name = "goup-test_0.1.0_darwin_arm64"
		checksum, _ = hex.DecodeString("3edda6fd54d12fc5c7e615100f54d3930035bb8e7b82a6d9ea1e51c335a9f675")
	} else if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		name = "goup-test_0.1.0_linux_amd64"
		checksum, _ = hex.DecodeString("9f58cb7e7c064629d0a4830206d34be1daf595570522ba477903727639d0ce90")
	}
	update := &Update{
		GetFile:  FromFile(path.Join("test", name)),
		Checksum: checksum,
	}
	err := Apply(update, &ApplyConfig{})
	require.NoError(t, err)
}
