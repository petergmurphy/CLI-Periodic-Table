package element

import (
	"fmt"
	"periodic-table/ui/grid"

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

type Element struct {
	data            Data
	selectedStyle   lipgloss.Style
	unSelectedStyle lipgloss.Style
	searchString    string
	isSelected      bool
	isPaddingCell   bool
}

func (c *Element) GetSearchString() string {
	return c.searchString
}

func (c *Element) GetData() interface{} {
	return c.data
}

func (c *Element) GetView() string {
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

func (c *Element) GetUnselectedStyle() lipgloss.Style {
	return c.unSelectedStyle
}

func (c *Element) GetSelectedStyle() lipgloss.Style {
	return c.selectedStyle
}

func (c *Element) SetStyle(selectedStyle lipgloss.Style, unSelectedStyle lipgloss.Style) {
	c.selectedStyle = selectedStyle
	c.unSelectedStyle = unSelectedStyle
}

func (c *Element) SetSelected(isSelected bool) {
	c.isSelected = isSelected
}

func styleText(atomicNumber string, symbol string) string {
	text := lipgloss.Place(width, height, 1, 1, symbol)
	text = lipgloss.JoinVertical(0, lipgloss.Place(0, 0, 0, 0, atomicNumber), text)
	return text
}

func (c *Element) IsPaddingCell() bool {
	return c.isPaddingCell
}

func ElementInfoView(elmt Data) string {
	heading := lipgloss.Place(11, 5, 1, 1, elmt.Symbol)
	heading = lipgloss.JoinVertical(0, heading, lipgloss.Place(20, 5, 0.5, 0, elmt.Element))

	body := elmt.GetDataAsString()
	body = lipgloss.Place(11, 10, 0, 0, body)

	text := lipgloss.JoinVertical(0, heading, body)
	style = style.BorderForeground(TypeColors[elmt.Type]).Width(22)

	return style.Render(text)
}

func CreateElement(data Data, isPaddingCell bool) grid.Cell {
	unSelectedStyle := style.Copy().BorderForeground(TypeColors[data.Type])
	selectedStyle := unSelectedStyle.Copy().Background(TypeColors[data.Type])

	cell := &Element{
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
