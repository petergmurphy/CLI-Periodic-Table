package table

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	noElementsError                 = fmt.Errorf("grid cannot be empty")
	elementsDoNotFitBoundariesError = fmt.Errorf("number of cells does not fit grid size")
	helpHeight                      = 5
)

type GridSettings struct {
	Rows    int
	Columns int
}

type Model struct {
	cells                []Cell
	grid                 [][]*Cell
	selectedX, selectedY int
	viewport             viewport.Model
	help                 help.Model
}

func (m *Model) SelectCell(move string) {
	(*m.grid[m.selectedY][m.selectedX]).SetSelected(false)
	m.selectedX, m.selectedY = m.getNextNonHiddenCell(move)
	(*m.grid[m.selectedY][m.selectedX]).SetSelected(true)
}

func (m *Model) getNextNonHiddenCell(direction string) (int, int) {
	coordNum := m.selectedX
	if direction == "up" || direction == "down" {
		coordNum = m.selectedY
	}

	rows := len(m.grid)
	cols := len(m.grid[0])

	for {
		if direction == "up" || direction == "left" {
			coordNum--
		} else {
			coordNum++
		}

		var cell *Cell
		switch direction {
		case "up":
			if coordNum < 0 {
				return m.selectedX, m.selectedY
			}
			cell = m.grid[coordNum][m.selectedX]
			if cell != nil && !(*cell).IsPaddingCell() {
				return m.selectedX, coordNum
			}
		case "down":
			if coordNum >= rows {
				return m.selectedX, m.selectedY
			}
			cell = m.grid[coordNum][m.selectedX]
			if cell != nil && !(*cell).IsPaddingCell() {
				return m.selectedX, coordNum
			}
		case "right":
			if coordNum >= cols {
				return m.selectedX, m.selectedY
			}
			cell = m.grid[m.selectedY][coordNum]
			if cell != nil && !(*cell).IsPaddingCell() {
				return coordNum, m.selectedY
			}
		case "left":
			if coordNum < 0 {
				return m.selectedX, m.selectedY
			}
			cell = m.grid[m.selectedY][coordNum]
			if cell != nil && !(*cell).IsPaddingCell() {
				return coordNum, m.selectedY
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

	m.grid = fillGrid(&m.cells, settings)

	return nil
}

func fillGrid(elements *[]Cell, settings GridSettings) (grid [][]*Cell) {
	for i := 0; i < settings.Rows; i++ {
		var row []*Cell
		for j := 0; j < settings.Columns; j++ {
			row = append(row, &(*elements)[i*settings.Columns+j])
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
		if key == "q" {
			return m, tea.Quit
		}
		if key == "up" || key == "down" || key == "left" || key == "right" {
			m.SelectCell(key)
		}
		m.viewport.SetContent(m.viewGrid())
	case tea.WindowSizeMsg:
		m.viewport = viewport.New(msg.Width-helpHeight, msg.Height)
		m.viewport.HighPerformanceRendering = false
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
	return m.viewport.View()
}

func CreateModel(cells []Cell, gridSettings GridSettings) (Model, error) {
	model := Model{
		cells: cells,
		help:  help.New(),
	}
	err := model.SetGrid(gridSettings)
	return model, err
}
