package filefind

import (
	"log/slog"
	"path/filepath"

	"github.com/lintnet/lintnet/pkg/config"
)

func (f *FileFinder) FindLintFiles(logger *slog.Logger, cfg *config.Config, cfgDir string) ([]*config.LintFile, error) {
	arr := make([]*config.LintFile, 0, len(cfg.Targets))
	for _, target := range cfg.Targets {
		lintFiles, err := f.findFilesFromLintFiles(logger, target.LintFiles, cfgDir, cfg.IgnoredPatterns)
		if err != nil {
			return nil, err
		}
		for _, lintFile := range lintFiles {
			if !filepath.IsAbs(lintFile.Path) {
				lintFile.Path = filepath.Join(cfgDir, lintFile.Path)
			}
		}
		arr = append(arr, lintFiles...)
	}
	return arr, nil
}
