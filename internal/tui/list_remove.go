package tui

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
