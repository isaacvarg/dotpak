package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// lipgloss stylin

var (
	questionStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("250"))
	errorStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Red)
	controlsStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("255"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#eed49f"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#a5adcb"))
	boxStyle      = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#f5a97f")).
			Padding(1, 3)
)

// install type colors
// i for some reason called it source in the tui package out of nowhere...
var sourceColors = map[string]color.Color{
	"pacman":  lipgloss.Color("#f0c6c6"),
	"aur":     lipgloss.Color("#f5bde6"),
	"mise":    lipgloss.Color("#ee99a0"),
	"omarchy": lipgloss.Color("#f5a97f"),
	"flatpak": lipgloss.Color("#a6da95"),
}
