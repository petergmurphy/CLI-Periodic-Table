package table

import "github.com/charmbracelet/lipgloss"

// Interface
type Cell interface {
	GetView() string
	SetStyle(selectedStyle lipgloss.Style, unSelectedStyle lipgloss.Style)
	SetSelected(isSelected bool)
	GetUnselectedStyle() lipgloss.Style
	GetSearchString() string
	GetData() interface{}
	IsPaddingCell() bool
}

// Implementing interface
