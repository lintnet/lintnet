package lint

import (
	"errors"
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
}

func (m *Module) ID() string {
	return strings.Join([]string{m.Host, m.RepoOwner, m.RepoName, m.Ref}, "/")
}

func parseModuleLine(line string) (*Module, error) {
	// github.com/<repo owner>/<repo name>/<path>@<ref>
	elems := strings.Split(line, "/")
	if len(elems) < 4 { //nolint:gomnd
		return nil, errors.New("line is invalid")
	}
	if elems[0] != "github.com" {
		return nil, errors.New("module must start with 'github.com/'")
	}
	size := len(elems)
	baseName, ref, ok := strings.Cut(elems[size-1], "@")
	if !ok {
		return nil, errors.New("ref is required")
	}
	return &Module{
		Type:      "github",
		Host:      "github.com",
		RepoOwner: elems[1],
		RepoName:  elems[2],
		Path:      strings.Join(append(elems[3:size-1], baseName), "/"),
		Ref:       ref,
	}, nil
}

func (c *Controller) listModules(logE *logrus.Entry, cfg *config.Config) ([][]*Module, map[string]*Module, error) {
	modulesList := make([][]*Module, len(cfg.Targets))
	modules := map[string]*Module{}
	for i, target := range cfg.Targets {
		arr := make([]*Module, 0, len(target.Modules))
		for _, line := range target.Modules {
			mod, err := parseModuleLine(line)
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
