package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func newQueryPrompt() textinput.Model {
	prompt := textinput.New()
	prompt.Prompt = "/"
	prompt.Focus()
	return prompt
}

func (m listModel) beginQuery() listModel {
	m.snapshotQuery = m.query
	m.snapshotCursor = m.cursor
	m.prompting = true
	m.prompt = newQueryPrompt()
	return m
}

func (m listModel) acceptQuery() listModel {
	if m.prompt.Value() == "" {
		m.query = ""
		m.cursor = m.snapshotCursor
	} else {
		m.query = m.prompt.Value()
	}
	m.prompting = false
	return m
}

func (m listModel) abortQuery() listModel {
	m.query = m.snapshotQuery
	m.cursor = m.snapshotCursor
	m.prompting = false
	return m
}

func (m listModel) deliverPrompt(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.prompt, cmd = m.prompt.Update(msg)
	return m.syncPromptCursor(), cmd
}

func (m listModel) syncPromptCursor() listModel {
	if m.prompt.Value() == "" {
		m.cursor = m.snapshotCursor
		return m
	}
	return m.settleCursor(m.prompt.Value())
}

func (m listModel) updatePrompt(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		return m.abortQuery(), nil
	case "enter":
		return m.acceptQuery(), nil
	case "up", "down", "ctrl+p", "ctrl+n":
		return m, nil
	}
	if m.prompt.Value() == "" && (key.String() == "backspace" || key.String() == "ctrl+h") {
		return m.abortQuery(), nil
	}
	var cmd tea.Cmd
	m.prompt, cmd = m.prompt.Update(key)
	return m.syncPromptCursor(), cmd
}

func (m listModel) settleCursor(query string) listModel {
	if query == "" || len(m.entries) == 0 {
		return m
	}
	if m.cursor < 0 || m.cursor >= len(m.entries) {
		m.cursor = 0
	}
	if containsFold(m.entries[m.cursor].Name, query) {
		return m
	}
	n := len(m.entries)
	for step := 1; step <= n; step++ {
		i := (m.cursor + step) % n
		if containsFold(m.entries[i].Name, query) {
			m.cursor = i
			return m
		}
	}
	return m
}

func (m listModel) activeQuery() string {
	if !m.prompting {
		return m.query
	}
	if m.prompt.Value() == "" {
		return m.snapshotQuery
	}
	return m.prompt.Value()
}

func (m listModel) visibleIndices() []int {
	query := m.activeQuery()
	indices := make([]int, 0, len(m.entries))
	for i, entry := range m.entries {
		if containsFold(entry.Name, query) {
			indices = append(indices, i)
		}
	}
	return indices
}

func (m listModel) cursorVisible() bool {
	if len(m.entries) == 0 || m.cursor < 0 || m.cursor >= len(m.entries) {
		return false
	}
	return containsFold(m.entries[m.cursor].Name, m.activeQuery())
}

func containsFold(name, query string) bool {
	if query == "" {
		return true
	}
	names := []rune(name)
	queries := []rune(query)
	if len(queries) > len(names) {
		return false
	}
	for i := 0; i+len(queries) <= len(names); i++ {
		if strings.EqualFold(string(names[i:i+len(queries)]), query) {
			return true
		}
	}
	return false
}

// queryLine draws "/" and the query on one line. A known width cuts the head
// of a long query with an ellipsis and keeps the tail. While a cursor is
// shown it stays in the last on-screen cell, including the cell after the text.
func queryLine(prefix, text string, pos, width int, showCursor bool) (string, int) {
	runes := []rune(text)
	if pos < 0 {
		pos = 0
	}
	if pos > len(runes) {
		pos = len(runes)
	}
	prefixW := lipgloss.Width(prefix)
	if width <= 0 {
		return prefix + text, prefixW + lipgloss.Width(string(runes[:pos]))
	}
	avail := width - prefixW
	if avail <= 0 {
		return ansi.Truncate(prefix, width, ""), 0
	}
	caret := 0
	if showCursor && pos == len(runes) {
		caret = 1
	}
	if lipgloss.Width(text)+caret <= avail {
		x := prefixW + lipgloss.Width(string(runes[:pos]))
		if x >= width {
			x = width - 1
		}
		return prefix + text, x
	}
	ellipsis := "…"
	ellipsisW := lipgloss.Width(ellipsis)
	budget := avail - ellipsisW - caret
	if budget < 1 {
		budget = avail - caret
		if budget < 1 {
			budget = 1
		}
		ellipsis = ""
		ellipsisW = 0
	}
	start := 0
	for start < pos && lipgloss.Width(string(runes[start:pos])) > budget {
		start++
	}
	end := len(runes)
	for end > pos && lipgloss.Width(string(runes[start:end])) > budget {
		end--
	}
	for start < pos && lipgloss.Width(string(runes[start:end])) > budget {
		start++
	}
	if start == 0 {
		ellipsis = ""
		ellipsisW = 0
	}
	shown := string(runes[start:end])
	cursorX := prefixW + lipgloss.Width(string(runes[start:pos]))
	if start > 0 {
		shown = ellipsis + shown
		cursorX += ellipsisW
	}
	if cursorX >= width {
		cursorX = width - 1
	}
	return prefix + shown, cursorX
}
