package goup

import (
	"errors"
	"golang.org/x/mod/semver"
	"os"
)

var ErrAmbiguousRelease = errors.New("multiple releases were selected, but checker only allows one release")

type Checker interface {
	GetLatestUpdate() (*Update, error)
}

type CheckConfig struct {
	CurrentVersion  string
	VersionComparer Comparer
}

func Check(checker Checker, config *CheckConfig) (*Update, error) {
	if config == nil || config.CurrentVersion == "" {
		return nil, os.ErrInvalid
	}
	if config.VersionComparer == nil {
		config.VersionComparer = &SemverComparer{}
	}

	update, err := checker.GetLatestUpdate()
	if err != nil {
		return nil, err
	}

	if config.VersionComparer.Compare(update.Version, config.CurrentVersion) > 1 {
		return update, nil
	}
	return nil, nil
}

// Comparer can be used to compare two versions to determine who is latest.
type Comparer interface {
	// Compare returns an integer comparing two versions.
	// The result will be 0 if a == b, -1 if a < b, or +1 if a > b.
	//
	// An invalid version string is considered less than a valid one. All invalid version strings compare equal to each other.
	Compare(a, b string) int
}

type SemverComparer struct{}

func (s *SemverComparer) Compare(a, b string) int {
	return semver.Compare(a, b)
}
