package tui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// model
type state int

const (
	stateName = iota
	stateSource
	stateGroup
	stateCommand
	stateDone
)

type model struct {
	state state

	// bubble elements
	nameInput textinput.Model

	// answers
	name string
	err  string
}

// initial stuff
func initialModel() model {
	m := model{
		state:     stateName,
		nameInput: newNameInput(),
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}
