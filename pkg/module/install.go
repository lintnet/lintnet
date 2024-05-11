package module

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/google/go-github/v62/github"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/osfile"
	"github.com/mholt/archiver/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

var errSubDirMustBeOne = errors.New("the number of sub directories must be one")

type ParamInstall struct {
	BaseDir string
}

type GitHub interface {
	GetArchiveLink(ctx context.Context, owner, repo string, archiveformat github.ArchiveFormat, opts *github.RepositoryContentGetOptions, maxRedirects int) (*url.URL, *github.Response, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Installer struct {
	fs   afero.Fs
	gh   GitHub
	http HTTPClient
}

func NewInstaller(fs afero.Fs, gh GitHub, httpClient HTTPClient) *Installer {
	return &Installer{
		fs:   fs,
		gh:   gh,
		http: httpClient,
	}
}

func (mi *Installer) Installs(ctx context.Context, logE *logrus.Entry, param *ParamInstall, modules map[string]*config.ModuleArchive) error {
	for _, mod := range modules {
		modID := mod.String()
		logE := logE.WithField("module_id", modID)
		if err := mi.Install(ctx, logE, param, mod); err != nil {
			return fmt.Errorf("install a module: %w", logerr.WithFields(err, logrus.Fields{
				"module_id": modID,
			}))
		}
	}
	return nil
}

func (mi *Installer) Install(ctx context.Context, logE *logrus.Entry, param *ParamInstall, mod *config.ModuleArchive) error { //nolint:funlen,cyclop
	// Check if the module is already downloaded
	dest := filepath.Join(param.BaseDir, filepath.FromSlash(mod.FilePath()))
	f, err := afero.DirExists(mi.fs, dest)
	if err != nil {
		return fmt.Errorf("check if the module is already installed: %w", err)
	}
	if f {
		return nil
	}
	if err := osfile.MkdirAll(mi.fs, filepath.Dir(dest)); err != nil {
		return fmt.Errorf("create parent directories: %w", err)
	}
	// Download Module
	u, _, err := mi.gh.GetArchiveLink(ctx, mod.RepoOwner, mod.RepoName, github.Tarball, &github.RepositoryContentGetOptions{
		Ref: mod.Ref,
	}, 5) //nolint:gomnd
	if err != nil {
		return fmt.Errorf("get an archive link by GitHub API: %w", logerr.WithFields(err, logrus.Fields{
			"moduel_repo_owner": mod.RepoOwner,
			"moduel_repo_name":  mod.RepoName,
			"moduel_ref":        mod.Ref,
		}))
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("create a HTTP request: %w", err)
	}
	resp, err := mi.http.Do(req)
	if err != nil {
		return fmt.Errorf("send a HTTP request: %w", err)
	}
	if resp.StatusCode >= 300 { //nolint:gomnd
		return errors.New("HTTP status code >= 300")
	}
	defer resp.Body.Close()
	tempDir, err := afero.TempDir(mi.fs, "", "")
	if err != nil {
		return fmt.Errorf("create a temporal directory: %w", err)
	}
	defer func() {
		if err := mi.fs.RemoveAll(tempDir); err != nil {
			logE.WithError(err).Warn("delete a temporal directory")
		}
	}()
	tempDest := filepath.Join(tempDir, "module.tar.gz")
	tempFile, err := mi.fs.Create(tempDest)
	if err != nil {
		return fmt.Errorf("create a temporal file: %w", err)
	}
	defer tempFile.Close()
	logE.Info("downloading a module")
	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		return fmt.Errorf("download a module on a temporal directory: %w", err)
	}
	tarGz := archiver.NewTarGz()
	unarchiveDest := filepath.Join(tempDir, "unarchived_dir")
	if err := tarGz.Unarchive(tempDest, unarchiveDest); err != nil {
		return fmt.Errorf("unarchive a tarball: %w", err)
	}
	dirs, err := os.ReadDir(unarchiveDest)
	if err != nil {
		return fmt.Errorf("read a directory: %w", err)
	}
	if len(dirs) != 1 {
		return errSubDirMustBeOne
	}
	if err := osfile.Copy(mi.fs, filepath.Join(unarchiveDest, dirs[0].Name()), dest); err != nil {
		return fmt.Errorf("copy a module from a teporal directory: %w", err)
	}
	return nil
}
