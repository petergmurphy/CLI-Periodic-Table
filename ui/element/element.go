package element

import (
	"fmt"
	"periodic-table/ui/table"

	"github.com/charmbracelet/lipgloss"
)

type Data struct {
	AtomicNumber      string
	Element           string
	Symbol            string
	AtomicMass        string
	NumberOfNeutrons  string
	NumberOfProtons   string
	NumberOfElectrons string
	Period            string
	Group             string
	Phase             string
	Radioactive       string
	Natural           string
	Metal             string
	Nonmetal          string
	Metalloid         string
	Type              string
	AtomicRadius      string
	Electronegativity string
	FirstIonization   string
	Density           string
	MeltingPoint      string
	BoilingPoint      string
	NumberOfIsotopes  string
	Discoverer        string
	Year              string
	SpecificHeat      string
	NumberOfShells    string
	NumberOfValence   string
}

func (d *Data) GetDataAsString() string {
	text := fmt.Sprintf(`Type: %s
Atomic number: %s
Atomic mass: %s
Electrons: %s
Protons: %s
Neutrons: %s
Group: %s
Density: %s
Atomic Radius: %s
Melting Point: %s
Specific Heat: %s
	`, d.Type, d.AtomicNumber, d.AtomicMass, d.NumberOfElectrons, d.NumberOfProtons, d.NumberOfNeutrons, d.Group, d.Density, d.AtomicRadius, d.MeltingPoint, d.SpecificHeat)

	return text

}

type Cell struct {
	data            Data
	selectedStyle   lipgloss.Style
	unSelectedStyle lipgloss.Style
	searchString    string
	isSelected      bool
	isPaddingCell   bool
}

func (c *Cell) GetSearchString() string {
	return c.searchString
}

func (c *Cell) GetData() interface{} {
	return c.data
}

func (c *Cell) GetView() string {
	var text string
	// Make cell text here
	text = styleText(c.data.AtomicNumber, c.data.Symbol)
	// Put formatting/styling here
	if c.isSelected {
		text = c.selectedStyle.Render(text)
	} else {
		text = c.unSelectedStyle.Render(text)
	}

	return text
}

func (c *Cell) SetStyle(selectedStyle lipgloss.Style, unSelectedStyle lipgloss.Style) {
	c.selectedStyle = selectedStyle
	c.unSelectedStyle = unSelectedStyle
}

func (c *Cell) SetSelected(isSelected bool) {
	c.isSelected = isSelected
}

func styleText(atomicNumber string, symbol string) string {
	text := lipgloss.Place(width, height, 1, 1, symbol)
	text = lipgloss.JoinVertical(0, lipgloss.Place(0, 0, 0, 0, atomicNumber), text)
	return text
}

func (c *Cell) IsPaddingCell() bool {
	return c.isPaddingCell
}

func CreateElement(data Data, isPaddingCell bool) table.Cell {
	unSelectedStyle := style.Copy().BorderForeground(TypeColors[data.Type])
	selectedStyle := unSelectedStyle.Copy().Background(TypeColors[data.Type])

	cell := &Cell{
		data:            data,
		selectedStyle:   selectedStyle,
		unSelectedStyle: unSelectedStyle,
		searchString:    data.Element,
		isSelected:      false,
		isPaddingCell:   isPaddingCell,
	}

	if isPaddingCell {
		cell.unSelectedStyle = empty
	}

	return cell
}
