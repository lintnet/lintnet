package lint

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/lintnet/pkg/config"
	"github.com/suzuki-shunsuke/lintnet/pkg/module"
)

func (c *Controller) downloadModules(ctx context.Context, logE *logrus.Entry, cfg *config.Config) (map[string]string, error) {
	modules := make(map[string]string, len(cfg.Modules))
	for _, mod := range cfg.Modules {
		m, err := module.ParseSource(mod.Source)
		if err != nil {
			return nil, fmt.Errorf("parse a module source: %w", err)
		}
		p, err := c.getModulePath("", m)
		if err != nil {
			return nil, fmt.Errorf("get a module path: %w", err)
		}
		modules[mod.ID] = p
		f, err := afero.Exists(c.fs, p)
		if err != nil {
			return nil, fmt.Errorf("check if a module exists: %w", err)
		}
		if f {
			continue
		}
		if err := c.downloadModule(ctx, logE, m, p); err != nil {
			return nil, fmt.Errorf("download a module: %w", err)
		}
	}
	return modules, nil
}

func (c *Controller) downloadModule(ctx context.Context, logE *logrus.Entry, m *module.Module, modulePath string) error {
	return nil
}

func (c *Controller) getModulePath(rootDir string, mod *module.Module) (string, error) {
	// e.g. ~/.local/share/lintnet/modules/github_content/github.com/suzuki-shunsuke/lintnet-example/v1.0.0/foo.jsonnet
	switch mod.Type {
	case "github_content":
		return filepath.Join(rootDir, "modules", "github_content", "github.com", mod.RepoOwner, mod.RepoName, mod.Ref, mod.Path), nil
	case "github_archive":
		return filepath.Join(rootDir, "modules", "github_archive", "github.com", mod.RepoOwner, mod.RepoName, mod.Ref), nil
	}
	return "", errors.New("unknown module type")
}
