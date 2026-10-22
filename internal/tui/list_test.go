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

const filteredManifest = `{
  "manifest": [
    {
      "name": "neovim",
      "group": "all",
      "install_type": "pacman",
      "install_command": "neovim"
    },
    {
      "name": "spotify",
      "group": "media",
      "install_type": "flatpak",
      "install_command": "com.spotify.Client"
    },
    {
      "name": "node",
      "group": "dev",
      "install_type": "mise",
      "install_command": "node@20"
    },
    {
      "name": "ripgrep",
      "group": "dev",
      "install_type": "pacman",
      "install_command": "ripgrep"
    }
  ]
}
`

func TestRunListFiltersByGroupAndInstallType(t *testing.T) {
	tests := []struct {
		name        string
		group       string
		installType string
		wantNames   []string
		absentNames []string
	}{
		{
			name:        "group",
			group:       "all",
			wantNames:   []string{"neovim"},
			absentNames: []string{"spotify", "node", "ripgrep"},
		},
		{
			name:        "install type",
			installType: "flatpak",
			wantNames:   []string{"spotify"},
			absentNames: []string{"neovim", "node", "ripgrep"},
		},
		{
			name:        "group and install type",
			group:       "dev",
			installType: "mise",
			wantNames:   []string{"node"},
			absentNames: []string{"neovim", "spotify", "ripgrep"},
		},
		{
			name:        "group keeps manifest order",
			group:       "dev",
			wantNames:   []string{"node", "ripgrep"},
			absentNames: []string{"neovim", "spotify"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := writeManifest(t, filteredManifest)
			view, err := RunList(tt.group, tt.installType, true, []tea.KeyPressMsg{{Text: "q"}})
			if err != nil {
				t.Fatal(err)
			}
			assertManifestUnchanged(t, before)
			plain := stripANSI(view)
			prev := -1
			for _, name := range tt.wantNames {
				at := strings.Index(plain, name)
				if at < 0 || at <= prev {
					t.Fatalf("names %v not in manifest order in %q", tt.wantNames, plain)
				}
				prev = at
			}
			for _, name := range tt.absentNames {
				if strings.Contains(plain, name) {
					t.Fatalf("entry %s is in the list %q", name, plain)
				}
			}
		})
	}
}

const markedManifest = `{
  "manifest": [
    {
      "name": "neovim",
      "group": "all",
      "install_type": "pacman",
      "install_command": "neovim"
    },
    {
      "name": "zsh",
      "group": "all",
      "install_type": "pacman",
      "install_command": ""
    },
    {
      "name": "Node",
      "group": "dev",
      "install_type": "mise",
      "install_command": "node"
    },
    {
      "name": "node",
      "group": "dev",
      "install_type": "mise",
      "install_command": "node@20"
    },
    {
      "name": "spotify",
      "group": "media",
      "install_type": "flatpak",
      "install_command": "com.spotify.Client"
    },
    {
      "name": "webapp",
      "group": "apps",
      "install_type": "omarchy",
      "install_command": "webapp-install"
    }
  ]
}
`

