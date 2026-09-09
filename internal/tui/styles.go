package tui

import "charm.land/lipgloss/v2"

// lipgloss stylin

// from lipgloss docs tutorial
var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	PaddingLeft(4).
	Width(22)

var (
	questionStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("250"))
	errorStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Red)
	controlsStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255"))
)
