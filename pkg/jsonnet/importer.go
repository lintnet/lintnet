package jsonnet

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-jsonnet"
	"github.com/lintnet/lintnet/pkg/config"
	"github.com/lintnet/lintnet/pkg/module"
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type FileImporter = jsonnet.FileImporter

type Importer struct {
	ctx             context.Context //nolint:containedctx
	logE            *logrus.Entry
	param           *module.ParamInstall
	importer        jsonnet.Importer
	moduleInstaller *module.Installer
}

func NewImporter(ctx context.Context, logE *logrus.Entry, param *module.ParamInstall, importer jsonnet.Importer, installer *module.Installer) *Importer {
	return &Importer{
		ctx:             ctx,
		logE:            logE,
		param:           param,
		importer:        importer,
		moduleInstaller: installer,
	}
}

func (ip *Importer) Import(importedFrom, importedPath string) (jsonnet.Contents, string, error) {
	contents, foundAt, err := ip.importer.Import(importedFrom, importedPath)
	if err == nil {
		return contents, foundAt, nil
	}
	if !strings.HasPrefix(importedPath, "github.com/") {
		return contents, foundAt, err //nolint:wrapcheck
	}
	mod, err := config.ParseModuleLine(importedPath)
	if err != nil {
		return contents, foundAt, fmt.Errorf("parse a module import path: %w", err)
	}
	if err := ip.moduleInstaller.Install(ip.ctx, ip.logE, ip.param, mod.Archive); err != nil {
		return contents, foundAt, fmt.Errorf("install a module: %w", logerr.WithFields(err, logrus.Fields{
			"module_id": mod.Archive.ID,
			"import":    importedPath,
		}))
	}
	return ip.importer.Import(importedFrom, mod.SlashPath) //nolint:wrapcheck
}
