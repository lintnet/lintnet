package info

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type ParamInfo struct {
	RootDir        string
	DataRootDir    string
	ConfigFilePath string
	PWD            string
}

type Info struct {
	Version     string            `json:"vesrion,omitempty"`
	ConfigFile  string            `json:"config_file,omitempty"`
	RootDir     string            `json:"root_dir"`
	DataRootDir string            `json:"data_root_dir"`
	Env         map[string]string `json:"env"`
}

func (c *Controller) Info(_ context.Context, logE *logrus.Entry, param *ParamInfo) error { //nolint:cyclop
	if param.ConfigFilePath == "" {
		for _, p := range []string{"lintnet.jsonnet", ".lintnet.jsonnet"} {
			f, err := afero.Exists(c.fs, p)
			if err != nil {
				return fmt.Errorf("check if a configuration file exists: %w", err)
			}
			if f {
				param.ConfigFilePath = p
			}
		}
	}

	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("get a current user: %w", err)
	}
	userName := currentUser.Username

	envs := []string{
		"LINTNET_CONFIG",
		"LINTNET_ERROR_LEVEL",
		"LINTNET_SHOWN_ERROR_LEVEL",
		"LINTNET_OUTPUT_SUCCESS",
		"LINTNET_LOG_LEVEL",
		"LINTNET_LOG_COLOR",
		"LINTNET_ROOT_DIR",
	}
	m := make(map[string]string, len(envs))
	for _, e := range envs {
		if v, b := os.LookupEnv(e); b {
			m[e] = maskUser(v, userName)
		}
	}
	for _, e := range []string{"GITHUB_TOKEN", "LINTNET_GITHUB_TOKEN"} {
		if _, b := os.LookupEnv(e); b {
			m[e] = "(masked)"
		}
	}
	info := &Info{
		Version:     c.param.Version,
		ConfigFile:  maskUser(param.ConfigFilePath, userName),
		RootDir:     maskUser(param.RootDir, userName),
		DataRootDir: maskUser(param.DataRootDir, userName),
		Env:         m,
	}
	encoder := json.NewEncoder(c.stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(info); err != nil {
		return fmt.Errorf("encode info as JSON to stdout: %w", err)
	}
	return nil
}

func maskUser(s, username string) string {
	return strings.ReplaceAll(s, username, "(USER)")
}
