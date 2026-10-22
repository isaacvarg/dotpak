package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

const (
	colorSelected = "\x1b[38;2;238;212;159m" // #eed49f
	colorPacman   = "\x1b[38;2;240;198;198m" // #f0c6c6
	colorAur      = "\x1b[38;2;245;189;230m" // #f5bde6
)

const emptyListView = "Name                              Group  Install type  \uEA85\n\nj/k move • d remove • e edit • q quit"

const oneEntryManifest = `{
  "manifest": [
    {
      "name": "neovim",
      "group": "all",
      "install_type": "pacman",
      "install_command": "neovim"
    }
  ]
}
`

func TestRunListReturnsBeforeDrawing(t *testing.T) {
	tests := []struct {
		name        string
		installType string
		terminal    bool
		wantErr     string
	}{
		{
			name:     "non-terminal stdout",
			terminal: false,
			wantErr:  "stdout is not a terminal",
		},
		{
			name:        "invalid install type",
			installType: "snap",
			terminal:    true,
			wantErr:     `invalid install type: "snap"`,
		},
		{
			name:        "non-terminal stdout with a valid install type",
			installType: "pacman",
			terminal:    false,
			wantErr:     "stdout is not a terminal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("XDG_CONFIG_HOME", t.TempDir())
			view, err := RunList("", tt.installType, tt.terminal, []tea.KeyPressMsg{})
			if err == nil {
				t.Fatal("expected an error before drawing")
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("error %q, want %q", err.Error(), tt.wantErr)
			}
			if view != "" {
				t.Fatalf("view %q, want empty before drawing", view)
			}
			if _, statErr := os.Stat(filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "dotpak", "manifest.json")); !os.IsNotExist(statErr) {
				t.Fatalf("manifest file appeared, stat err %v", statErr)
			}
		})
	}
}

func TestRunListDrawsEmptyList(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
		group    string
	}{
		{name: "no manifest", manifest: ""},
		{name: "unknown group", manifest: oneEntryManifest, group: "missing"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := writeManifest(t, tt.manifest)
			view, err := RunList(tt.group, "", true, []tea.KeyPressMsg{{Text: "q"}})
			if err != nil {
				t.Fatalf("quit with an unchanged manifest, got %v", err)
			}
			if view != emptyListView {
				t.Fatalf("view\n%q\nwant\n%q", view, emptyListView)
			}
			assertManifestUnchanged(t, before)
		})
	}
}

const orderedManifest = `{
  "manifest": [
    {
      "name": "zsh",
      "group": "shell",
      "install_type": "aur",
      "install_command": "zsh"
    },
    {
      "name": "neovim",
      "group": "all",
      "install_type": "pacman",
      "install_command": "neovim"
    }
  ]
}
`

func TestRunListShowsEntriesInManifestOrder(t *testing.T) {
	before := writeManifest(t, orderedManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{{Text: "q"}})
	if err != nil {
		t.Fatalf("quit with an unchanged manifest, got %v", err)
	}
	assertManifestUnchanged(t, before)

	lines := strings.Split(view, "\n")
	if len(lines) < 4 {
		t.Fatalf("view %q", view)
	}
	header := lines[0]
	if strings.Contains(header, "\x1b") {
		t.Fatalf("header is styled %q", header)
	}
	if header != "Name                              Group  Install type  \uEA85" {
		t.Fatalf("header %q", header)
	}
	zsh := lines[1]
	neovim := lines[2]
	if !strings.Contains(zsh, colorSelected+"zsh") || !strings.Contains(zsh, colorSelected+"shell") {
		t.Fatalf("first entry is not selected %q", zsh)
	}
	if !strings.Contains(zsh, colorAur+"aur") || strings.Contains(zsh, colorSelected+"aur") {
		t.Fatalf("install type color %q", zsh)
	}
	if strings.Contains(neovim, colorSelected+"neovim") || strings.Contains(neovim, colorSelected+"all") {
		t.Fatalf("second entry is selected %q", neovim)
	}
	if !strings.Contains(neovim, colorPacman+"pacman") {
		t.Fatalf("pacman color %q", neovim)
	}
	plainZsh := stripANSI(zsh)
	plainNeovim := stripANSI(neovim)
	if strings.Contains(plainZsh, "\uEA85") || strings.Contains(plainNeovim, "\uEA85") {
		t.Fatalf("mark shown for an install command equal to the name\nzsh %q\nneovim %q", plainZsh, plainNeovim)
	}
	assertColumnOrder(t, plainZsh, "zsh", "shell", "aur")
	assertColumnOrder(t, plainNeovim, "neovim", "all", "pacman")
	if lines[len(lines)-1] != "j/k move • d remove • e edit • q quit" {
		t.Fatalf("hint %q", lines[len(lines)-1])
	}
}

func assertColumnOrder(t *testing.T, line string, name, group, installType string) {
	t.Helper()
	nameAt := strings.Index(line, name)
	groupAt := strings.Index(line, group)
	typeAt := strings.Index(line, installType)
	if nameAt < 0 || groupAt < 0 || typeAt < 0 || !(nameAt < groupAt && groupAt < typeAt) {
		t.Fatalf("columns Name, Group, Install type out of order on %q", line)
	}
}

func stripANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != 0x1b {
			b.WriteByte(s[i])
			continue
		}
		for i < len(s) && s[i] != 'm' {
			i++
		}
	}
	return b.String()
}

func writeManifest(t *testing.T, body string) []byte {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	if body == "" {
		return nil
	}
	cfg := filepath.Join(dir, "dotpak")
	if err := os.MkdirAll(cfg, 0o755); err != nil {
		t.Fatal(err)
	}
	payload := []byte(body)
	if err := os.WriteFile(filepath.Join(cfg, "manifest.json"), payload, 0o644); err != nil {
		t.Fatal(err)
	}
	return payload
}

func assertManifestUnchanged(t *testing.T, before []byte) {
	t.Helper()
	path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "dotpak", "manifest.json")
	after, err := os.ReadFile(path)
	if before == nil {
		if !os.IsNotExist(err) {
			t.Fatalf("manifest file appeared: %v", err)
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatalf("manifest changed\nbefore %s\nafter %s", before, after)
	}
}
