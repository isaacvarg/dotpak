package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestRunListEditOpensOnTheSelectedName(t *testing.T) {
	tests := []struct {
		name     string
		keys     []tea.KeyPressMsg
		wantName string
		absent   string
	}{
		{
			name:     "first entry",
			keys:     []tea.KeyPressMsg{{Text: "e"}},
			wantName: "neovim",
			absent:   "spotify",
		},
		{
			name:     "moved entry",
			keys:     []tea.KeyPressMsg{{Text: "j"}, {Text: "e"}},
			wantName: "spotify",
			absent:   "neovim",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := writeManifest(t, filteredManifest)
			view, err := RunList("", "", true, tt.keys)
			if err != nil {
				t.Fatal(err)
			}
			assertManifestUnchanged(t, before)
			if !strings.Contains(view, "What is the name?") {
				t.Fatalf("name step missing\n%s", view)
			}
			if !strings.Contains(view, tt.wantName) {
				t.Fatalf("name step does not start on %s\n%s", tt.wantName, view)
			}
			if strings.Contains(stripANSI(view), tt.absent) {
				t.Fatalf("name step shows %s\n%s", tt.absent, view)
			}
			if strings.Contains(view, "What is the install type?") {
				t.Fatalf("edit skipped the name step\n%s", view)
			}
		})
	}
}

func TestRunListEditNameStepTakesLetters(t *testing.T) {
	before := writeManifest(t, oneEntryManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Text: "j"},
		{Text: "k"},
		{Text: "d"},
		{Text: "e"},
		{Code: 'x'},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertManifestUnchanged(t, before)
	plain := stripANSI(view)
	if !strings.Contains(plain, "neovimjkde") {
		t.Fatalf("name step did not take j, k, d, and e as characters\n%s", plain)
	}
	if strings.Contains(plain, "neovimjkdex") {
		t.Fatalf("a code-only key inserted a character\n%s", plain)
	}
	if strings.Contains(view, "What is the install type?") {
		t.Fatalf("letters left the name step\n%s", view)
	}
}

func TestRunListEditConfirmsTheName(t *testing.T) {
	before := writeManifest(t, oneEntryManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Code: tea.KeyEnter},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertManifestUnchanged(t, before)
	if !strings.Contains(view, "What is the install type?") {
		t.Fatalf("enter did not leave the name step\n%s", view)
	}
	if strings.Contains(view, "What is the name?") {
		t.Fatalf("name step still showing\n%s", view)
	}
}

func TestRunListEditEmptyNameStays(t *testing.T) {
	tests := []struct {
		name string
		keys []tea.KeyPressMsg
	}{
		{
			name: "cleared",
			keys: []tea.KeyPressMsg{{Text: "e"}, {Code: 'u', Mod: tea.ModCtrl}, {Code: tea.KeyEnter}},
		},
		{
			name: "spaces",
			keys: []tea.KeyPressMsg{{Text: "e"}, {Code: 'u', Mod: tea.ModCtrl}, {Text: "   "}, {Code: tea.KeyEnter}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := writeManifest(t, oneEntryManifest)
			view, err := RunList("", "", true, tt.keys)
			if err != nil {
				t.Fatal(err)
			}
			assertManifestUnchanged(t, before)
			if !strings.Contains(view, "What is the name?") {
				t.Fatalf("left the name step\n%s", view)
			}
			if strings.Contains(view, "What is the install type?") {
				t.Fatalf("empty name advanced\n%s", view)
			}
			if !strings.Contains(stripANSI(view), "a name is required") {
				t.Fatalf("name step does not say a name is required\n%s", view)
			}
		})
	}
}

