package info

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

type ParamInfo struct {
	RootDir        string
	DataRootDir    string
	ConfigFilePath string
	PWD            string
	ModuleRootDir  bool
	MaskUser       bool
}

type Info struct {
	Version    string            `json:"vesrion,omitempty"`
	Commit     string            `json:"commit,omitempty"`
	Runtime    string            `json:"runtime"`
	ConfigFile string            `json:"config_file,omitempty"`
	RootDir    string            `json:"root_dir"`
	Env        map[string]string `json:"env"`
}

func (c *Controller) Info(_ context.Context, param *ParamInfo) error { //nolint:cyclop,funlen
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("get a current user: %w", err)
	}
	userName := currentUser.Username

	if param.ModuleRootDir {
		p := filepath.Join(param.RootDir, "modules")
		if param.MaskUser {
			p = maskUser(p, userName)
		}
		fmt.Fprintln(c.stdout, p)
		return nil
	}

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
			if param.MaskUser {
				v = maskUser(v, userName)
			}
			m[e] = v
		}
	}
	for _, e := range []string{"GITHUB_TOKEN", "LINTNET_GITHUB_TOKEN"} {
		if _, b := os.LookupEnv(e); b {
			m[e] = "(masked)"
		}
	}
	info := &Info{
		Version:    c.param.Version,
		Commit:     c.param.Commit,
		Runtime:    c.param.Env,
		ConfigFile: param.ConfigFilePath,
		RootDir:    param.RootDir,
		Env:        m,
	}
	if param.MaskUser {
		info.ConfigFile = maskUser(info.ConfigFile, userName)
		info.RootDir = maskUser(info.RootDir, userName)
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
