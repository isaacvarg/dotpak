package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/isaacvarg/dotpak/internal/manifest"
)

func (m listModel) askRemove() listModel {
	if len(m.entries) == 0 {
		return m
	}
	m.confirming = true
	return m
}

func (m listModel) cancelRemove() listModel {
	m.confirming = false
	return m
}

func (m listModel) removePrompt() string {
	return "remove " + m.entries[m.cursor].Name + "?\n\n" +
		controlsStyle.Render("y confirm • any other key cancel")
}

func (m listModel) answerRemove(key tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "y", "Y":
		return m.removeSelected()
	default:
		return m.cancelRemove(), nil
	}
}

func (m listModel) removeSelected() (tea.Model, tea.Cmd) {
	name := m.entries[m.cursor].Name
	mf, err := manifest.Load()
	if err != nil {
		return m.cancelRemove(), nil
	}
	if err := mf.Remove(name); err != nil {
		return m.cancelRemove(), nil
	}
	if err := mf.Save(); err != nil {
		return m.cancelRemove(), nil
	}
	m.entries = append(m.entries[:m.cursor], m.entries[m.cursor+1:]...)
	if len(m.entries) > 0 && m.cursor >= len(m.entries) {
		m.cursor = len(m.entries) - 1
	}
	m.confirming = false
	return m, nil
}
