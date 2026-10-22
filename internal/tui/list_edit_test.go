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