func TestRunListEditDuplicateNameStays(t *testing.T) {
	tests := []struct {
		name  string
		group string
		typed string
		stay  bool
	}{
		{name: "another entry", typed: "spotify", stay: true},
		{name: "another entry ignoring case", typed: "SPOTIFY", stay: true},
		{name: "entry outside the filter", group: "all", typed: "node", stay: true},
		{name: "case-only change of this entry", typed: "Neovim", stay: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := writeManifest(t, filteredManifest)
			view, err := RunList(tt.group, "", true, []tea.KeyPressMsg{
				{Text: "e"},
				{Code: 'u', Mod: tea.ModCtrl},
				{Text: tt.typed},
				{Code: tea.KeyEnter},
			})
			if err != nil {
				t.Fatal(err)
			}
			assertManifestUnchanged(t, before)
			if tt.stay {
				if !strings.Contains(view, "What is the name?") || strings.Contains(view, "What is the install type?") {
					t.Fatalf("duplicate name left the name step\n%s", view)
				}
				if !strings.Contains(stripANSI(view), "entry already exists in the manifest") {
					t.Fatalf("duplicate name has no error\n%s", view)
				}
				return
			}
			if !strings.Contains(view, "What is the install type?") {
				t.Fatalf("case-only rename did not leave the name step\n%s", view)
			}
		})
	}
}

const editTypesManifest = `{
  "manifest": [
    {
      "name": "neovim",
      "group": "all",
      "install_type": "pacman",
      "install_command": "neovim"
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
    }
  ]
}
`

func TestRunListEditInstallTypeStartsOnTheCurrentValue(t *testing.T) {
	tests := []struct {
		name string
		keys []tea.KeyPressMsg
		want string
	}{
		{name: "pacman", keys: []tea.KeyPressMsg{{Text: "e"}, {Code: tea.KeyEnter}}, want: "pacman"},
		{name: "mise", keys: []tea.KeyPressMsg{{Text: "j"}, {Text: "e"}, {Code: tea.KeyEnter}}, want: "mise"},
		{name: "flatpak", keys: []tea.KeyPressMsg{{Text: "j"}, {Text: "j"}, {Text: "e"}, {Code: tea.KeyEnter}}, want: "flatpak"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := writeManifest(t, editTypesManifest)
			view, err := RunList("", "", true, tt.keys)
			if err != nil {
				t.Fatal(err)
			}
			assertManifestUnchanged(t, before)
			plain := stripANSI(view)
			if !strings.Contains(view, "install type") || strings.Contains(view, "installation type") {
				t.Fatalf("install type label\n%s", view)
			}
			if got := selectedChoice(t, plain); got != tt.want {
				t.Fatalf("selected %s, want %s\n%s", got, tt.want, plain)
			}
			if tt.want == "pacman" {
				assertChoiceOrder(t, plain, []string{"pacman", "aur", "mise", "omarchy", "flatpak"})
			}
		})
	}
}

func TestRunListEditInstallTypeMovesAndWraps(t *testing.T) {
	enter := tea.KeyPressMsg{Code: tea.KeyEnter}
	tests := []struct {
		name string
		keys []tea.KeyPressMsg
		want string
	}{
		{name: "j moves down", keys: []tea.KeyPressMsg{{Text: "e"}, enter, {Text: "j"}}, want: "aur"},
		{name: "down moves down", keys: []tea.KeyPressMsg{{Text: "e"}, enter, {Code: tea.KeyDown}}, want: "aur"},
		{name: "k wraps to flatpak", keys: []tea.KeyPressMsg{{Text: "e"}, enter, {Text: "k"}}, want: "flatpak"},
		{name: "up wraps to flatpak", keys: []tea.KeyPressMsg{{Text: "e"}, enter, {Code: tea.KeyUp}}, want: "flatpak"},
		{name: "j wraps from flatpak", keys: []tea.KeyPressMsg{{Text: "j"}, {Text: "j"}, {Text: "e"}, enter, {Text: "j"}}, want: "pacman"},
		{name: "down wraps from flatpak", keys: []tea.KeyPressMsg{{Text: "j"}, {Text: "j"}, {Text: "e"}, enter, {Code: tea.KeyDown}}, want: "pacman"},
		{name: "j moves from mise to omarchy", keys: []tea.KeyPressMsg{{Text: "j"}, {Text: "e"}, enter, {Text: "j"}}, want: "omarchy"},
		{name: "k moves from mise to aur", keys: []tea.KeyPressMsg{{Text: "j"}, {Text: "e"}, enter, {Text: "k"}}, want: "aur"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := writeManifest(t, editTypesManifest)
			view, err := RunList("", "", true, tt.keys)
			if err != nil {
				t.Fatal(err)
			}
			assertManifestUnchanged(t, before)
			plain := stripANSI(view)
			if !strings.Contains(plain, "install type") {
				t.Fatalf("left the install type step\n%s", plain)
			}
			if got := selectedChoice(t, plain); got != tt.want {
				t.Fatalf("selected %s, want %s\n%s", got, tt.want, plain)
			}
		})
	}
}

