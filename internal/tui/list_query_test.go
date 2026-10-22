package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/isaacvarg/dotpak/internal/manifest"
)

func TestRunListQueryFiltersTheLoadedNames(t *testing.T) {
	before := writeManifest(t, filteredManifest)
	baseline, err := RunList("", "", true, []tea.KeyPressMsg{{Text: "q"}})
	if err != nil {
		t.Fatal(err)
	}
	opened, err := RunList("", "", true, []tea.KeyPressMsg{{Text: "/"}})
	if err != nil {
		t.Fatal(err)
	}
	if opened == baseline {
		t.Fatal("slash left the list unchanged")
	}
	if strings.Contains(opened, listHint) {
		t.Fatalf("hint shown while the prompt is open\n%s", opened)
	}
	if lastLine(opened) != "/" {
		t.Fatalf("prompt %q", lastLine(opened))
	}
	for _, name := range []string{"neovim", "spotify", "node", "ripgrep"} {
		if !strings.Contains(stripANSI(opened), name) {
			t.Fatalf("opening slash hid %s\n%s", name, opened)
		}
	}

	aborted, err := RunList("", "", true, []tea.KeyPressMsg{{Text: "/"}, {Code: tea.KeyEscape}})
	if err != nil {
		t.Fatal(err)
	}
	if aborted != baseline {
		t.Fatalf("esc did not restore the list\n%s\nwant\n%s", aborted, baseline)
	}
	assertManifestUnchanged(t, before)

	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "n"},
		{Text: "o"},
		{Text: "d"},
		{Text: "e"},
		{Code: tea.KeyEnter},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := selectedEntry(t, view, []string{"node"}); got != "node" {
		t.Fatalf("selected %s", got)
	}
	for _, name := range []string{"neovim", "spotify", "ripgrep"} {
		if strings.Contains(stripANSI(view), name) {
			t.Fatalf("query node shows %s\n%s", name, stripANSI(view))
		}
	}
	lines := strings.Split(view, "\n")
	if lines[len(lines)-1] != listHint {
		t.Fatalf("hint %q", lines[len(lines)-1])
	}
	if lines[len(lines)-2] != "/node" {
		t.Fatalf("accepted query %q", lines[len(lines)-2])
	}
	if strings.Contains(lines[len(lines)-2], "\x1b") {
		t.Fatalf("query line is styled %q", lines[len(lines)-2])
	}
	assertManifestUnchanged(t, before)
}

func TestRunListQueryMatchesNamesLiterally(t *testing.T) {
	writeManifest(t, filteredManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "N"},
		{Text: "e"},
		{Text: "O"},
		{Code: tea.KeyEnter},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stripANSI(view), "neovim") || strings.Contains(stripANSI(view), "node") {
		t.Fatalf("case-insensitive name match\n%s", stripANSI(view))
	}
	if !strings.Contains(view, "/NeO\n") {
		t.Fatalf("query line did not keep the typed case\n%s", view)
	}

	hidden, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "r"},
		{Text: "g"},
		{Code: tea.KeyEnter},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stripANSI(hidden), "ripgrep") {
		t.Fatalf("non-contiguous letters matched\n%s", stripANSI(hidden))
	}

	command, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "C"},
		{Text: "l"},
		{Text: "i"},
		{Text: "e"},
		{Text: "n"},
		{Text: "t"},
		{Code: tea.KeyEnter},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stripANSI(command), "spotify") {
		t.Fatalf("install command matched\n%s", stripANSI(command))
	}
}

func TestRunListQuerySelectionFollowsTheLoadedOrder(t *testing.T) {
	toNode := []tea.KeyPressMsg{{Text: "j"}, {Text: "j"}, {Text: "/"}}
	tests := []struct {
		name     string
		extra    []tea.KeyPressMsg
		selected string
		visible  []string
		absent   []string
	}{
		{
			name:     "e stays on node",
			extra:    []tea.KeyPressMsg{{Text: "e"}},
			selected: "node",
			visible:  []string{"neovim", "node", "ripgrep"},
			absent:   []string{"spotify"},
		},
		{
			name:     "p moves to ripgrep",
			extra:    []tea.KeyPressMsg{{Text: "p"}},
			selected: "ripgrep",
			visible:  []string{"spotify", "ripgrep"},
			absent:   []string{"neovim", "node"},
		},
		{
			name:     "v wraps to neovim",
			extra:    []tea.KeyPressMsg{{Text: "v"}},
			selected: "neovim",
			visible:  []string{"neovim"},
			absent:   []string{"spotify", "node", "ripgrep"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeManifest(t, filteredManifest)
			view, err := RunList("", "", true, append(append([]tea.KeyPressMsg{}, toNode...), tt.extra...))
			if err != nil {
				t.Fatal(err)
			}
			if got := selectedEntry(t, view, tt.visible); got != tt.selected {
				t.Fatalf("selected %s, want %s\n%s", got, tt.selected, stripANSI(view))
			}
			plain := stripANSI(view)
			for _, name := range tt.absent {
				if strings.Contains(plain, name) {
					t.Fatalf("showed %s\n%s", name, plain)
				}
			}
		})
	}
}

