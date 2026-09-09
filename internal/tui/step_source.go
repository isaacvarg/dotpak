package tui

func (m model) selectSource() model {
	m.source = m.sources[m.sourceIdx]
	m.state = stateGroup
	return m
}

func (m model) viewSource() []string {
	var output []string

	output = append(output,
		questionStyle.Render("What is the installation type?"),
		"",
		m.selection(m.sources, m.sourceIdx, sourceColors),
		controlsStyle.Render("j/k move • enter select"),
	)
	return output
}