const storedGroups = `{
  "groups": [
    {"name": "shell"},
    {"name": "dev"}
  ]
}
`

const storedGroupsWithAll = `{
  "groups": [
    {"name": "all"},
    {"name": "shell"}
  ]
}
`

func TestRunListEditGroupOffersExistingGroups(t *testing.T) {
	tests := []struct {
		name    string
		entry   string
		groups  string
		want    []string
		current string
	}{
		{
			name:    "only all",
			entry:   `{"name":"neovim","group":"all","install_type":"pacman","install_command":"neovim"}`,
			want:    []string{"all"},
			current: "all",
		},
		{
			name:    "stored groups",
			entry:   `{"name":"zsh","group":"shell","install_type":"aur","install_command":"zsh"}`,
			groups:  storedGroups,
			want:    []string{"all", "shell", "dev"},
			current: "shell",
		},
		{
			name:    "current group missing from the list",
			entry:   `{"name":"widget","group":"orphans","install_type":"omarchy","install_command":"widget"}`,
			groups:  storedGroups,
			want:    []string{"all", "shell", "dev", "orphans"},
			current: "orphans",
		},
		{
			name:    "stored all is not duplicated",
			entry:   `{"name":"neovim","group":"all","install_type":"pacman","install_command":"neovim"}`,
			groups:  storedGroupsWithAll,
			want:    []string{"all", "shell"},
			current: "all",
		},
		{
			name:    "current group case is kept",
			entry:   `{"name":"neovim","group":"ALL","install_type":"pacman","install_command":"neovim"}`,
			groups:  storedGroupsWithAll,
			want:    []string{"all", "shell", "ALL"},
			current: "ALL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeManifest(t, `{"manifest":[`+tt.entry+`]}`)
			var before []byte
			if tt.groups != "" {
				before = writeGroups(t, tt.groups)
			}
			view, err := RunList("", "", true, []tea.KeyPressMsg{
				{Text: "e"},
				{Code: tea.KeyEnter},
				{Code: tea.KeyEnter},
			})
			if err != nil {
				t.Fatal(err)
			}
			plain := stripANSI(view)
			if !strings.Contains(plain, "Which group?") {
				t.Fatalf("group step missing\n%s", plain)
			}
			if strings.Contains(plain, "What is the install type?") {
				t.Fatalf("still on install type\n%s", plain)
			}
			if got := choiceLines(plain); strings.Join(got, ",") != strings.Join(tt.want, ",") {
				t.Fatalf("groups %v, want %v\n%s", got, tt.want, plain)
			}
			if got := selectedChoice(t, plain); got != tt.current {
				t.Fatalf("selected %s, want %s", got, tt.current)
			}
			if tt.groups == "" {
				assertGroupsAbsent(t)
			} else {
				assertGroupsUnchanged(t, before)
			}
		})
	}
}

