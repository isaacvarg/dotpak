package tui

func (m model) selectGroup() model {
	m.group = m.groups[m.groupIdx]
	m.state = stateCommand
	return m
}

func (m model) viewGroup() []string {
	var output []string

	output = append(output,
		questionStyle.Render("Which group does it belong to?"),
		"",
		m.selection(m.groups, m.groupIdx, nil),
		controlsStyle.Render("j/k move • enter select"),
	)
	return output
}
