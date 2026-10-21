package install

import (
	"fmt"
	"os"
	"os/exec"
)

func installMiseTools(specs []string) error {
	if _, err := exec.LookPath("mise"); err != nil {
		return fmt.Errorf("mise is not installed")
	}

	args := append([]string{"use", "-g"}, specs...)
	command := exec.Command("mise", args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("mise failed: %w", err)
	}
	return nil
}