func TestRunListEditGroupMovesAndWraps(t *testing.T) {
	enter := tea.KeyPressMsg{Code: tea.KeyEnter}
	// widget is selected. Choices are all, shell, dev, orphans, with orphans selected.
	open := []tea.KeyPressMsg{{Text: "e"}, enter, enter}
	tests := []struct {
		name string
		keys []tea.KeyPressMsg
		want string
	}{
		{name: "j wraps from the last group", keys: append(append([]tea.KeyPressMsg{}, open...), tea.KeyPressMsg{Text: "j"}), want: "all"},
		{name: "down wraps from the last group", keys: append(append([]tea.KeyPressMsg{}, open...), tea.KeyPressMsg{Code: tea.KeyDown}), want: "all"},
		{name: "k moves up", keys: append(append([]tea.KeyPressMsg{}, open...), tea.KeyPressMsg{Text: "k"}), want: "dev"},
		{name: "up moves up", keys: append(append([]tea.KeyPressMsg{}, open...), tea.KeyPressMsg{Code: tea.KeyUp}), want: "dev"},
		{name: "k wraps from all", keys: []tea.KeyPressMsg{{Text: "j"}, {Text: "e"}, enter, enter, {Text: "k"}}, want: "dev"},
		{name: "j moves down from all", keys: []tea.KeyPressMsg{{Text: "j"}, {Text: "e"}, enter, enter, {Text: "j"}}, want: "shell"},
	}
	manifest := `{
  "manifest": [
    {"name":"widget","group":"orphans","install_type":"omarchy","install_command":"widget"},
    {"name":"zsh","group":"all","install_type":"pacman","install_command":"zsh"}
  ]
}`
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeManifest(t, manifest)
			before := writeGroups(t, storedGroups)
			view, err := RunList("", "", true, tt.keys)
			if err != nil {
				t.Fatal(err)
			}
			assertGroupsUnchanged(t, before)
			plain := stripANSI(view)
			if got := selectedChoice(t, plain); got != tt.want {
				t.Fatalf("selected %s, want %s\n%s", got, tt.want, plain)
			}
		})
	}
}

const editPairManifest = `{
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
    }
  ]
}
`

func TestRunListEditClearingTheCommandUsesTheNewName(t *testing.T) {
	writeManifest(t, editPairManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Code: 'u', Mod: tea.ModCtrl},
		{Text: "nvim"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: 'u', Mod: tea.ModCtrl},
		{Code: tea.KeyEnter},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertBackOnList(t, view)
	if got := selectedEntry(t, view, []string{"nvim", "spotify"}); got != "nvim" {
		t.Fatalf("selected %s, want nvim", got)
	}
	assertManifestJSON(t, `{
  "manifest": [
    {
      "name": "nvim",
      "group": "all",
      "install_type": "pacman",
      "install_command": "nvim"
    },
    {
      "name": "spotify",
      "group": "media",
      "install_type": "flatpak",
      "install_command": "com.spotify.Client"
    }
  ]
}`)
}

