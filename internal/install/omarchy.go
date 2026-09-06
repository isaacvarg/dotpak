package install

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func InstallOmarchyPlugins(names []string) error {
      if err := DeterctOmarchyPluginReady(); err != nil {
              return err
      }

      var errs []error
      for _, name := range names {
              command := exec.Command("omarchy", "plugin", "add", name, "--enable")
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

	cmd := exec.Command("bash", "-c", "omarchy", "version")

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("omarchy not installed: %w", err)
	}

	v, err := parseOmarchyVersion(string(out))
	if err != nil {
		return err
	}

	if !v.atLeast(version{4,0,0}) {
		return fmt.Errorf("omarchy plugins require omarchy v4.0.0+, while the version installed is %s", strings.TrimSpace(string(out)))
	}

// 	if _, err := exec.LookPath("omarchy-plugin-add"); err != nil {
// 		return fmt.Errorf("omarchy plugin system not installed")
// 	}
// 
	//
	
	return nil
}

// omarchy version parsing
type version struct {
      major, minor, patch int
}

func parseOmarchyVersion(raw string) (version, error) {
      v := strings.TrimSpace(raw)
      v, _, _ = strings.Cut(v, "-") // drop pkgrel/build suffix, e.g. "4.0.2-1" -> "4.0.2"

      parts := strings.Split(v, ".")
      if len(parts) < 3 {
              return version{}, fmt.Errorf("unexpected omarchy version format: %q", raw)
      }

      nums := make([]int, 3)
      for i := 0; i < 3; i++ {
              n, err := strconv.Atoi(parts[i])
              if err != nil {
                      return version{}, fmt.Errorf("invalid omarchy version %q: %w", raw, err)
              }
              nums[i] = n
      }
      return version{nums[0], nums[1], nums[2]}, nil
}

func (v version) atLeast(other version) bool {
      if v.major != other.major {
              return v.major > other.major
      }
      if v.minor != other.minor {
              return v.minor > other.minor
      }
      return v.patch >= other.patch
}
