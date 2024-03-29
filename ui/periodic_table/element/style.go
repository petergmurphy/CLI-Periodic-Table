package element

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	width  = 6
	height = 1
)

var (
	style      = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).Width(width).Height(height)
	empty      = lipgloss.NewStyle().Width(width + 2).Height(height + 1)
	TypeColors = map[string]lipgloss.Color{
		"Nonmetal":             lipgloss.Color("#cf53a4"),
		"Noble Gas":            lipgloss.Color("#697a90"),
		"Alkali Metal":         lipgloss.Color("#a86d69"),
		"Alkaline Earth Metal": lipgloss.Color("#7b524c"),
		"Metalloid":            lipgloss.Color("#e4a54d"),
		"Transition Metal":     lipgloss.Color("#c778a5"),
		"Transactinide":        lipgloss.Color("#4c36a7"),
		"Actinide":             lipgloss.Color("#1aa29d"),
		"Lanthanide":           lipgloss.Color("#4b93c6"),
		"Metal":                lipgloss.Color("#136F63"),
		"Halogen":              lipgloss.Color("#052cbf"),
		"artificial":           lipgloss.Color("#f9c857"),
	}
)

func getView(atomicNumber string, symbol string, isSelected bool, elementType string) string {
	if symbol == "" {
		return empty.Render("")
	}

	text := styleText(atomicNumber, symbol)
	if isSelected {
		style = style.Background(TypeColors[elementType])
	} else {
		style = style.UnsetBackground()
	}

	return style.Render(text)
}
