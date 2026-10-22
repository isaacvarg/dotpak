package tui

import (
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
