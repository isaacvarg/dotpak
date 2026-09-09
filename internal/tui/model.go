package tui

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/isaacvarg/dotpak/internal/groups"
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
	nameInput    textinput.Model
	commandInput textinput.Model

	// choices and currently selected index
	sources   []string
	sourceIdx int
	groups    []string
	groupIdx  int

	// answers
	name    string
	source  string
	group   string
	command string

	err string
}

// initial stuff
func initialModel() (model, error) {
	lg, err := groups.Load()
	if err != nil {
		return model{}, err
	}

	if len(lg.Groups) == 0 {
		return model{}, fmt.Errorf("whoops, please add some groups first. see dotpak group --help")
	}

	groupNames := make([]string, 0, len(lg.Groups))
	for _, g := range lg.Groups {
		groupNames = append(groupNames, g.Name)
	}

	m := model{
		state:        stateName,
		nameInput:    newNameInput(),
		sources:      []string{"pacman", "aur", "mise", "omarchy", "flatpak"},
		groups:       groupNames,
		commandInput: newCommandInput(),
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) selection(options []string, idx int, colors map[string]color.Color) string {
	var builder strings.Builder

	for i, option := range options {
		if i == idx {
			style := selectedStyle

			if c, ok := colors[option]; ok {
				style = style.Foreground(c)
			}
			builder.WriteString("- ")
			builder.WriteString(style.Render(option))
			builder.WriteString("\n")
			continue
		}

		// not selected
		builder.WriteString("  ")
		builder.WriteString(option)
		builder.WriteString("\n")

	}

	return builder.String()
}
