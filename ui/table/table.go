package table

import (
	"encoding/csv"
	"github.com/charmbracelet/bubbles/help"
	"log"
	"os"
	"periodic-table/ui/element"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	emptyEntryRange = map[int]int{
		1:   17,
		20:  30,
		38:  48,
		126: 148,
		162: 166,
	}
	emptyElement = element.Model{}
	helpHeight   = 5
)

type Model struct {
	elements             []element.Model
	grid                 [][]element.Model
	selectedX, selectedY int
	viewport             viewport.Model
	help                 help.Model
}

func (m *Model) SelectElement(move string) {
	m.grid[m.selectedY][m.selectedX].IsSelected = false
	m.selectedX, m.selectedY = m.getNextNonHiddenElement(move)
	m.grid[m.selectedY][m.selectedX].IsSelected = true
}

func (m *Model) getNextNonHiddenElement(direction string) (int, int) {
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

		switch direction {
		case "up":
			if coordNum < 0 {
				return m.selectedX, m.selectedY
			}
			if m.grid[coordNum][m.selectedX].Data.Type != "" {
				return m.selectedX, coordNum
			}
		case "down":
			if coordNum >= rows {
				return m.selectedX, m.selectedY
			}
			if m.grid[coordNum][m.selectedX].Data.Type != "" {
				return m.selectedX, coordNum
			}
		case "right":
			if coordNum >= cols {
				return m.selectedX, m.selectedY
			}
			if m.grid[m.selectedY][coordNum].Data.Type != "" {
				return coordNum, m.selectedY
			}
		case "left":
			if coordNum < 0 {
				return m.selectedX, m.selectedY
			}
			if m.grid[m.selectedY][coordNum].Data.Type != "" {
				return coordNum, m.selectedY
			}
		}
	}
}

func readElements() []element.Model {
	// open file
	f, err := os.Open("data/elements.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// convert records to array of structs
	return createElements(data)
}

func createElements(data [][]string) []element.Model {
	var elements []element.Model
	var count int
	var lastGroups []element.Model
	for i, line := range data {
		if skipLines := emptyEntryRange[count]; skipLines != 0 {
			for j := count; j < skipLines; j++ {
				elements = append(elements, emptyElement)
				count++
			}
		}

		if i > 0 {
			data := createElementData(line)

			elements = append(elements, element.Model{
				Data:       data,
				IsSelected: i == 1,
			})
			count++
		}
	}

	elements = append(elements, lastGroups...)

	return elements
}

func createElementData(line []string) element.Data {
	return element.Data{
		AtomicNumber:      line[0],
		Element:           line[1],
		Symbol:            line[2],
		AtomicMass:        line[3],
		NumberOfNeutrons:  line[4],
		NumberOfProtons:   line[5],
		NumberOfElectrons: line[6],
		Period:            line[7],
		Group:             line[8],
		Phase:             line[9],
		Radioactive:       line[10],
		Natural:           line[11],
		Metal:             line[12],
		Nonmetal:          line[13],
		Metalloid:         line[14],
		Type:              line[15],
		AtomicRadius:      line[16],
		Electronegativity: line[17],
		FirstIonization:   line[18],
		Density:           line[19],
		MeltingPoint:      line[20],
		BoilingPoint:      line[21],
		NumberOfIsotopes:  line[22],
		Discoverer:        line[23],
		Year:              line[24],
		SpecificHeat:      line[25],
		NumberOfShells:    line[26],
		NumberOfValence:   line[27],
	}
}

func (m *Model) GetSelectedElement() element.Model {
	if m.grid == nil {
		return emptyElement
	}

	if m.checkSelectedInBounds() {
		return m.grid[m.selectedY][m.selectedX]
	}
	return emptyElement
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
			m.SelectElement(key)
		}
		m.viewport.SetContent(m.viewGrid())
	case tea.WindowSizeMsg:
		m.viewport = viewport.New(msg.Width-helpHeight, msg.Height)
		m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
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
			rowString = lipgloss.JoinHorizontal(0, rowString, elmt.View())
		}
		text = lipgloss.JoinVertical(0, text, rowString)
	}

	return text
}

func (m Model) View() string {
	return m.viewport.View()
}

func CreateModel() (Model, error) {
	elements := readElements()
	model := Model{
		elements: elements,
		help:     help.New(),
	}
	err := model.SetGrid(elements, GridSettings{Rows: 10, Columns: 18})
	return model, err
}
