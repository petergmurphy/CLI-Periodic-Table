package element

import (
	tea "github.com/charmbracelet/bubbletea"
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

type Model struct {
	Data       Data
	IsSelected bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return getView(m.Data.AtomicNumber, m.Data.Symbol, m.IsSelected, m.Data.Type)
}
