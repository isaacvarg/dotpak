package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
)

// name input stuff
func newNameInput() textinput.Model {
	name := textinput.New()
	name.Placeholder = "zoxide"
	name.SetWidth(40)
	name.Prompt = "› "
	name.Focus()

	return name
}

func (m model) viewName() ([]string, int) {
	var output []string

	output = append(output, questionStyle.Render("What are we adding?"), "")
	inputLine := len(output)
	output = append(output, m.nameInput.View())
	if m.err != "" {
		output = append(output, "", errorStyle.Render("x "+m.err))
	}
	output = append(output, "", controlsStyle.Render("enter confirm • q quit"))

	return output, inputLine
}

// state stuff
func (m model) setName(input string) model {
	name := strings.TrimSpace(input)
	if name == "" {
		m.err = "whoops, a name is required"
		return m
	}

	m.err = ""
	m.name = name
	m.state = stateSource
	return m
}
