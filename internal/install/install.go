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
	var aur []string
	var flatpak []string
	var omarchy []string

	for _, e := range groupData {
		switch e.InstallType {

		case "pacman":
			pacman = append(pacman, e.InstallCommand)

		case "aur":
			aur = append(aur, e.InstallCommand)
		case "flatpak":
			flatpak = append(flatpak, e.InstallCommand)
		case "omarchy":
			omarchy = append(omarchy, e.InstallCommand)
		}
	}

	if len(pacman) != 0 {
		err = installPacmanPackages(pacman)
		if err != nil {
			return fmt.Errorf("pacman failed: %w", err)
		}
	}

	if len(aur) != 0 {
		err = InstallAURPackages(aur)
		if err != nil {
			return fmt.Errorf("aur installer failed: %w", err)
		}
	}

	if len(flatpak) != 0 {
		err = InstallFlatpakPackages(flatpak)
		if err != nil {
			return fmt.Errorf("flatpak installer failed: %w", err)
		}
	}

	if len(omarchy) != 0 {
		err = InstallOmarchyPlugins(omarchy)
		if err != nil {
			return fmt.Errorf("omarchy plugins installer failed: %w", err)
		}
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
