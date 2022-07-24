package grid

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/slices"
	"strings"
)

const (
	gridMode = iota
	searchMode
)

var (
	noElementsError                 = fmt.Errorf("grid cannot be empty")
	elementsDoNotFitBoundariesError = fmt.Errorf("number of cells does not fit grid size")
	bottomBarHeight                 = 5
)

type GridSettings struct {
	Rows    int
	Columns int
}

type Model struct {
	cells                    []Cell
	grid                     [][]*Cell
	selectedX, selectedY     int
	viewport                 viewport.Model
	search                   textinput.Model
	help                     help.Model
	state                    int
	maxHeight                int
	highPerformanceRendering bool
}

func getDirectionFromKey(directionKey string) (direction string) {
	switch directionKey {
	case "h":
		direction = "left"
	case "j":
		direction = "down"
	case "k":
		direction = "up"
	case "l":
		direction = "right"
	default:
		direction = directionKey
	}

	return direction
}

func (m *Model) SelectCell(directionKey string) {
	(*m.grid[m.selectedY][m.selectedX]).SetSelected(false)
	direction := getDirectionFromKey(directionKey)
	m.selectedX, m.selectedY = m.getNextNonHiddenCell(direction)
	(*m.grid[m.selectedY][m.selectedX]).SetSelected(true)
}

func (m *Model) setSelectedCell(idx int) {
	y := idx / len(m.grid[0])
	x := idx - y*len(m.grid[0])

	(*m.grid[m.selectedY][m.selectedX]).SetSelected(false)
	m.selectedX, m.selectedY = x, y
	(*m.grid[m.selectedY][m.selectedX]).SetSelected(true)
}

func (m *Model) getNextNonHiddenCell(direction string) (int, int) {
	planeValue := m.selectedX
	planeValueModifier := func() { planeValue++ }

	if direction == "up" || direction == "left" {
		planeValueModifier = func() { planeValue-- }
	}

	if direction == "up" || direction == "down" {
		planeValue = m.selectedY
	}

	rows := len(m.grid)
	cols := len(m.grid[0])

	for {
		planeValueModifier()

		var cell *Cell
		switch direction {
		case "up":
			if planeValue < 0 {
				return m.selectedX, m.selectedY
			}
			cell = m.grid[planeValue][m.selectedX]
			if cell != nil && !(*cell).IsPaddingCell() {
				return m.selectedX, planeValue
			}
		case "down":
			if planeValue >= rows {
				return m.selectedX, m.selectedY
			}
			cell = m.grid[planeValue][m.selectedX]
			if cell != nil && !(*cell).IsPaddingCell() {
				return m.selectedX, planeValue
			}
		case "right":
			if planeValue >= cols {
				return m.selectedX, m.selectedY
			}
			cell = m.grid[m.selectedY][planeValue]
			if cell != nil && !(*cell).IsPaddingCell() {
				return planeValue, m.selectedY
			}
		case "left":
			if planeValue < 0 {
				return m.selectedX, m.selectedY
			}
			cell = m.grid[m.selectedY][planeValue]
			if cell != nil && !(*cell).IsPaddingCell() {
				return planeValue, m.selectedY
			}
		}
	}
}

func (m *Model) SetGrid(settings GridSettings) error {
	if m.cells == nil {
		return fmt.Errorf("grid cannot be nil")
	}

	err := allocateColumnsAndRows(len(m.cells), &settings)
	if err != nil {
		return err
	}

	m.grid = fillGrid(m.cells, settings)

	return nil
}

func fillGrid(elements []Cell, settings GridSettings) (grid [][]*Cell) {

	for i := 0; i < settings.Rows; i++ {
		if i == 0 {
			elements[0].SetSelected(true)
		}

		var row []*Cell
		for j := 0; j < settings.Columns; j++ {
			row = append(row, &elements[i*settings.Columns+j])
		}
		grid = append(grid, row)
	}

	return grid
}

func allocateColumnsAndRows(count int, settings *GridSettings) error {
	if count == 0 {
		return noElementsError
	}

	if settings.Columns == 0 && settings.Rows == 0 {
		settings.Columns = count
		settings.Rows = 1
	}

	if settings.Rows == 0 {
		settings.Rows = count / settings.Columns
	} else if settings.Columns == 0 {
		settings.Columns = count / settings.Rows
	}

	if settings.Rows*settings.Columns != count {
		return elementsDoNotFitBoundariesError
	}

	return nil
}

func (m *Model) GetSelectedElement() *Cell {
	if m.grid == nil {
		return nil
	}

	if m.checkSelectedInBounds() {
		return m.grid[m.selectedY][m.selectedX]
	}
	return nil
}

func (m *Model) checkSelectedInBounds() bool {
	if m.selectedY >= 0 && m.selectedY < len(m.grid) && m.selectedX >= 0 && m.selectedX < len(m.grid[m.selectedY]) {
		return true
	}
	return false
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if m.state == gridMode {
			switch key {
			case "up", "down", "left", "right", "h", "j", "k", "l":
				m.SelectCell(key)
			case "/":
				m.state = searchMode
				m.search.Focus()
			case "q":
				return m, tea.Quit
			}
		} else {
			if key == "esc" || key == "enter" {
				m.state = gridMode
				m.search.Reset()
			} else {
				m.search, cmd = m.search.Update(msg)
				if value := m.search.Value(); value != "" {
					m.SearchCells(value)
				}
				cmds = append(cmds, cmd)
			}
		}
		m.viewport.SetContent(m.viewGrid())
	case tea.WindowSizeMsg:
		maxHeight := len(m.grid) * (m.cells[0].GetUnselectedStyle().GetVerticalFrameSize() + m.cells[0].GetUnselectedStyle().GetHeight() + 1)
		fmt.Print()
		if msg.Height > maxHeight {
			m.viewport = viewport.New(msg.Width, maxHeight)
		} else {
			m.viewport = viewport.New(msg.Width, msg.Height-bottomBarHeight)
		}
		m.viewport.HighPerformanceRendering = m.highPerformanceRendering
		m.viewport.SetContent(m.viewGrid())
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) viewGrid() string {
	var text string
	for _, row := range m.grid {
		var rowString string
		for _, elmt := range row {
			rowString = lipgloss.JoinHorizontal(0, rowString, (*elmt).GetView())
		}
		text = lipgloss.JoinVertical(0, text, rowString)
	}

	return text
}

func (m Model) View() string {
	text := m.viewport.View()
	if m.state == searchMode {
		text = lipgloss.JoinVertical(0, text, m.search.View())
	}

	return text
}

func (m *Model) SearchCells(searchText string) {
	searchText = strings.ToLower(searchText)
	idx := slices.IndexFunc(m.cells, func(c Cell) bool {
		cellString := strings.ToLower(c.GetSearchString())
		if len(cellString) < len(searchText) {
			return false
		}
		return cellString[:len(searchText)] == searchText
	})
	if idx != -1 {
		m.setSelectedCell(idx)
	}
}

func CreateModel(cells []Cell, gridSettings GridSettings) (Model, error) {
	search := textinput.New()
	search.Prompt = "/"

	model := Model{
		cells:                    cells,
		help:                     help.New(),
		search:                   search,
		highPerformanceRendering: false,
	}
	err := model.SetGrid(gridSettings)

	return model, err
}
