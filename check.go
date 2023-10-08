package goup

import "errors"

var ErrAmbiguousRelease = errors.New("multiple releases were selected, but checker only allows one release")

type Checker interface {
	GetLatestRelease() (*Release, error)
}

type CheckConfig struct {
	Version string
}

func Check(checker Checker, config *CheckConfig) (*Update, error) {
	panic("not implemented")
}
