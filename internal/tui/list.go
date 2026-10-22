package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/isaacvarg/dotpak/internal/manifest"
)

const (
	listNameWidth = 32
	listHint      = "j/k move • / query • d remove • e edit • q quit"
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
	m.filterGroup = group
	m.filterInstallType = installType
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
	entries           []manifest.ManifestEntry
	cursor            int
	groupWidth        int
	typeWidth         int
	filterGroup       string
	filterInstallType string
	height            int
	width             int
	query             string
	prompting         bool
	prompt            textinput.Model
	snapshotQuery     string
	snapshotCursor    int
	confirming        bool
	edit              *listEdit
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
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.height = size.Height
		m.width = size.Width
		return m, nil
	}
	if m.edit != nil {
		return m.updateEdit(msg)
	}
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		if m.prompting {
			return m.deliverPrompt(msg)
		}
		return m, nil
	}
	if m.confirming {
		return m.answerRemove(key)
	}
	if m.prompting {
		return m.updatePrompt(key)
	}
	switch key.String() {
	case "q", "esc", "ctrl+c":
		return m, tea.Quit
	case "j", "down":
		return m.move(1), nil
	case "k", "up":
		return m.move(-1), nil
	case "d":
		return m.askRemove(), nil
	case "e":
		if !m.cursorVisible() {
			return m, nil
		}
		return m.beginEdit(), nil
	case "/":
		return m.beginQuery(), nil
	default:
		return m, nil
	}
}

func (m listModel) move(delta int) listModel {
	visible := m.visibleIndices()
	if len(visible) == 0 {
		return m
	}
	pos := -1
	for i, idx := range visible {
		if idx == m.cursor {
			pos = i
			break
		}
	}
	if pos < 0 {
		return m
	}
	pos += delta
	if pos < 0 {
		pos = len(visible) - 1
	} else if pos >= len(visible) {
		pos = 0
	}
	m.cursor = visible[pos]
	return m
}

func (m listModel) View() tea.View {
	content := m.render()
	if m.edit != nil {
		content = m.renderEdit()
	}
	view := tea.NewView(content)
	if m.height > 0 {
		view.AltScreen = true
	}
	if m.prompting && m.edit == nil {
		_, x := queryLine("/", m.prompt.Value(), m.prompt.Position(), m.width)
		view.Cursor = tea.NewCursor(x, strings.Count(content, "\n"))
	}
	return view
}

func (m listModel) render() string {
	if m.confirming {
		return m.removePrompt()
	}
	visible := m.visibleIndices()
	pos := -1
	for i, idx := range visible {
		if idx == m.cursor {
			pos = i
			break
		}
	}
	start, end := 0, len(visible)
	if m.height > 0 {
		room := m.height - 3
		if m.query != "" && !m.prompting {
			room--
		}
		if room < 1 {
			room = 1
		}
		if pos < 0 {
			if room > len(visible) {
				room = len(visible)
			}
			end = room
		} else {
			start, end = rowWindow(pos, len(visible), room)
		}
	}
	var b strings.Builder
	b.WriteString(m.headerLine())
	b.WriteByte('\n')
	for _, i := range visible[start:end] {
		b.WriteString(m.rowLine(i, m.entries[i]))
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	if m.prompting {
		line, _ := queryLine("/", m.prompt.Value(), m.prompt.Position(), m.width)
		b.WriteString(line)
		return b.String()
	}
	if m.query != "" {
		line, _ := queryLine("/", m.query, len([]rune(m.query)), m.width)
		b.WriteString(line)
		b.WriteByte('\n')
	}
	b.WriteString(listHint)
	return b.String()
}

// rowWindow is the slice of rows that fit under the header and above the hint,
// always including the selected row.
func rowWindow(cursor, n, room int) (start, end int) {
	if n == 0 {
		return 0, 0
	}
	if room > n {
		room = n
	}
	start = cursor - room + 1
	if start < 0 {
		start = 0
	}
	end = start + room
	if end > n {
		end = n
		start = end - room
	}
	return start, end
}

func (m listModel) headerLine() string {
	return padRight("Name", listNameWidth) + "  " +
		padRight("Group", m.groupWidth) + "  " +
		padRight("Install type", m.typeWidth) + "  " +
		"\uEA85"
}

func (m listModel) rowLine(i int, e manifest.ManifestEntry) string {
	name := padRight(cutEntryName(e.Name), listNameWidth)
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

func cutEntryName(name string) string {
	if lipgloss.Width(name) <= listNameWidth {
		return name
	}
	return ansi.Truncate(name, listNameWidth, "…")
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
