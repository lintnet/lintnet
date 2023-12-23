package lint

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/go-github/v57/github"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/sirupsen/logrus"
)

type ParamDownloadModule struct {
	BaseDir string
}

type GitHub interface {
	GetArchiveLink(ctx context.Context, owner, repo string, archiveformat github.ArchiveFormat, opts *github.RepositoryContentGetOptions, maxRedirects int) (*url.URL, *github.Response, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c *Controller) installModules(ctx context.Context, logE *logrus.Entry, param *module.ParamInstall, modMap map[string]*module.Module) error {
	for modID, mod := range modMap {
		logE := logE.WithField("module_id", modID)
		if err := c.moduleInstaller.Install(ctx, logE, param, modID, mod); err != nil {
			return err
		}
	}
	return nil
}
