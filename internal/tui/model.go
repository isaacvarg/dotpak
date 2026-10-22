package tui

import (
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
	stateSaveFailed
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

	m := model{
		state:        stateName,
		nameInput:    newNameInput(),
		sources:      []string{"pacman", "aur", "mise", "omarchy", "flatpak"},
		groups:       lg.Names(),
		commandInput: newCommandInput(),
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) selection(options []string, idx int, colors map[string]color.Color) string {
	return renderChoices(options, idx, colors)
}

func renderChoices(options []string, idx int, colors map[string]color.Color) string {
	var b strings.Builder
	for i, option := range options {
		if i == idx {
			style := selectedStyle
			if c, ok := colors[option]; ok {
				style = style.Foreground(c)
			}
			b.WriteString("- ")
			b.WriteString(style.Render(option))
		} else {
			b.WriteString("  ")
			b.WriteString(option)
		}
		b.WriteByte('\n')
	}
	return b.String()
}
