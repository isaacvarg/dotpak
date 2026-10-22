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
	return "remove " + m.entries[m.cursor].Name + "?"
}

func (m listModel) answerRemove(key tea.KeyPressMsg) (tea.Model, tea.Cmd, bool) {
	switch key.String() {
	case "y", "Y":
		return m.removeSelected()
	case "j", "k", "up", "down":
		return m, nil, true
	case "q":
		return m.cancelRemove(), nil, true
	default:
		return m, nil, false
	}
}

func (m listModel) removeSelected() (tea.Model, tea.Cmd, bool) {
	name := m.entries[m.cursor].Name
	mf, err := manifest.Load()
	if err != nil {
		return m.cancelRemove(), nil, true
	}
	if err := mf.Remove(name); err != nil {
		return m.cancelRemove(), nil, true
	}
	if err := mf.Save(); err != nil {
		return m.cancelRemove(), nil, true
	}
	m.entries = append(m.entries[:m.cursor], m.entries[m.cursor+1:]...)
	m.confirming = false
	return m, nil, true
}
