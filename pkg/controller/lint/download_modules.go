package lint

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/google/go-github/v57/github"
	"github.com/mholt/archiver/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const dirPermission os.FileMode = 0o775

type ParamDownloadModule struct {
	BaseDir string
}

type GitHub interface {
	GetArchiveLink(ctx context.Context, owner, repo string, archiveformat github.ArchiveFormat, opts *github.RepositoryContentGetOptions, maxRedirects int) (*url.URL, *github.Response, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c *Controller) downloadModules(ctx context.Context, logE *logrus.Entry, param *ParamDownloadModule, modMap map[string]*Module) error {
	for modID, mod := range modMap {
		if err := c.downloadModule(ctx, logE, param, modID, mod); err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) downloadModule(ctx context.Context, logE *logrus.Entry, param *ParamDownloadModule, modID string, mod *Module) error {
	// Check if the module is already downloaded
	dest := filepath.Join(param.BaseDir, filepath.FromSlash(modID))
	f, err := afero.DirExists(c.fs, dest)
	if err != nil {
		return err
	}
	if f {
		return nil
	}
	if err := c.fs.MkdirAll(filepath.Dir(dest), dirPermission); err != nil {
		return fmt.Errorf("create parent directories: %w", err)
	}
	// Download Module
	u, _, err := c.gh.GetArchiveLink(ctx, mod.RepoOwner, mod.RepoName, github.Tarball, nil, 5) //nolint:gomnd
	if err != nil {
		return fmt.Errorf("get an archive link by GitHub API: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("create a HTTP request: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("send a HTTP request: %w", err)
	}
	if resp.StatusCode >= 300 { //nolint:gomnd
		return errors.New("HTTP status code >= 300")
	}
	defer resp.Body.Close()
	tempDir, err := afero.TempDir(c.fs, "", "")
	if err != nil {
		return fmt.Errorf("create a temporal directory: %w", err)
	}
	defer c.fs.RemoveAll(tempDir)
	tempDest := filepath.Join(tempDir, "module.tar.gz")
	tempFile, err := c.fs.Create(tempDest)
	if err != nil {
		return fmt.Errorf("create a temporal file: %w", err)
	}
	defer tempFile.Close()
	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		return fmt.Errorf("download a module on a temporal directory: %w", err)
	}
	tarGz := archiver.NewTarGz()
	if err := tarGz.Unarchive(tempDest, dest); err != nil {
		return fmt.Errorf("unarchive a tarball: %w", err)
	}
	return nil
}
