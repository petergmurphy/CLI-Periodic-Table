package table

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"periodic-table/ui/grid"
	"periodic-table/ui/periodic_table/element"
	"periodic-table/ui/periodic_table/keys"
)

const (
	gridMode = iota
	searchMode
)

const bottomBarHeight = 1

type model struct {
	grid           grid.Model
	viewport       viewport.Model
	state          int
	help           help.Model
	keys           keys.KeyMap
	search         textinput.Model
	terminalHeight int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.grid, cmd = m.grid.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		maxHeight := lipgloss.Height(m.grid.View())
		m.terminalHeight = msg.Height
		if msg.Height > maxHeight {
			m.viewport = viewport.New(msg.Width, maxHeight)
		} else {
			m.viewport = viewport.New(msg.Width, m.terminalHeight-1)
		}
	case tea.KeyMsg:
		key := msg.String()
		if m.state == gridMode {
			switch key {
			case "/":
				m.state = searchMode
				m.search.Focus()
			case "q":
				return m, tea.Quit
			}
		} else if m.state == searchMode {
			if key == "esc" || key == "enter" {
				m.state = gridMode
				m.search.Reset()
			} else {
				m.search, cmd = m.search.Update(msg)
				if value := m.search.Value(); value != "" {
					m.grid.SearchCells(value)
				}
				return m, cmd
			}
		}
	}
	m.viewport.SetContent(m.grid.View())
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	text := lipgloss.JoinHorizontal(0, m.viewport.View(), m.getElementInfoView())
	relativeBottomBarPos := m.terminalHeight - lipgloss.Height(m.grid.View())
	if m.state == searchMode {
		searchBar := lipgloss.PlaceVertical(relativeBottomBarPos, lipgloss.Bottom, m.search.View())
		text = lipgloss.JoinVertical(0, text, searchBar)
	} else if m.state == gridMode {
		helpBar := lipgloss.PlaceVertical(relativeBottomBarPos, lipgloss.Bottom, m.help.View(m.keys))
		text = lipgloss.JoinVertical(0, text, helpBar)
	}

	return text
}

func (m model) getElementInfoView() string {
	switch elementData := (*m.grid.GetActiveCell()).GetData().(type) {
	case element.Data:
		return element.ElementInfoView(elementData)
	}
	return ""
}

func CreateModel(cells []grid.Cell) (tea.Model, error) {
	search := textinput.New()
	search.Prompt = "Search: "

	g, err := grid.CreateModel(cells, grid.GridSettings{Rows: 10, Columns: 18})
	if err != nil {
		return nil, err
	}

	model := model{
		help:   help.New(),
		search: search,
		keys:   keys.CreateKeys(),
		grid:   g,
	}
	return model, nil
}
