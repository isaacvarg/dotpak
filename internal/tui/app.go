// Package tui handles the bubbletea interface
package tui

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

// new to bubbletea so following
// https://github.com/charmbracelet/bubbletea/tree/main/tutorials/basics
// https://github.com/charmbracelet/bubbletea/tree/main/examples/list-simple

func App() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
