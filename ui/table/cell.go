package table

import "github.com/charmbracelet/lipgloss"

// Interface
type Cell[T any] interface {
	GetView() string
	SetStyle(style lipgloss.Style)
	GetSearchString()
	GetDataOnSelect() interface{}
}

// Implementing interface
type ElementCell struct {
	Data         string
	Style        lipgloss.Style
	SearchString string
}

func (c *ElementCell) GetView() string {
	var text string
	// Put formatting/styling here
	return text
}

func (c *ElementCell) SetStyle(style lipgloss.Style) {
	c.Style = style
}

func (c *ElementCell) GetSearchString() string {
	return c.SearchString
}
