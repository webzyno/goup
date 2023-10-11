package goup

import (
	"crypto"
	"github.com/minio/selfupdate"
	"os"
)

// ApplyConfig is used to configure the path of executables and verification settings during apply process.
type ApplyConfig struct {
	// Path defines the path to the file to update.
	// The empty string means 'the executable file of the running program'.
	Path string

	// Mode is the file mode applied during Path replacement. If zero, defaults to 0755.
	Mode os.FileMode

	// Hash is used to generate the update checksum to match provided one.
	// If not set, SHA256 is used.
	Hash crypto.Hash

	// OldPath is where the old executable file at this path after a successful update.
	// The empty string means the old executable file will be removed after the update.
	OldPath string
}

func Apply(update *Update, config *ApplyConfig) error {
	// Download new update
	data, err := update.GetFile.Download()
	if err != nil {
		return err
	}
	defer data.Close()

	if err := selfupdate.Apply(data, selfupdate.Options{
		TargetPath:  config.Path,
		TargetMode:  config.Mode,
		Checksum:    update.Checksum,
		Hash:        config.Hash,
		OldSavePath: config.OldPath,
	}); err != nil {
		if rollbackErr := selfupdate.RollbackError(err); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	return nil
}
