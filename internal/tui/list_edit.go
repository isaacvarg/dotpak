package tui

import (
	"image/color"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/isaacvarg/dotpak/internal/groups"
	"github.com/isaacvarg/dotpak/internal/manifest"
)

var editInstallTypes = []string{"pacman", "aur", "mise", "omarchy", "flatpak"}

type editStep int

const (
	editName editStep = iota
	editInstallType
	editGroup
	editCommand
	editSaveFailed
)

type listEdit struct {
	step           editStep
	originalName   string
	name           string
	installType    string
	group          string
	command        string
	err            string
	nameInput      textinput.Model
	commandInput   textinput.Model
	installTypeIdx int
	groups         []string
	groupIdx       int
}

func (m listModel) beginEdit() listModel {
	entry := m.entries[m.cursor]
	lg, err := groups.Load()
	if err != nil {
		m.edit = &listEdit{step: editSaveFailed, err: err.Error()}
		return m
	}
	name := textinput.New()
	name.Prompt = "› "
	name.SetWidth(40)
	name.Focus()
	name.SetValue(entry.Name)
	command := textinput.New()
	command.Prompt = "› "
	command.SetWidth(40)
	command.Focus()
	command.SetValue(entry.InstallCommand)
	choices := editGroupChoices(lg.Names(), entry.Group)
	m.edit = &listEdit{
		step:           editName,
		originalName:   entry.Name,
		nameInput:      name,
		commandInput:   command,
		installTypeIdx: installTypeIndex(entry.InstallType),
		groups:         choices,
		groupIdx:       groupIndex(choices, entry.Group),
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
		case "enter":
			m.edit.installType = editInstallTypes[m.edit.installTypeIdx]
			m.edit.step = editGroup
		}
		return m, nil
	case editGroup:
		switch key.String() {
		case "j", "down":
			m.edit.groupIdx = wrap(m.edit.groupIdx+1, len(m.edit.groups))
		case "k", "up":
			m.edit.groupIdx = wrap(m.edit.groupIdx-1, len(m.edit.groups))
		case "enter":
			m.edit.group = m.edit.groups[m.edit.groupIdx]
			m.edit.step = editCommand
		}
		return m, nil
	case editCommand:
		if key.String() == "enter" {
			return m.confirmEditCommand(), nil
		}
		m.edit.commandInput, _ = m.edit.commandInput.Update(key)
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

func editGroupChoices(existing []string, current string) []string {
	names := append([]string(nil), existing...)
	for _, name := range names {
		if name == current {
			return names
		}
	}
	return append(names, current)
}

func groupIndex(names []string, current string) int {
	for i, name := range names {
		if name == current {
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

func (m listModel) confirmEditCommand() listModel {
	command := m.edit.commandInput.Value()
	if command == "" {
		command = m.edit.name
	}
	m.edit.command = command
	return m.saveEdit()
}

func (m listModel) saveEdit() listModel {
	mf, err := manifest.Load()
	if err != nil {
		return m.failEdit(err)
	}
	idx := -1
	for i, existing := range mf.Manifest {
		if strings.EqualFold(existing.Name, m.edit.originalName) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return m.failEdit(manifest.ErrNotFound)
	}
	mf.Manifest[idx] = manifest.ManifestEntry{
		Name:           m.edit.name,
		Group:          m.edit.group,
		InstallType:    manifest.InstallType(m.edit.installType),
		InstallCommand: m.edit.command,
	}
	if err := mf.Save(); err != nil {
		return m.failEdit(err)
	}
	return m.finishEdit(mf)
}

func (m listModel) failEdit(err error) listModel {
	m.edit.step = editSaveFailed
	m.edit.err = err.Error()
	return m
}

func (m listModel) finishEdit(mf *manifest.Manifest) listModel {
	oldCursor := m.cursor
	name := m.edit.name
	entries := filterEntries(mf.Manifest, m.filterGroup, m.filterInstallType)
	next := newListModel(entries)
	next.filterGroup = m.filterGroup
	next.filterInstallType = m.filterInstallType
	found := -1
	for i, entry := range entries {
		if entry.Name == name {
			found = i
			break
		}
	}
	switch {
	case found >= 0:
		next.cursor = found
	case len(entries) == 0:
		next.cursor = 0
	case oldCursor < len(entries):
		next.cursor = oldCursor
	default:
		next.cursor = len(entries) - 1
	}
	return next
}

func (m listModel) renderEdit() string {
	switch m.edit.step {
	case editInstallType:
		return m.renderEditInstallType()
	case editGroup:
		return m.renderEditGroup()
	case editCommand:
		return m.renderEditCommand()
	case editSaveFailed:
		return errorStyle.Render(m.edit.err)
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

func (m listModel) renderEditGroup() string {
	var b strings.Builder
	b.WriteString(questionStyle.Render("Which group?"))
	b.WriteString("\n\n")
	b.WriteString(editChoices(m.edit.groups, m.edit.groupIdx, nil))
	b.WriteString(controlsStyle.Render("j/k move • enter select"))
	return b.String()
}

func (m listModel) renderEditCommand() string {
	var b strings.Builder
	b.WriteString(questionStyle.Render("What is the install command?"))
	b.WriteString("\n\n")
	b.WriteString(m.edit.commandInput.View())
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("leave blank to use the name: " + m.edit.name))
	b.WriteString("\n\n")
	b.WriteString(controlsStyle.Render("enter confirm • esc back"))
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
