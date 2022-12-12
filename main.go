package main

import (
	"fmt"
	"os"
	"periodic-table/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model, err := ui.CreateModel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := tea.NewProgram(model, tea.WithAltScreen()).Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
