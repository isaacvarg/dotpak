package tui

import (
	"charm.land/bubbles/v2/textinput"
)

func newCommandInput() textinput.Model {
	command := textinput.New()
	command.SetWidth(40)
	command.Prompt = "› "
	command.Focus()

	return command
}

func (m model) viewCommand() ([]string, int) {
	var output []string

	output = append(output, questionStyle.Render("What is the install command?"), "")
	inputLine := len(output)
	output = append(output, m.commandInput.View())
	output = append(output, "", helpStyle.Render("leave blank to default to name: ", m.name))
	output = append(output, "", controlsStyle.Render("enter confirm • esc quit"))

	return output, inputLine
}

// state stuff
func (m model) setCommand(input string) model {
	command := input
	if input == "" {
		command = m.name
	}

	m.command = command
	m.commandInput.Blur()
	m.state = stateDone
	return m
}