func TestRunListEditCaseOnlyRenameReplacesTheEntry(t *testing.T) {
	writeManifest(t, editPairManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Code: 'u', Mod: tea.ModCtrl},
		{Text: "Neovim"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertBackOnList(t, view)
	if got := selectedEntry(t, view, []string{"Neovim", "spotify"}); got != "Neovim" {
		t.Fatalf("selected %s, want Neovim", got)
	}
	if !strings.Contains(stripANSI(entryLine(t, view, "Neovim")), "\uEA85") {
		t.Fatal("case-only rename lost the different install command")
	}
	assertManifestJSON(t, `{
  "manifest": [
    {
      "name": "Neovim",
      "group": "all",
      "install_type": "pacman",
      "install_command": "neovim"
    },
    {
      "name": "spotify",
      "group": "media",
      "install_type": "flatpak",
      "install_command": "com.spotify.Client"
    }
  ]
}`)
}

func TestRunListEditWritesEveryFieldOnce(t *testing.T) {
	writeManifest(t, editPairManifest)
	beforeGroups := writeGroups(t, `{
  "groups": [
    {"name": "shell"}
  ]
}
`)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Code: tea.KeyEnter},
		{Text: "j"},
		{Code: tea.KeyEnter},
		{Text: "j"},
		{Code: tea.KeyEnter},
		{Code: 'u', Mod: tea.ModCtrl},
		{Text: "nvim-bin"},
		{Code: tea.KeyEnter},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertGroupsUnchanged(t, beforeGroups)
	assertBackOnList(t, view)
	if got := selectedEntry(t, view, []string{"neovim", "spotify"}); got != "neovim" {
		t.Fatalf("selected %s, want neovim", got)
	}
	assertManifestJSON(t, `{
  "manifest": [
    {
      "name": "neovim",
      "group": "shell",
      "install_type": "aur",
      "install_command": "nvim-bin"
    },
    {
      "name": "spotify",
      "group": "media",
      "install_type": "flatpak",
      "install_command": "com.spotify.Client"
    }
  ]
}`)
}

func TestRunListEditDoesNotCreateAGroup(t *testing.T) {
	writeManifest(t, `{
  "manifest": [
    {
      "name": "widget",
      "group": "orphans",
      "install_type": "omarchy",
      "install_command": "widget"
    }
  ]
}
`)
	beforeGroups := writeGroups(t, `{
  "groups": [
    {"name": "shell"}
  ]
}
`)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Code: 'u', Mod: tea.ModCtrl},
		{Text: "gadget"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertGroupsUnchanged(t, beforeGroups)
	assertBackOnList(t, view)
	assertManifestJSON(t, `{
  "manifest": [
    {
      "name": "gadget",
      "group": "orphans",
      "install_type": "omarchy",
      "install_command": "widget"
    }
  ]
}`)
}

func TestRunListEditReturnsToTheEntryInsideTheFilter(t *testing.T) {
	writeManifest(t, editPairManifest)
	view, err := RunList("media", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Code: 'u', Mod: tea.ModCtrl},
		{Text: "music"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: 'u', Mod: tea.ModCtrl},
		{Code: tea.KeyEnter},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertBackOnList(t, view)
	plain := stripANSI(view)
	if strings.Contains(plain, "neovim") {
		t.Fatalf("group filter was forgotten\n%s", plain)
	}
	if !strings.Contains(entryLine(t, view, "music"), colorSelected+"music") {
		t.Fatalf("edited entry is not selected\n%s", view)
	}
	assertGroupsAbsent(t)
	assertManifestJSON(t, `{
  "manifest": [
    {
      "name": "neovim",
      "group": "all",
      "install_type": "pacman",
      "install_command": "neovim"
    },
    {
      "name": "music",
      "group": "media",
      "install_type": "flatpak",
      "install_command": "music"
    }
  ]
}`)
}

func assertBackOnList(t *testing.T, view string) {
	t.Helper()
	if strings.Contains(view, "What is the ") || strings.Contains(view, "Which group?") {
		t.Fatalf("edit did not return to the list\n%s", view)
	}
	if !strings.Contains(view, "j/k move • d remove • e edit • q quit") {
		t.Fatalf("list hint missing\n%s", view)
	}
}

func assertManifestJSON(t *testing.T, want string) {
	t.Helper()
	path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "dotpak", "manifest.json")
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("manifest\n%s\nwant\n%s", got, want)
	}
}

