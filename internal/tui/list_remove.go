package tui

import tea "charm.land/bubbletea/v2"

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
	case "j", "k", "up", "down":
		return m, nil, true
	case "q":
		return m.cancelRemove(), nil, true
	default:
		return m, nil, false
	}
}
