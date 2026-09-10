package tui

import (
	"fmt"
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/isaacvarg/dotpak/internal/manifest"
)

type savedMsg struct{ err error }

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

func (m model) failure() string {
	body := strings.Join([]string{
		errorStyle.Render("Something went wrong!"),
		"",
		errorStyle.Render(m.err),
	}, "\n")

	return boxStyle.Render(body) + "\n"
}

func (m model) save() tea.Cmd {
	return func() tea.Msg {
		mf, err := manifest.Load()
		if err != nil {
			return savedMsg{fmt.Errorf("something went wrong reading manifest: %w", err)}
		}

		if err := mf.Add(
			m.name,
			manifest.InstallType(m.source),
			m.group,
			m.command,
		); err != nil {
			return savedMsg{fmt.Errorf("something went wrong creating %q: %w", m.name, err)}
		}

		if err := mf.Save(); err != nil {
			return savedMsg{fmt.Errorf("someting went wrong saving: %w", err)}
		}

		return savedMsg{nil}
	}
}
