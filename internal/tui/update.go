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
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			return m.confirm()
		case "j", "down":
			return m.move(+1), nil
		case "k", "up":
			return m.move(-1), nil
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

	case stateSource:
		return m.selectSource(), nil

	case stateGroup:
		m = m.selectGroup()
		return m, nil

	case stateCommand:
		return m.setCommand(m.commandInput.Value()), nil
	}

	return m, nil
}

// handles moving selection
func (m model) move(delta int) model {
	switch m.state {
	case stateSource:
		m.sourceIdx = wrap(m.sourceIdx+delta, len(m.sources))
	case stateGroup:
		m.groupIdx = wrap(m.groupIdx+delta, len(m.groups))
	}

	return m
}

// so we can go from bottom to top and other way
// there is a shorter way to write this but i am keeping
// it as this for readability, also % is different than in Typescript
func wrap(i, n int) int {
	remainder := i % n

	if remainder < 0 {
		remainder = remainder + n
	}

	return remainder
}

func (m model) cursor(line int) *tea.Cursor {
	if line < 0 {
		return nil
	}

	var cursor *tea.Cursor
	switch m.state {
	case stateName:
		cursor = m.nameInput.Cursor()

	case stateCommand:
		cursor = m.commandInput.Cursor()
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
	case stateSource:
		output = m.viewSource()
	case stateGroup:
		output = m.viewGroup()
	case stateCommand:
		output, inputLine = m.viewCommand()
	}

	view := tea.NewView(strings.Join(output, "\n"))
	view.Cursor = m.cursor(inputLine)
	return view
}