func TestRunListEditCommandStartsOnTheCurrentValue(t *testing.T) {
	before := writeManifest(t, editTypesManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "j"},
		{Text: "e"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertManifestUnchanged(t, before)
	plain := stripANSI(view)
	if !strings.Contains(plain, "What is the install command?") {
		t.Fatalf("install command step missing\n%s", plain)
	}
	if !strings.Contains(plain, "node@20") {
		t.Fatalf("install command does not start on the current value\n%s", plain)
	}
	if strings.Contains(plain, "Which group?") {
		t.Fatalf("still on the group step\n%s", plain)
	}
}

func TestRunListEditCommandStepTakesLetters(t *testing.T) {
	before := writeManifest(t, oneEntryManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Text: "j"},
		{Text: "k"},
		{Text: "d"},
		{Text: "e"},
		{Code: 'x'},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertManifestUnchanged(t, before)
	plain := stripANSI(view)
	if !strings.Contains(plain, "What is the install command?") {
		t.Fatalf("left the install command step\n%s", plain)
	}
	if !strings.Contains(plain, "neovimjkde") {
		t.Fatalf("install command did not take j, k, d, and e as characters\n%s", plain)
	}
	if strings.Contains(plain, "neovimjkdex") {
		t.Fatalf("a code-only key inserted a character\n%s", plain)
	}
}

func choiceLines(plain string) []string {
	var got []string
	for _, line := range strings.Split(plain, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Which ") || strings.Contains(line, "move") {
			continue
		}
		line = strings.TrimPrefix(line, "- ")
		got = append(got, strings.TrimSpace(line))
	}
	return got
}

func writeGroups(t *testing.T, body string) []byte {
	t.Helper()
	dir := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "dotpak")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	payload := []byte(body)
	if err := os.WriteFile(filepath.Join(dir, "groups.json"), payload, 0o644); err != nil {
		t.Fatal(err)
	}
	return payload
}

func assertGroupsAbsent(t *testing.T) {
	t.Helper()
	path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "dotpak", "groups.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("groups file appeared: %v", err)
	}
}

func assertGroupsUnchanged(t *testing.T, before []byte) {
	t.Helper()
	path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "dotpak", "groups.json")
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatalf("groups changed\nbefore %s\nafter %s", before, after)
	}
}

func selectedChoice(t *testing.T, plain string) string {
	t.Helper()
	var got string
	for _, line := range strings.Split(plain, "\n") {
		if !strings.HasPrefix(line, "- ") {
			continue
		}
		if got != "" {
			t.Fatalf("two choices selected in %q", plain)
		}
		got = strings.TrimSpace(strings.TrimPrefix(line, "- "))
	}
	if got == "" {
		t.Fatalf("no selected choice in %q", plain)
	}
	return got
}

func assertChoiceOrder(t *testing.T, plain string, choices []string) {
	t.Helper()
	prev := -1
	for _, choice := range choices {
		at := strings.Index(plain, choice)
		if at < 0 || at <= prev {
			t.Fatalf("choices %v not in order in %q", choices, plain)
		}
		prev = at
	}
}

const dropManifest = `{
  "manifest": [
    {
      "name": "alpha",
      "group": "all",
      "install_type": "pacman",
      "install_command": "alpha"
    },
    {
      "name": "beta",
      "group": "all",
      "install_type": "pacman",
      "install_command": "beta"
    },
    {
      "name": "gamma",
      "group": "all",
      "install_type": "mise",
      "install_command": "gamma"
    },
    {
      "name": "delta",
      "group": "media",
      "install_type": "pacman",
      "install_command": "delta"
    }
  ]
}
`

const mediaGroup = `{
  "groups": [
    {"name": "media"}
  ]
}
`

