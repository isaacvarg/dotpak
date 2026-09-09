package tui

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

func (m model) summary() string {
	// styles the response answers
	value := func(s string, c color.Color) string {
		return lipgloss.NewStyle().Bold(true).Foreground(c).Render(s)
	}

	// renders the data set as the key value (or question / answer) with style
	row := func(key, value string) string {
		return labelStyle.Render(key) + value
	}

	body := strings.Join([]string{
		gradientStyle(m.name + " was added to the manifest!"),
		"",
		row("name", value(m.name, lipgloss.Color("#cad3f5"))),
		row("source", value(m.source, sourceColors[m.source])),
		row("group", value(m.group, lipgloss.Color("#cad3f5"))),
		row("command", value(m.command, lipgloss.Color("#7dc4e4"))),
	}, "\n")

	return boxStyle.Render(body) + "\n"
}
