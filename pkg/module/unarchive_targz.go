package module

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lintnet/lintnet/pkg/osfile"
	"github.com/spf13/afero"
)

const maxFileSize = 1073741824 // 1GB

func extractTarGz(fs afero.Fs, src, dest string) error {
	f, err := fs.Open(src)
	if err != nil {
		return fmt.Errorf("open a tarball: %w", err)
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("read a gzip: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("get a new entry from a tar archive: %w", err)
		}

		targetPath := filepath.Join(dest, filepath.Clean(hdr.Name))

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := fs.MkdirAll(targetPath, os.FileMode(hdr.Mode)); err != nil { //nolint:gosec
				return fmt.Errorf("create a directory: %w", err)
			}
		case tar.TypeReg:
			if err := osfile.MkdirAll(fs, filepath.Dir(targetPath)); err != nil {
				return fmt.Errorf("create a directory: %w", err)
			}

			outFile, err := fs.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode)) //noilnt:gosec
			if err != nil {
				return fmt.Errorf("open a file: %w", err)
			}
			if _, err := io.CopyN(outFile, tr, maxFileSize); err != nil {
				outFile.Close()
				if errors.Is(err, io.EOF) {
					continue
				}
				return fmt.Errorf("copy a file: %w", err)
			}
			return errors.New("file size exceeds the limit")
		}
	}
}