func TestRunListQueryRemembersTheSelectionThroughNoMatch(t *testing.T) {
	writeManifest(t, filteredManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "j"}, {Text: "j"}, {Text: "j"},
		{Text: "/"},
		{Text: "p"},
		{Text: "z"},
		{Code: tea.KeyBackspace},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := selectedEntry(t, view, []string{"spotify", "ripgrep"}); got != "ripgrep" {
		t.Fatalf("selected %s\n%s", got, stripANSI(view))
	}
	plain := stripANSI(view)
	if strings.Contains(plain, "neovim") || strings.Contains(plain, "node") {
		t.Fatalf("rows\n%s", plain)
	}
}

func TestRunListQueryEmptyPromptRestoresTheSnapshot(t *testing.T) {
	writeManifest(t, filteredManifest)
	baseline, err := RunList("", "", true, []tea.KeyPressMsg{{Text: "j"}, {Text: "q"}})
	if err != nil {
		t.Fatal(err)
	}
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "j"},
		{Text: "/"},
		{Text: "v"},
		{Code: tea.KeyBackspace},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := selectedEntry(t, view, []string{"neovim", "spotify", "node", "ripgrep"}); got != "spotify" {
		t.Fatalf("selected %s\n%s", got, stripANSI(view))
	}
	if strings.Contains(view, listHint) {
		t.Fatal("empty prompt closed")
	}
	if lastLine(view) != "/" {
		t.Fatalf("prompt %q", lastLine(view))
	}
	for _, name := range []string{"neovim", "spotify", "node", "ripgrep"} {
		if !strings.Contains(stripANSI(view), name) {
			t.Fatalf("missing %s\n%s", name, stripANSI(view))
		}
	}

	aborted, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "j"},
		{Text: "/"},
		{Text: "v"},
		{Code: tea.KeyBackspace},
		{Code: tea.KeyBackspace},
	})
	if err != nil {
		t.Fatal(err)
	}
	if aborted != baseline {
		t.Fatalf("second backspace\n%s\nwant\n%s", aborted, baseline)
	}
}

func TestRunListQueryEnterKeepsAndClears(t *testing.T) {
	writeManifest(t, filteredManifest)
	kept, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "j"},
		{Text: "/"},
		{Text: "n"},
		{Text: "o"},
		{Text: "d"},
		{Text: "e"},
		{Code: tea.KeyEnter},
		{Text: "/"},
	})
	if err != nil {
		t.Fatal(err)
	}
	plain := stripANSI(kept)
	if strings.Contains(plain, "neovim") || strings.Contains(plain, "spotify") || strings.Contains(plain, "ripgrep") {
		t.Fatalf("opening a new prompt flashed the loaded list\n%s", plain)
	}
	if strings.Contains(kept, "/node") {
		t.Fatalf("accepted query stayed on screen\n%s", kept)
	}
	if lastLine(kept) != "/" {
		t.Fatalf("prompt %q", lastLine(kept))
	}

	cleared, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "n"},
		{Text: "o"},
		{Text: "d"},
		{Text: "e"},
		{Code: tea.KeyEnter},
		{Text: "/"},
		{Code: tea.KeyEnter},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(cleared, "/node") {
		t.Fatalf("empty enter kept the query\n%s", cleared)
	}
	if got := selectedEntry(t, cleared, []string{"neovim", "spotify", "node", "ripgrep"}); got != "node" {
		t.Fatalf("selected %s, want the entry remembered when slash was pressed", got)
	}
	for _, name := range []string{"neovim", "spotify", "node", "ripgrep"} {
		if !rowShows(cleared, name) {
			t.Fatalf("empty enter hid %s\n%s", name, stripANSI(cleared))
		}
	}
}

