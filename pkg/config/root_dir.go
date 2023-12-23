package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

func GetRootDir() (string, error) {
	// ${LINTNET_ROOT_DIR:-${XDG_DATA_HOME:-$HOME/.local/share}/lintnet}
	xdgDataHome := xdg.DataHome
	if xdgDataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get the current user home directory: %w", err)
		}
		xdgDataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(xdgDataHome, "lintnet"), nil
}