func TestRunListEditDropsAnEntryThatLeavesTheFilter(t *testing.T) {
	enter := tea.KeyPressMsg{Code: tea.KeyEnter}
	toMedia := []tea.KeyPressMsg{{Text: "e"}, enter, enter, {Text: "j"}, enter, enter}
	tests := []struct {
		name        string
		group       string
		installType string
		keys        []tea.KeyPressMsg
		selected    string
		visible     []string
		absent      []string
		empty       bool
		wantJSON    string
	}{
		{
			name:     "next entry after the first drops",
			group:    "all",
			keys:     toMedia,
			selected: "beta",
			visible:  []string{"beta", "gamma"},
			absent:   []string{"alpha", "delta"},
			wantJSON: dropJSON("media", "pacman", "all", "pacman", "all", "mise", "media", "pacman"),
		},
		{
			name:     "next entry after the middle drops",
			group:    "all",
			keys:     append([]tea.KeyPressMsg{{Text: "j"}}, toMedia...),
			selected: "gamma",
			visible:  []string{"alpha", "gamma"},
			absent:   []string{"beta", "delta"},
			wantJSON: dropJSON("all", "pacman", "media", "pacman", "all", "mise", "media", "pacman"),
		},
		{
			name:     "new last entry after the last drops",
			group:    "all",
			keys:     append([]tea.KeyPressMsg{{Text: "j"}, {Text: "j"}}, toMedia...),
			selected: "beta",
			visible:  []string{"alpha", "beta"},
			absent:   []string{"gamma", "delta"},
			wantJSON: dropJSON("all", "pacman", "all", "pacman", "media", "mise", "media", "pacman"),
		},
		{
			name:        "install type filter keeps the next entry",
			installType: "pacman",
			keys:        []tea.KeyPressMsg{{Text: "e"}, enter, {Text: "j"}, enter, enter, enter},
			selected:    "beta",
			visible:     []string{"beta", "delta"},
			absent:      []string{"alpha", "gamma"},
			wantJSON:    dropJSON("all", "aur", "all", "pacman", "all", "mise", "media", "pacman"),
		},
		{
			name:        "both filters hide the other rows",
			group:       "all",
			installType: "pacman",
			keys:        []tea.KeyPressMsg{{Text: "e"}, enter, {Text: "j"}, enter, enter, enter},
			selected:    "beta",
			visible:     []string{"beta"},
			absent:      []string{"alpha", "gamma", "delta"},
			wantJSON:    dropJSON("all", "aur", "all", "pacman", "all", "mise", "media", "pacman"),
		},
		{
			name:     "empty list when the last match drops",
			group:    "media",
			keys:     []tea.KeyPressMsg{{Text: "e"}, enter, enter, {Text: "k"}, enter, enter},
			empty:    true,
			wantJSON: dropJSON("all", "pacman", "all", "pacman", "all", "mise", "all", "pacman"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeManifest(t, dropManifest)
			writeGroups(t, mediaGroup)
			view, err := RunList(tt.group, tt.installType, true, tt.keys)
			if err != nil {
				t.Fatal(err)
			}
			assertManifestJSON(t, tt.wantJSON)
			plain := stripANSI(view)
			if tt.empty {
				if view != emptyListView {
					t.Fatalf("view\n%q\nwant\n%q", view, emptyListView)
				}
				return
			}
			assertBackOnList(t, view)
			if got := selectedEntry(t, view, tt.visible); got != tt.selected {
				t.Fatalf("selected %s, want %s\n%s", got, tt.selected, plain)
			}
			for _, name := range tt.absent {
				if strings.Contains(plain, name) {
					t.Fatalf("entry %s stayed on the list\n%s", name, plain)
				}
			}
		})
	}
}

func dropJSON(alphaGroup, alphaType, betaGroup, betaType, gammaGroup, gammaType, deltaGroup, deltaType string) string {
	return `{
  "manifest": [
    {
      "name": "alpha",
      "group": "` + alphaGroup + `",
      "install_type": "` + alphaType + `",
      "install_command": "alpha"
    },
    {
      "name": "beta",
      "group": "` + betaGroup + `",
      "install_type": "` + betaType + `",
      "install_command": "beta"
    },
    {
      "name": "gamma",
      "group": "` + gammaGroup + `",
      "install_type": "` + gammaType + `",
      "install_command": "gamma"
    },
    {
      "name": "delta",
      "group": "` + deltaGroup + `",
      "install_type": "` + deltaType + `",
      "install_command": "delta"
    }
  ]
}`
}

