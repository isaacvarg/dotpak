package install

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var devCheckoutVersion = regexp.MustCompile(`^dev \([0-9a-fA-F]+\)$`)

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

func omarchyVersionReady(line string) error {
	version := strings.TrimSpace(line)
	if devCheckoutVersion.MatchString(version) {
		return nil
	}
	pkgver := version
	if epoch, rest, ok := strings.Cut(pkgver, ":"); ok && digitsOnly(epoch) {
		pkgver = rest
	}
	if i := strings.LastIndexByte(pkgver, '-'); i >= 0 {
		pkgver = pkgver[:i]
	}
	if pkgver == "dev" || strings.HasPrefix(pkgver, "dev.") {
		return nil
	}
	if versionAtLeast(leadingVersionNumbers(pkgver), []int{4, 0, 0}) {
		return nil
	}
	return fmt.Errorf("version 4.0.0 or greater of omarchy is required: %q", line)
}

func leadingVersionNumbers(pkgver string) []int {
	var nums []int
	for _, part := range strings.Split(pkgver, ".") {
		if !digitsOnly(part) {
			break
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			break
		}
		nums = append(nums, n)
	}
	return nums
}

func digitsOnly(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func versionAtLeast(got, min []int) bool {
	if len(got) == 0 {
		return false
	}
	n := max(len(got), len(min))
	for i := range n {
		g, m := 0, 0
		if i < len(got) {
			g = got[i]
		}
		if i < len(min) {
			m = min[i]
		}
		if g != m {
			return g > m
		}
	}
	return true
}

func DeterctOmarchyPluginReady() error {
	out, err := exec.Command("omarchy", "version").Output()
	if err != nil {
		return fmt.Errorf("omarchy does not appear to be installed")
	}
	return omarchyVersionReady(string(out))
}
