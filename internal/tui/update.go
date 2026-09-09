package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// update component
// and controls

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// think of Msg as "something happened"
	switch msg := msg.(type) {
	// is it a keypress?
	case tea.KeyPressMsg:

		// cool, which key was pressed
		switch msg.String() {
		// exit
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m.confirm()
		}
	}

	// return the updated model to the bubble tea runtime for processing
	// no command returned
	return m.msgRelay(msg)
}

// relays the msg to the proper input since we separated the update/controllers form the state setters
func (m model) msgRelay(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateName:
		m.nameInput, cmd = m.nameInput.Update(msg)
	}

	return m, cmd
}

// confirms is the  advancer for the state
func (m model) confirm() (tea.Model, tea.Cmd) {
	switch m.state {
	case stateName:
		return m.setName(m.nameInput.Value()), nil
	}

	return m, nil
}

func (m model) cursor(line int) *tea.Cursor {
	if line < 0 {
		return nil
	}

	var cursor *tea.Cursor
	switch m.state {
	case stateName:
		cursor = m.nameInput.Cursor()
	}

	if cursor == nil {
		return nil
	}
	return cursor
}

// view component
func (m model) View() tea.View {
	var output []string
	inputLine := -1

	switch m.state {
	case stateName:
		output, inputLine = m.viewName()
	}

	view := tea.NewView(strings.Join(output, "\n"))
	view.Cursor = m.cursor(inputLine)
	return view
}
