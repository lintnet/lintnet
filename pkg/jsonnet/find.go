package jsonnet

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func FindFiles(logE *logrus.Entry, afs afero.Fs) ([]string, error) {
	filePaths := []string{}
	if err := fs.WalkDir(afero.NewIOFS(afs), "", func(p string, dirEntry fs.DirEntry, e error) error {
		if e != nil {
			logE.WithError(e).Warn("error occurred while searching files")
			return nil
		}
		if dirEntry.Type().IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".jsonnet") {
			return nil
		}
		filePaths = append(filePaths, p)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("walks the file tree of the unarchived package: %w", err)
	}
	return filePaths, nil
}
