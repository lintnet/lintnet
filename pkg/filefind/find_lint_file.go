package filefind

import (
	"path/filepath"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/sirupsen/logrus"
)

func (f *FileFinder) FindLintFiles(logE *logrus.Entry, cfg *config.Config, cfgDir string) ([]*config.LintFile, error) {
	arr := make([]*config.LintFile, 0, len(cfg.Targets))
	for _, target := range cfg.Targets {
		lintFiles, err := f.findFilesFromLintFiles(logE, target.LintFiles, cfgDir, cfg.IgnoredPatterns)
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