func TestRunListMarksADifferentInstallCommand(t *testing.T) {
	const (
		colorMise    = "\x1b[38;2;238;153;160m"
		colorFlatpak = "\x1b[38;2;166;218;149m"
		colorOmarchy = "\x1b[38;2;245;169;127m"
	)
	before := writeManifest(t, markedManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{{Text: "q"}})
	if err != nil {
		t.Fatal(err)
	}
	assertManifestUnchanged(t, before)
	if strings.Contains(view, "node@20") || strings.Contains(view, "com.spotify.Client") || strings.Contains(view, "webapp-install") {
		t.Fatalf("install command text is on the list %q", view)
	}
	if strings.Contains(stripANSI(entryLine(t, view, "Node")), "node") {
		t.Fatal("install command node is on the Node line")
	}

	plainNeovim := stripANSI(entryLine(t, view, "neovim"))
	plainZsh := stripANSI(entryLine(t, view, "zsh"))
	if strings.Contains(plainNeovim, "\uEA85") || strings.Contains(plainZsh, "\uEA85") {
		t.Fatalf("mark on an equal or blank install command\nneovim %q\nzsh %q", plainNeovim, plainZsh)
	}
	for _, name := range []string{"Node", "node", "spotify", "webapp"} {
		line := stripANSI(entryLine(t, view, name))
		if !strings.Contains(line, "\uEA85") {
			t.Fatalf("entry %s has no mark on %q", name, line)
		}
	}
	if !strings.Contains(entryLine(t, view, "node"), colorMise+"mise") {
		t.Fatalf("mise color %q", entryLine(t, view, "node"))
	}
	if !strings.Contains(entryLine(t, view, "spotify"), colorFlatpak+"flatpak") {
		t.Fatalf("flatpak color %q", entryLine(t, view, "spotify"))
	}
	if !strings.Contains(entryLine(t, view, "webapp"), colorOmarchy+"omarchy") {
		t.Fatalf("omarchy color %q", entryLine(t, view, "webapp"))
	}
}

func TestRunListMovesWrapsAndQuits(t *testing.T) {
	names := []string{"neovim", "spotify", "node", "ripgrep"}
	quit := tea.KeyPressMsg{Text: "q"}
	tests := []struct {
		name         string
		keys         []tea.KeyPressMsg
		wantSelected string
		unchanged    bool
	}{
		{name: "j moves down", keys: []tea.KeyPressMsg{{Text: "j"}, quit}, wantSelected: "spotify"},
		{name: "down moves down", keys: []tea.KeyPressMsg{{Code: tea.KeyDown}, quit}, wantSelected: "spotify"},
		{name: "k wraps to the last entry", keys: []tea.KeyPressMsg{{Text: "k"}, quit}, wantSelected: "ripgrep"},
		{name: "up wraps to the last entry", keys: []tea.KeyPressMsg{{Code: tea.KeyUp}, quit}, wantSelected: "ripgrep"},
		{
			name:         "j wraps from the last entry",
			keys:         []tea.KeyPressMsg{{Text: "k"}, {Text: "j"}, quit},
			wantSelected: "neovim",
		},
		{
			name:         "three j presses reach the last entry",
			keys:         []tea.KeyPressMsg{{Text: "j"}, {Text: "j"}, {Text: "j"}, quit},
			wantSelected: "ripgrep",
		},
		{name: "enter leaves the selection", keys: []tea.KeyPressMsg{{Code: tea.KeyEnter}, quit}, wantSelected: "neovim", unchanged: true},
		{name: "slash does nothing", keys: []tea.KeyPressMsg{{Text: "/"}, quit}, wantSelected: "neovim", unchanged: true},
		{name: "d does not change the manifest", keys: []tea.KeyPressMsg{{Text: "d"}, quit}, wantSelected: "neovim", unchanged: true},
		{name: "e then esc restores the list", keys: []tea.KeyPressMsg{{Text: "e"}, {Code: tea.KeyEscape}}, wantSelected: "neovim", unchanged: true},
		{name: "q quits", keys: []tea.KeyPressMsg{quit}, wantSelected: "neovim", unchanged: true},
		{name: "esc quits", keys: []tea.KeyPressMsg{{Code: tea.KeyEscape}}, wantSelected: "neovim", unchanged: true},
		{name: "ctrl+c quits", keys: []tea.KeyPressMsg{{Code: 'c', Mod: tea.ModCtrl}}, wantSelected: "neovim", unchanged: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := writeManifest(t, filteredManifest)
			baseline, err := RunList("", "", true, []tea.KeyPressMsg{quit})
			if err != nil {
				t.Fatal(err)
			}
			view, err := RunList("", "", true, tt.keys)
			if err != nil {
				t.Fatalf("quit with an unchanged manifest, got %v", err)
			}
			assertManifestUnchanged(t, before)
			if got := selectedEntry(t, view, names); got != tt.wantSelected {
				t.Fatalf("selected %s, want %s", got, tt.wantSelected)
			}
			if tt.unchanged && view != baseline {
				t.Fatalf("view changed\n%s\nwant\n%s", view, baseline)
			}
			if strings.Contains(view, "remove neovim?") {
				t.Fatalf("remove prompt shown %q", view)
			}
		})
	}
}

