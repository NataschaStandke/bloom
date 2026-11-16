package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	thoughts, err := loadThoughts()
	if err != nil {
		thoughts = []thought{}
	}

	p := tea.NewProgram(newModel(thoughts), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
