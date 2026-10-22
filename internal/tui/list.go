package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/isaacvarg/dotpak/internal/manifest"
)

const (
	listNameWidth = 32
	listHint      = "j/k move • d remove • e edit • q quit"
)

// RunList draws the manifest list until the session quits.
// An empty group keeps every entry. An empty install type keeps every install type.
// A nil keys slice reads the terminal; any other slice is the scripted session.
func RunList(group, installType string, terminal bool, keys []tea.KeyPressMsg) (string, error) {
	if installType != "" && !manifest.InstallType(installType).IsValid() {
		return "", fmt.Errorf("invalid install type: %q", installType)
	}
	if !terminal {
		return "", fmt.Errorf("stdout is not a terminal")
	}

	mf, err := manifest.Load()
	if err != nil {
		return "", err
	}

	m := newListModel(filterEntries(mf.Manifest, group, installType))
	if keys == nil {
		final, err := tea.NewProgram(m).Run()
		if err != nil {
			return "", err
		}
		return final.View().Content, nil
	}
	return runScripted(m, keys), nil
}

func filterEntries(entries []manifest.ManifestEntry, group, installType string) []manifest.ManifestEntry {
	var out []manifest.ManifestEntry
	for _, e := range entries {
		if group != "" && e.Group != group {
			continue
		}
		if installType != "" && e.InstallType != manifest.InstallType(installType) {
			continue
		}
		out = append(out, e)
	}
	return out
}

type listModel struct {
	entries    []manifest.ManifestEntry
	cursor     int
	groupWidth int
	typeWidth  int
}

func newListModel(entries []manifest.ManifestEntry) listModel {
	groupWidth := lipgloss.Width("Group")
	typeWidth := lipgloss.Width("Install type")
	for _, e := range entries {
		if w := lipgloss.Width(e.Group); w > groupWidth {
			groupWidth = w
		}
		if w := lipgloss.Width(string(e.InstallType)); w > typeWidth {
			typeWidth = w
		}
	}
	return listModel{
		entries:    entries,
		groupWidth: groupWidth,
		typeWidth:  typeWidth,
	}
}

func (m listModel) Init() tea.Cmd { return nil }

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "q", "esc", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m listModel) View() tea.View {
	return tea.NewView(m.render())
}

func (m listModel) render() string {
	var b strings.Builder
	b.WriteString(m.headerLine())
	b.WriteByte('\n')
	for i, e := range m.entries {
		b.WriteString(m.rowLine(i, e))
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	b.WriteString(listHint)
	return b.String()
}

func (m listModel) headerLine() string {
	return padRight("Name", listNameWidth) + "  " +
		padRight("Group", m.groupWidth) + "  " +
		padRight("Install type", m.typeWidth) + "  " +
		"\uEA85"
}

func (m listModel) rowLine(i int, e manifest.ManifestEntry) string {
	name := padRight(e.Name, listNameWidth)
	group := padRight(e.Group, m.groupWidth)
	install := installTypeStyle(e.InstallType).Render(padRight(string(e.InstallType), m.typeWidth))
	mark := " "
	if e.InstallCommand != "" && e.InstallCommand != e.Name {
		mark = "\uEA85"
	}
	if i == m.cursor {
		name = selectedStyle.Render(name)
		group = selectedStyle.Render(group)
		mark = selectedStyle.Render(mark)
	}
	return name + "  " + group + "  " + install + "  " + mark
}

func installTypeStyle(installType manifest.InstallType) lipgloss.Style {
	style := lipgloss.NewStyle()
	if c, ok := sourceColors[string(installType)]; ok {
		style = style.Foreground(c)
	}
	return style
}

func padRight(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

func runScripted(m listModel, keys []tea.KeyPressMsg) string {
	var current tea.Model = m
	for _, key := range keys {
		next, cmd := current.Update(key)
		current = next
		if cmd == nil {
			continue
		}
		if _, ok := cmd().(tea.QuitMsg); ok {
			break
		}
	}
	return current.View().Content
}
