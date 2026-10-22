package tui

import (
	"image/color"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/isaacvarg/dotpak/internal/manifest"
)

var editInstallTypes = []string{"pacman", "aur", "mise", "omarchy", "flatpak"}

type editStep int

const (
	editName editStep = iota
	editInstallType
)

type listEdit struct {
	step           editStep
	originalName   string
	name           string
	err            string
	nameInput      textinput.Model
	installTypeIdx int
}

func (m listModel) beginEdit() listModel {
	entry := m.entries[m.cursor]
	name := textinput.New()
	name.Prompt = "› "
	name.SetWidth(40)
	name.Focus()
	name.SetValue(entry.Name)
	m.edit = &listEdit{
		step:           editName,
		originalName:   entry.Name,
		nameInput:      name,
		installTypeIdx: installTypeIndex(entry.InstallType),
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
	case editInstallType:
		switch key.String() {
		case "j", "down":
			m.edit.installTypeIdx = wrap(m.edit.installTypeIdx+1, len(editInstallTypes))
		case "k", "up":
			m.edit.installTypeIdx = wrap(m.edit.installTypeIdx-1, len(editInstallTypes))
		}
		return m, nil
	default:
		return m, nil
	}
}

func installTypeIndex(current manifest.InstallType) int {
	for i, name := range editInstallTypes {
		if name == string(current) {
			return i
		}
	}
	return 0
}

func (m listModel) confirmEditName() listModel {
	name := strings.TrimSpace(m.edit.nameInput.Value())
	if name == "" {
		m.edit.err = "a name is required"
		return m
	}
	mf, err := manifest.Load()
	if err != nil {
		m.edit.err = err.Error()
		return m
	}
	for _, existing := range mf.Manifest {
		if strings.EqualFold(existing.Name, m.edit.originalName) {
			continue
		}
		if strings.EqualFold(existing.Name, name) {
			m.edit.err = manifest.ErrDuplicateEntry.Error()
			return m
		}
	}
	m.edit.err = ""
	m.edit.name = name
	m.edit.step = editInstallType
	return m
}

func (m listModel) renderEdit() string {
	switch m.edit.step {
	case editInstallType:
		return m.renderEditInstallType()
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

func (m listModel) renderEditInstallType() string {
	var b strings.Builder
	b.WriteString(questionStyle.Render("What is the install type?"))
	b.WriteString("\n\n")
	b.WriteString(editChoices(editInstallTypes, m.edit.installTypeIdx, sourceColors))
	b.WriteString(controlsStyle.Render("j/k move • enter select"))
	return b.String()
}

func editChoices(options []string, idx int, colors map[string]color.Color) string {
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
