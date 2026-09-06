package install

import (
	"fmt"
	"os"
	"os/exec"
)

func IsFlatpakInstalled() bool {
	_, err := exec.LookPath("flatpak")
	return  err == nil
}

func InstallFlatpakPackages(appID []string) error {

	isReady := IsFlatpakInstalled()
	if !isReady {
		return fmt.Errorf("flatpak not installed")
	}
	
// TODO make this work with remote and not only flathub
	args := append([]string{"install", "flathub", "-y"}, appID...)

	command := exec.Command("flatpak", args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()

}
