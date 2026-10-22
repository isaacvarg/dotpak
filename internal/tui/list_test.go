package tui

import (
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
)

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
