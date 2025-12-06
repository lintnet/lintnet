package newcmd

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

//go:embed lint.jsonnet
var lintTemplate []byte

//go:embed test.jsonnet
var testTemplate []byte

func (c *Controller) New(_ context.Context, _ *slog.Logger, fileName string) error {
	if !strings.HasSuffix(fileName, ".jsonnet") {
		return slogerr.With(errors.New("the file name must end with '.jsonnet'"), "file_name", fileName) //nolint:wrapcheck
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
		return fmt.Errorf("check if the file exists: %w", slogerr.With(err, "file_name", fileName))
	} else if f {
		return nil
	}
	if err := afero.WriteFile(c.fs, fileName, tpl, 0o644); err != nil { //nolint:mnd
		return fmt.Errorf("create a file: %w", err)
	}
	return nil
}