func TestRunListQueryEditsLikeAName(t *testing.T) {
	writeManifest(t, filteredManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "a"},
		{Text: "b"},
		{Code: tea.KeyLeft},
		{Text: "c"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if lastLine(view) != "/acb" {
		t.Fatalf("left did not move inside the query %q", lastLine(view))
	}

	home, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "a"},
		{Text: "b"},
		{Code: tea.KeyHome},
		{Text: "c"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if lastLine(home) != "/cab" {
		t.Fatalf("home %q", lastLine(home))
	}

	cleared, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "n"},
		{Text: "o"},
		{Code: 'u', Mod: tea.ModCtrl},
	})
	if err != nil {
		t.Fatal(err)
	}
	if lastLine(cleared) != "/" {
		t.Fatalf("ctrl+u %q", lastLine(cleared))
	}
	if strings.Contains(cleared, listHint) {
		t.Fatal("ctrl+u closed the prompt")
	}

	still, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Code: tea.KeyUp},
		{Code: tea.KeyDown},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := selectedEntry(t, still, []string{"neovim", "spotify", "node", "ripgrep"}); got != "neovim" {
		t.Fatalf("arrows moved the list to %s", got)
	}
	if lastLine(still) != "/" {
		t.Fatalf("arrows typed into the prompt %q", lastLine(still))
	}
}

func TestRunListQueryStaysInsideTheFilter(t *testing.T) {
	writeManifest(t, filteredManifest)
	view, err := RunList("dev", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "n"},
		{Code: tea.KeyEnter},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	plain := stripANSI(view)
	if !strings.Contains(plain, "node") || strings.Contains(plain, "ripgrep") || strings.Contains(plain, "neovim") {
		t.Fatalf("filter and query\n%s", plain)
	}
}

func TestRunListQueryKeepsLoadedColumnWidths(t *testing.T) {
	const manifest = `{
  "manifest": [
    {
      "name": "neovim",
      "group": "all",
      "install_type": "pacman",
      "install_command": "neovim"
    },
    {
      "name": "widget",
      "group": "development",
      "install_type": "pacman",
      "install_command": "widget"
    }
  ]
}`
	writeManifest(t, manifest)
	full, err := RunList("", "", true, []tea.KeyPressMsg{{Text: "q"}})
	if err != nil {
		t.Fatal(err)
	}
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "w"},
		{Code: tea.KeyEnter},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if headerOf(full) != headerOf(view) {
		t.Fatalf("header changed\nfull %q\nquery %q", headerOf(full), headerOf(view))
	}
	if strings.Contains(stripANSI(view), "neovim") || !strings.Contains(stripANSI(view), "widget") {
		t.Fatalf("rows\n%s", stripANSI(view))
	}
}

func TestRunListQueryRemoveAndEditUseTheShownRows(t *testing.T) {
	before := writeManifest(t, filteredManifest)
	removed, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "n"},
		{Text: "o"},
		{Text: "d"},
		{Text: "e"},
		{Code: tea.KeyEnter},
		{Text: "d"},
		{Text: "y"},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"neovim", "spotify", "node", "ripgrep"} {
		if rowShows(removed, name) {
			t.Fatalf("row %s still shown\n%s", name, stripANSI(removed))
		}
	}
	if !strings.Contains(removed, "/node\n"+listHint) {
		t.Fatalf("query did not survive the remove\n%s", removed)
	}
	if strings.Contains(readManifest(t), "node") {
		t.Fatalf("manifest\n%s", readManifest(t))
	}

	writeManifest(t, beforeManifest(before))
	cancelled, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "n"},
		{Text: "o"},
		{Text: "d"},
		{Text: "e"},
		{Code: tea.KeyEnter},
		{Text: "d"},
		{Text: "/"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(cancelled, "remove node?") {
		t.Fatalf("slash did not cancel the remove\n%s", cancelled)
	}
	if !strings.Contains(cancelled, "/node\n"+listHint) {
		t.Fatalf("slash opened a prompt from the remove\n%s", cancelled)
	}
	if !strings.Contains(readManifest(t), "node") {
		t.Fatal("remove was confirmed")
	}

	renamed, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "n"},
		{Text: "e"},
		{Text: "o"},
		{Code: tea.KeyEnter},
		{Text: "e"},
		{Code: 'u', Mod: tea.ModCtrl},
		{Text: "r"}, {Text: "u"}, {Text: "n"}, {Text: "t"}, {Text: "i"}, {Text: "m"}, {Text: "e"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEnter},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if rowShows(renamed, "neovim") || rowShows(renamed, "runtime") {
		t.Fatalf("renamed entry stayed on the query\n%s", stripANSI(renamed))
	}
	if !strings.Contains(renamed, "/neo\n"+listHint) {
		t.Fatalf("query line\n%s", renamed)
	}
	if !strings.Contains(readManifest(t), "runtime") {
		t.Fatalf("manifest\n%s", readManifest(t))
	}
}

