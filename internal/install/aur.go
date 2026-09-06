package install

import (
	"fmt"
	"os"
	"os/exec"
)

// these are the ones being checked
var aurHelpers = []string{"yay", "paru", "pikaur"}

func detectAURHelper() (string, error) {

	// returns first match so yay first, then paru, then blabla
	for _, h := range aurHelpers {
		if _, err := exec.LookPath(h); err == nil {
			return h, nil
		}
	}

	return "", fmt.Errorf("no aur helper found\nchecked: %v", aurHelpers)
}

func InstallAURPackages(names []string) error {
	aurHelper, err := detectAURHelper()
	if err != nil {
		return err
	}

	args := append([]string{ "-S", "--noconfirm"}, names...)

	command := exec.Command(aurHelper, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()

}
