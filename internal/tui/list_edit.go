package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type editStep int

const (
	editName editStep = iota
)

type listEdit struct {
	step      editStep
	nameInput textinput.Model
}

func (m listModel) beginEdit() listModel {
	entry := m.entries[m.cursor]
	name := textinput.New()
	name.Prompt = "› "
	name.SetWidth(40)
	name.Focus()
	name.SetValue(entry.Name)
	m.edit = &listEdit{
		step:      editName,
		nameInput: name,
	}
	return m
}

func (m listModel) updateEdit(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch key.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.edit = nil
		return m, nil
	default:
		m.edit.nameInput, _ = m.edit.nameInput.Update(key)
		return m, nil
	}
}

func (m listModel) renderEdit() string {
	var b strings.Builder
	b.WriteString(questionStyle.Render("What is the name?"))
	b.WriteString("\n\n")
	b.WriteString(m.edit.nameInput.View())
	b.WriteString("\n\n")
	b.WriteString(controlsStyle.Render("enter confirm • esc back"))
	return b.String()
}