func TestRunListQuerySlashInAnEditIsText(t *testing.T) {
	writeManifest(t, oneEntryManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "e"},
		{Text: "/"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stripANSI(view), "neovim/") {
		t.Fatalf("slash was not typed into the name\n%s", stripANSI(view))
	}
	if strings.Contains(view, listHint) {
		t.Fatalf("slash left the edit\n%s", view)
	}
}

func TestRunListQueryPasteFoldsTabsAndNewlines(t *testing.T) {
	m := newListModel([]manifest.ManifestEntry{{
		Name:           "node",
		Group:          "dev",
		InstallType:    manifest.InstallMise,
		InstallCommand: "node@20",
	}})
	opened, _ := m.Update(tea.KeyPressMsg{Text: "/"})
	pasted, _ := opened.Update(tea.PasteMsg{Content: "a\tb\nc"})
	if lastLine(pasted.View().Content) != "/a b c" {
		t.Fatalf("pasted query %q", lastLine(pasted.View().Content))
	}
}

func TestRunListQueryEscapeQuitsWhileAQueryIsActive(t *testing.T) {
	before := writeManifest(t, filteredManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "n"},
		{Text: "o"},
		{Text: "d"},
		{Text: "e"},
		{Code: tea.KeyEnter},
		{Code: tea.KeyEscape},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertManifestUnchanged(t, before)
	if !strings.Contains(view, "/node\n"+listHint) {
		t.Fatalf("esc cleared the query instead of quitting\n%s", view)
	}
	if rowShows(view, "neovim") || !rowShows(view, "node") {
		t.Fatalf("rows\n%s", stripANSI(view))
	}
}

func TestRunListQueryCutsALongLine(t *testing.T) {
	m := newListModel([]manifest.ManifestEntry{{
		Name:           "neovim",
		Group:          "all",
		InstallType:    manifest.InstallPacman,
		InstallCommand: "neovim",
	}})
	sized, _ := m.Update(tea.WindowSizeMsg{Width: 8, Height: 20})
	current := sized
	keys := []tea.KeyPressMsg{{Text: "/"}}
	for _, r := range "0123456789" {
		keys = append(keys, tea.KeyPressMsg{Text: string(r)})
	}
	keys = append(keys, tea.KeyPressMsg{Code: tea.KeyEnter})
	for _, key := range keys {
		next, _ := current.Update(key)
		current = next
	}
	line := strings.Split(current.View().Content, "\n")
	query := line[len(line)-2]
	if query != "/…456789" {
		t.Fatalf("accepted query %q", query)
	}
	if strings.Contains(query, "\x1b") {
		t.Fatalf("styled %q", query)
	}
}

func TestRunListQueryNothingDrawnStillIgnoresRemove(t *testing.T) {
	before := writeManifest(t, filteredManifest)
	view, err := RunList("", "", true, []tea.KeyPressMsg{
		{Text: "/"},
		{Text: "z"},
		{Code: tea.KeyEnter},
		{Text: "d"},
		{Text: "e"},
		{Text: "j"},
		{Text: "q"},
	})
	if err != nil {
		t.Fatal(err)
	}
	assertManifestUnchanged(t, before)
	for _, name := range []string{"neovim", "spotify", "node", "ripgrep"} {
		if rowShows(view, name) {
			t.Fatalf("showed %s\n%s", name, stripANSI(view))
		}
	}
	if !strings.Contains(view, "/z\n"+listHint) {
		t.Fatalf("view\n%s", view)
	}
	if strings.Contains(view, "remove z?") || strings.Contains(view, "What is the name?") {
		t.Fatalf("d or e acted\n%s", view)
	}
}

func lastLine(view string) string {
	lines := strings.Split(view, "\n")
	return lines[len(lines)-1]
}

func headerOf(view string) string {
	return strings.Split(view, "\n")[0]
}

func rowShows(view, name string) bool {
	for _, line := range strings.Split(view, "\n") {
		plain := stripANSI(line)
		if strings.HasPrefix(plain, name) && (len(plain) == len(name) || plain[len(name)] == ' ') {
			return true
		}
	}
	return false
}

func beforeManifest(before []byte) string {
	if before == nil {
		return ""
	}
	return string(before)
}
