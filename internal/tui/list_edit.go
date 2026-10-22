package tui

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type editStep int

const (
	editName editStep = iota
	editInstallType
)

type listEdit struct {
	step      editStep
	name      string
	err       string
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
	if key.String() == "ctrl+c" {
		return m, tea.Quit
	}
	if key.String() == "esc" {
		m.edit = nil
		return m, nil
	}
	switch m.edit.step {
	case editName:
		if key.String() == "enter" {
			return m.confirmEditName(), nil
		}
		m.edit.nameInput, _ = m.edit.nameInput.Update(key)
		return m, nil
	default:
		return m, nil
	}
}

func (m listModel) confirmEditName() listModel {
	name := strings.TrimSpace(m.edit.nameInput.Value())
	if name == "" {
		m.edit.err = "a name is required"
		return m
	}
	m.edit.err = ""
	m.edit.name = name
	m.edit.step = editInstallType
	return m
}

func (m listModel) renderEdit() string {
	switch m.edit.step {
	case editInstallType:
		return questionStyle.Render("What is the install type?")
	default:
		return m.renderEditName()
	}
}

func (m listModel) renderEditName() string {
	var b strings.Builder
	b.WriteString(questionStyle.Render("What is the name?"))
	b.WriteString("\n\n")
	b.WriteString(m.edit.nameInput.View())
	if m.edit.err != "" {
		b.WriteString("\n\n")
		b.WriteString(errorStyle.Render(m.edit.err))
	}
	b.WriteString("\n\n")
	b.WriteString(controlsStyle.Render("enter confirm • esc back"))
	return b.String()
}
