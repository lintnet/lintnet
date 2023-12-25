package module

import (
	"errors"
	"regexp"
	"strings"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type Module struct {
	Type      string
	Host      string
	RepoOwner string
	RepoName  string
	Path      string
	Ref       string
	Tag       string
}

var fullCommitHashPattern = regexp.MustCompile("[a-fA-F0-9]{40}")

func validateRef(ref string) error {
	if fullCommitHashPattern.MatchString(ref) {
		return nil
	}
	return errors.New("ref must be full commit hash")
}

func (m *Module) ID() string {
	return strings.Join([]string{m.Host, m.RepoOwner, m.RepoName, m.Ref}, "/")
}

func ParseModuleLine(line string) (*Module, error) {
	// github.com/<repo owner>/<repo name>/<path>@<commit hash>[:<tag>]
	elems := strings.Split(line, "/")
	if len(elems) < 4 { //nolint:gomnd
		return nil, errors.New("line is invalid")
	}
	if elems[0] != "github.com" {
		return nil, errors.New("module must start with 'github.com/'")
	}
	size := len(elems)
	baseName, refAndTag, ok := strings.Cut(elems[size-1], "@")
	if !ok {
		return nil, errors.New("ref is required")
	}
	ref, tag, _ := strings.Cut(refAndTag, ":")
	if err := validateRef(ref); err != nil {
		return nil, err
	}
	return &Module{
		Type:      "github",
		Host:      "github.com",
		RepoOwner: elems[1],
		RepoName:  elems[2],
		Path:      strings.Join(append(elems[3:size-1], baseName), "/"),
		Ref:       ref,
		Tag:       tag,
	}, nil
}

func ListModules(cfg *config.Config) ([][]*Module, map[string]*Module, error) {
	modulesList := make([][]*Module, len(cfg.Targets))
	modules := map[string]*Module{}
	for i, target := range cfg.Targets {
		arr := make([]*Module, 0, len(target.Modules))
		for _, line := range target.Modules {
			mod, err := ParseModuleLine(line)
			if err != nil {
				return nil, nil, logerr.WithFields(err, logrus.Fields{ //nolint:wrapcheck
					"module": line,
				})
			}
			arr = append(arr, mod)
			modules[mod.ID()] = mod
		}
		modulesList[i] = arr
	}
	return modulesList, modules, nil
}
