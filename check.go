package goup

type Checker interface {
	GetLatestRelease() (Release, error)
}

type CheckConfig struct {
	Version string
}

func Check(checker Checker, config *CheckConfig) (Update, error) {
	panic("not implemented")
}
