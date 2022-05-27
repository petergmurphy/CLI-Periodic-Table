package ui

import (
	"fmt"
	"periodic-table/src/elements"
	"periodic-table/ui/element"
	"periodic-table/ui/table"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	tableView = iota
	elementView
)

type Model struct {
	table table.Model
	state int
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch m.state {
	case tableView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				m.state = elementView
			}
		}
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func elementInfoView(elem element.Data) string {
	text := fmt.Sprintf(`Type: %s
Atomic number: %s
Atomic mass: %s
Group: %s
`, elem.Type, elem.AtomicNumber, elem.AtomicMass, elem.Group)

	heading := lipgloss.Place(10, 5, 0.5, 1, elem.Symbol)
	heading = lipgloss.JoinVertical(0, heading, lipgloss.Place(10, 5, 0.5, 0, elem.Element))
	text = lipgloss.Place(10, 10, 0.5, 0.5, text)
	text = lipgloss.JoinVertical(0, heading, text)
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

	table, err := table.CreateModel(elements, table.GridSettings{Rows: 10, Columns: 18})
	return Model{table: table}, err
}
