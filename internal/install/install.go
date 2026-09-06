// Package install
package install

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/isaacvarg/dotpak/internal/manifest"
)

func Install(group string) error {
	data, err := manifest.Load()
	if err != nil {
		return err
	}
	var groupData []manifest.ManifestEntry
	for _, e := range data.Manifest {
		if e.Group == "all" || e.Group == group {
			groupData = append(groupData, e)
		}
	}

	var pacman []string
	for _, e	 := range groupData {
		if e.InstallType == "pacman" {
			pacman = append(pacman, e.InstallCommand)
		}
	}

	err = installPacmanPackages(pacman)
	if err != nil {
	return fmt.Errorf("pacman failed: %w", err )
	}

	return nil
}

func installPacmanPackages(names []string) error {
	args := append([]string{"pacman", "-S", "--noconfirm"}, names...)

	command := exec.Command("sudo", args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()

}