func TestRunListEditEscDiscardsTheEdit(t *testing.T) {
	enter := tea.KeyPressMsg{Code: tea.KeyEnter}
	esc := tea.KeyPressMsg{Code: tea.KeyEscape}
	tests := []struct {
		name string
		keys []tea.KeyPressMsg
	}{
		{name: "name", keys: []tea.KeyPressMsg{{Text: "e"}, {Text: "x"}, esc}},
		{name: "install type", keys: []tea.KeyPressMsg{{Text: "e"}, enter, {Text: "j"}, esc}},
		{name: "group", keys: []tea.KeyPressMsg{{Text: "e"}, enter, enter, {Text: "j"}, esc}},
		{name: "install command", keys: []tea.KeyPressMsg{{Text: "e"}, enter, enter, enter, {Text: "x"}, esc}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := writeManifest(t, editPairManifest)
			writeGroups(t, `{
  "groups": [
    {"name": "shell"}
  ]
}
`)
			baseline, err := RunList("", "", true, []tea.KeyPressMsg{{Text: "q"}})
			if err != nil {
				t.Fatal(err)
			}
			view, err := RunList("", "", true, tt.keys)
			if err != nil {
				t.Fatal(err)
			}
			assertManifestUnchanged(t, before)
			if view != baseline {
				t.Fatalf("view\n%s\nwant\n%s", view, baseline)
			}
		})
	}
}

func TestRunListEditCtrlCQuitsWithoutSaving(t *testing.T) {
	enter := tea.KeyPressMsg{Code: tea.KeyEnter}
	quit := tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
	tests := []struct {
		name string
		keys []tea.KeyPressMsg
	}{
		{name: "name step", keys: []tea.KeyPressMsg{{Text: "e"}, quit}},
		{
			name: "command step",
			keys: []tea.KeyPressMsg{
				{Text: "e"},
				{Code: 'u', Mod: tea.ModCtrl},
				{Text: "nvim"},
				enter,
				{Text: "j"},
				enter,
				enter,
				quit,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := writeManifest(t, oneEntryManifest)
			view, err := RunList("", "", true, tt.keys)
			if err != nil {
				t.Fatalf("ctrl+c returned %v", err)
			}
			assertManifestUnchanged(t, before)
			if strings.Contains(stripANSI(view), "nvim") && strings.Contains(view, "j/k move • d remove") {
				t.Fatal("ctrl+c saved the edit and returned to the list")
			}
		})
	}
}

func TestRunListEditFailedSaveKeepsTheError(t *testing.T) {
	before := writeManifest(t, oneEntryManifest)
	path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "dotpak", "manifest.json")
	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatal(err)
	}
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Code: 'u', Mod: tea.ModCtrl},
		{Text: "nvim"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Text: "j"},
		{Text: "k"},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertManifestUnchanged(t, before)
	if !strings.Contains(view, "permission denied") {
		t.Fatalf("save error missing\n%s", view)
	}
	if strings.Contains(view, "j/k move • d remove • e edit • q quit") {
		t.Fatalf("j/k left the error\n%s", view)
	}
}

func TestRunListEditEscLeavesAFailedSave(t *testing.T) {
	before := writeManifest(t, oneEntryManifest)
	path := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "dotpak", "manifest.json")
	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatal(err)
	}
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Code: 'u', Mod: tea.ModCtrl},
		{Text: "nvim"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Text: "j"},
		{Code: tea.KeyEscape},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertManifestUnchanged(t, before)
	assertBackOnList(t, view)
	if strings.Contains(view, "permission denied") {
		t.Fatalf("error stayed after esc\n%s", view)
	}
	if got := selectedEntry(t, view, []string{"neovim"}); got != "neovim" {
		t.Fatalf("selected %s", got)
	}
}
