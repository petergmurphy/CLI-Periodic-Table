package ui

import (
	"periodic-table/src/elements"
	"periodic-table/ui/element"
	"periodic-table/ui/grid"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	tableView = iota
	elementView
)

type Model struct {
	table grid.Model
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

func elementInfoView(elem element.Data) string {
	heading := lipgloss.Place(10, 5, 1, 1, elem.Symbol)
	heading = lipgloss.JoinVertical(0, heading, lipgloss.Place(20, 5, 0.5, 0, elem.Element))

	body := elem.GetDataAsString()
	body = lipgloss.Place(10, 10, 0, 0, body)

	text := lipgloss.JoinVertical(0, heading, body)
	style = style.BorderForeground(element.TypeColors[elem.Type])

	return style.Render(text)
}

func (m Model) View() string {
	text := m.table.View()

	switch elementData := (*m.table.GetSelectedElement()).GetData().(type) {
	case element.Data:
		elementText := elementInfoView(elementData)
		text = lipgloss.JoinHorizontal(0, text, elementText)
	}

	return text
}

func CreateModel() (tea.Model, error) {
	elements := elements.ReadElements()

	table, err := grid.CreateModel(elements, grid.GridSettings{Rows: 10, Columns: 18})
	return Model{table: table}, err
}
