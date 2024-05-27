package newcmd

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

//go:embed lint.jsonnet
var lintTemplate []byte

//go:embed test.jsonnet
var testTemplate []byte

func (c *Controller) New(_ context.Context, _ *logrus.Entry, fileName string) error {
	if !strings.HasSuffix(fileName, ".jsonnet") {
		return logerr.WithFields(errors.New("the file name must end with '.jsonnet'"), logrus.Fields{ //nolint:wrapcheck
			"file_name": fileName,
		})
	}

	if err := c.create(fileName, lintTemplate); err != nil {
		return err
	}

	testFileName := fileName[:len(fileName)-len(".jsonnet")] + "_test.jsonnet"
	if err := c.create(testFileName, testTemplate); err != nil {
		return err
	}

	return nil
}

func (c *Controller) create(fileName string, tpl []byte) error {
	if f, err := afero.Exists(c.fs, fileName); err != nil {
		return fmt.Errorf("check if the file exists: %w", logerr.WithFields(err, logrus.Fields{
			"file_name": fileName,
		}))
	} else if f {
		return nil
	}
	if err := afero.WriteFile(c.fs, fileName, tpl, 0o644); err != nil { //nolint:mnd
		return fmt.Errorf("create a file: %w", err)
	}
	return nil
}
