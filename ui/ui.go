package ui

import (
	"periodic-table/src/elements"
	"periodic-table/ui/periodic_table/table"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	tableView = iota
	elementView
)

type Model struct {
	table tea.Model
	state int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.table.View()
}

func CreateModel() (tea.Model, error) {
	elmts := elements.ReadElements()

	t, err := table.CreateModel(elmts)
	return Model{table: t}, err
}
