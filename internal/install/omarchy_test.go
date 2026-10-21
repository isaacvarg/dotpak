package install

import (
	"strconv"
	"strings"
	"testing"
)

func TestOmarchyVersionReady(t *testing.T) {
	tests := []struct {
		name string
		line string
		ok   bool
	}{
		{name: "release with pkgrel", line: "4.0.0-1", ok: true},
		{name: "newer release", line: "4.0.6-1", ok: true},
		{name: "major only", line: "4", ok: true},
		{name: "major minor", line: "4.0", ok: true},
		{name: "edge pkgver", line: "4.0.0.r1836.g0ae1694-1", ok: true},
		{name: "alpha tail", line: "4.0.0.alpha", ok: true},
		{name: "trimmed command output", line: "  4.0.6-1\n", ok: true},
		{name: "epoch", line: "1:4.0.6-1", ok: true},
		{name: "dev checkout", line: "dev", ok: true},
		{name: "dev checkout hash", line: "dev (abc1234)", ok: true},
		{name: "dev checkout hash uppercase", line: "dev (ABC1234)", ok: true},
		{name: "dev test package", line: "dev.abc1234.dirty-1", ok: true},
		{name: "older release", line: "3.8.5", ok: false},
		{name: "older release with pkgrel", line: "3.8.5-1", ok: false},
		{name: "non-numeric epoch", line: "hello:4.0.0-1", ok: false},
		{name: "empty", line: "", ok: false},
		{name: "unrecognised", line: "nope", ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := omarchyVersionReady(tt.line)
			if tt.ok {
				if err != nil {
					t.Fatalf("line is new enough, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("line is below 4.0.0, got nil")
			}
			got := err.Error()
			if !strings.Contains(got, "version 4.0.0 or greater of omarchy is required") || !strings.Contains(got, strconv.Quote(tt.line)) {
				t.Fatalf("error %q does not report the line", got)
			}
		})
	}
}
