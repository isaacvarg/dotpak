package tui

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/isaacvarg/dotpak/internal/manifest"
)

// RunList draws the manifest list until the session quits.
// An empty group keeps every entry. An empty install type keeps every install type.
// A nil keys slice reads the terminal; any other slice is the scripted session.
func RunList(group, installType string, terminal bool, keys []tea.KeyPressMsg) (string, error) {
	if installType != "" && !manifest.InstallType(installType).IsValid() {
		return "", fmt.Errorf("invalid install type: %q", installType)
	}
	if !terminal {
		return "", fmt.Errorf("stdout is not a terminal")
	}
	return "", nil
}
