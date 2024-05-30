package filefind

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/domain"
	"github.com/spf13/afero"
)

func (f *FileFinder) findDataFiles(dataBasePath string, files []*config.DataFile, cfgDir string, ignorePatterns []string) ([][]*domain.Path, error) { //nolint:cyclop
	if dataBasePath == "" {
		dataFiles, err := f.findFilesFromPaths(files, cfgDir, ignorePatterns)
		if err != nil {
			return nil, err
		}
		return [][]*domain.Path{dataFiles}, nil
	}

	matches := map[string]struct{}{}
	if err := doublestar.GlobWalk(afero.NewIOFS(f.fs), filepath.Join(cfgDir, filepath.FromSlash(dataBasePath)), func(path string, d fs.DirEntry) error {
		if err := ignorePath(path, ignorePatterns); err != nil {
			return err
		}
		if !d.IsDir() {
			path = filepath.Dir(path)
		}
		matches[path] = struct{}{}
		return nil
	}, doublestar.WithNoFollow()); err != nil {
		return nil, fmt.Errorf("search files: %w", err)
	}
	paths := make([][]*domain.Path, 0, len(matches))
	for rootPath := range matches {
		dataFiles, err := f.findFilesFromPaths(files, rootPath, ignorePatterns)
		if err != nil {
			return nil, err
		}
		base, err := filepath.Rel(cfgDir, rootPath)
		if err != nil {
			return nil, fmt.Errorf("get a relative path from configuration file: %w", err)
		}
		for _, dataFile := range dataFiles {
			if !filepath.IsAbs(dataFile.Raw) {
				dataFile.Raw = filepath.Join(base, dataFile.Raw)
			}
		}
		paths = append(paths, dataFiles)
	}
	return paths, nil
}
