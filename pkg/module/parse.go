package module

import (
	"errors"
	"strings"
)

func ParseSource(source string) (*Module, error) {
	switch {
	case strings.HasPrefix(source, "github_archive>"):
		// github_archive>suzuki-shunsuke/lintnet-example-3#v0.1.0
		m, err := parseGitHub(source[len("github_archive>"):])
		if err != nil {
			return nil, err
		}
		m.Type = "github_archive"
		return m, nil
	case strings.HasPrefix(source, "github_content>"):
		// github_content>suzuki-shunsuke/lintnet-example/toml.jsonnet#v0.1.0
		m, err := parseGitHub(source[len("github_content>"):])
		if err != nil {
			return nil, err
		}
		m.Type = "github_content"
		return m, nil
	default:
		return nil, errors.New("invalid source")
	}
}

func parseGitHub(source string) (*Module, error) {
	// suzuki-shunsuke/lintnet-example-3#v0.1.0
	arr := strings.Split(source, "/")
	size := len(arr)
	if size < 2 { //nolint:gomnd
		return nil, errors.New("invalid source")
	}
	repoOwner := arr[0]
	lastElem := arr[size-1]
	a, ref, ok := strings.Cut(lastElem, "#")
	if !ok {
		return nil, errors.New("ref is required")
	}
	if len(arr) == 2 { //nolint:gomnd
		return &Module{
			RepoOwner: repoOwner,
			RepoName:  a,
			Ref:       ref,
		}, nil
	}
	return &Module{
		RepoOwner: repoOwner,
		RepoName:  arr[1],
		Path:      strings.Join(append(arr[2:size-1], a), "/"),
		Ref:       ref,
	}, nil
}