func TestRunListEmptyListIgnoresMovement(t *testing.T) {
	writeManifest(t, "")
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "j"},
		{Text: "k"},
		{Code: tea.KeyDown},
		{Code: tea.KeyUp},
		{Text: "d"},
		{Text: "e"},
		{Text: "/"},
		{Code: tea.KeyEnter},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if view != emptyListView {
		t.Fatalf("view\n%q\nwant\n%q", view, emptyListView)
	}
}

const longNameManifest = `{
  "manifest": [
    {
      "name": "abcdefghijklmnopqrstuvwxyz0123456789",
      "group": "widgets",
      "install_type": "omarchy",
      "install_command": "custom-identifier"
    },
    {
      "name": "abcdefghijklmnopqrstuvwxyz012345",
      "group": "widgets",
      "install_type": "pacman",
      "install_command": "abcdefghijklmnopqrstuvwxyz012345"
    }
  ]
}
`

func TestRunListCutsALongName(t *testing.T) {
	const (
		fullName  = "abcdefghijklmnopqrstuvwxyz0123456789"
		cutName   = "abcdefghijklmnopqrstuvwxyz01234…"
		exactName = "abcdefghijklmnopqrstuvwxyz012345"
	)
	before := writeManifest(t, longNameManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{{Text: "q"}})
	if err != nil {
		t.Fatal(err)
	}
	assertManifestUnchanged(t, before)
	if strings.Contains(view, fullName) {
		t.Fatal("long name is not cut")
	}
	if strings.Contains(view, "custom-identifier") {
		t.Fatal("install command text is on the list")
	}
	longLine := stripANSI(lineContaining(t, view, cutName))
	if strings.Contains(longLine, "\n") {
		t.Fatalf("long entry wraps %q", longLine)
	}
	for _, part := range []string{cutName, "widgets", "omarchy", "\uEA85"} {
		if !strings.Contains(longLine, part) {
			t.Fatalf("long line %q missing %s", longLine, part)
		}
	}
	exactLine := stripANSI(lineContaining(t, view, exactName))
	if strings.Contains(exactLine, "…") {
		t.Fatalf("name that fits was cut %q", exactLine)
	}
	for _, part := range []string{exactName, "widgets", "pacman"} {
		if !strings.Contains(exactLine, part) {
			t.Fatalf("exact line %q missing %s", exactLine, part)
		}
	}
	if strings.Contains(exactLine, "\uEA85") {
		t.Fatalf("mark shown for an install command equal to the name %q", exactLine)
	}
}

func lineContaining(t *testing.T, view, needle string) string {
	t.Helper()
	for _, line := range strings.Split(view, "\n") {
		if strings.Contains(stripANSI(line), needle) {
			return line
		}
	}
	t.Fatalf("no line contains %q in %q", needle, view)
	return ""
}

func selectedEntry(t *testing.T, view string, names []string) string {
	t.Helper()
	var got string
	for _, name := range names {
		if strings.Contains(entryLine(t, view, name), colorSelected+name) {
			if got != "" {
				t.Fatalf("both %s and %s are selected", got, name)
			}
			got = name
		}
	}
	if got == "" {
		t.Fatal("no entry is selected")
	}
	return got
}

func entryLine(t *testing.T, view, name string) string {
	t.Helper()
	for _, line := range strings.Split(view, "\n") {
		plain := stripANSI(line)
		if strings.HasPrefix(plain, name) && (len(plain) == len(name) || plain[len(name)] == ' ') {
			return line
		}
	}
	t.Fatalf("no line for entry %s in %q", name, view)
	return ""
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
