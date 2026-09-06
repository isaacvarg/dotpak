package install

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func InstallOmarchyPlugins(names []string) error {
	if err := DeterctOmarchyPluginReady(); err != nil {
		return err
	}

	 var errs []error
	 for _, name := range names {
	         command := exec.Command("omarchy", "plugin", "install", name, "--enable")
	         command.Stdin = os.Stdin
	         command.Stdout = os.Stdout
	         command.Stderr = os.Stderr
	         if err := command.Run(); err != nil {
	                 errs = append(errs, fmt.Errorf("%s: %w", name, err))
	         }
	 }

	return errors.Join(errs...)
}

func DeterctOmarchyPluginReady() error {
	cmd := exec.Command("bash", "-c", "omarchy version")

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("omarchy does not appeared to be installed")
	}

	runes := []rune(string(out))
	major := int(runes[0])
	if major < 4 {
		return fmt.Errorf("version 4.0.0 or greater of omarchy is required")
	}

	// 	if _, err := exec.LookPath("omarchy-plugin-add"); err != nil {
	// 		return fmt.Errorf("omarchy plugin system not installed")
	// 	}
	//
	//

	return nil
}
