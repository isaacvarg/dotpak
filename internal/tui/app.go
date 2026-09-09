// Package tui handles the bubbletea interface
package tui

import (
	"charm.land/lipgloss/v2"
)

// new to bubbletea so following
// https://github.com/charmbracelet/bubbletea/tree/main/tutorials/basics
// https://github.com/charmbracelet/bubbletea/tree/main/examples/list-simple

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

	nameInput string

}

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingLeft(4).
	Width(22)

func TestTUI() {
	lipgloss.Println(style.Render("hey there"))
}

func initialModel() model {
	return model{}
}



// init component
// update component
// view component
