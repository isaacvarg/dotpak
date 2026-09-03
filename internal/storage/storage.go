// Package storage handles directory parsing
package storage

import (
	"os"
	"path"
)

func ConfigDir() string {
	var configHome string
	configHome = os.Getenv("XDG_CONFIG_HOME")

	if configHome == "" {
		configHome = path.Join(os.Getenv("HOME"), ".config")
	}

	return path.Join(configHome, "dotpak")
}
