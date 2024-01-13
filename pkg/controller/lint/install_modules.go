package lint

import (
	"context"
	"fmt"

	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type ParamDownloadModule struct {
	BaseDir string
}

func (c *Controller) installModules(ctx context.Context, logE *logrus.Entry, param *module.ParamInstall, modules map[string]*config.ModuleArchive) error {
	for _, mod := range modules {
		logE := logE.WithField("module_id", mod.ID)
		if err := c.moduleInstaller.Install(ctx, logE, param, mod); err != nil {
			return fmt.Errorf("install a module: %w", logerr.WithFields(err, logrus.Fields{
				"module_id": mod.ID,
			}))
		}
	}
	return nil
}
