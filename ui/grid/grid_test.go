package grid

import (
	"github.com/charmbracelet/lipgloss"
	"testing"
)

type mockCell struct {
	SelectedStyle   lipgloss.Style
	UnSelectedStyle lipgloss.Style
	SearchString    string
	IsSelected      bool
	IsPaddingCell   bool
	View            string
}

func (c *mockCell) GetSearchString() string {
	return c.searchString
}

func (c *mockCell) GetData() interface{} {
	return "no data"
}

func (c *mockCell) GetView() string {
	return c.View
}

func (c *mockCell) GetUnselectedStyle() lipgloss.Style {
	return c.unSelectedStyle
}

func (c *mockCell) SetStyle(selectedStyle lipgloss.Style, unSelectedStyle lipgloss.Style) {
	c.selectedStyle = selectedStyle
	c.unSelectedStyle = unSelectedStyle
}

func (c *mockCell) SetSelected(isSelected bool) {
	c.isSelected = isSelected
}

func styleText(atomicNumber string, symbol string) string {
	text := lipgloss.Place(width, height, 1, 1, symbol)
	text = lipgloss.JoinVertical(0, lipgloss.Place(0, 0, 0, 0, atomicNumber), text)
	return text
}

func TestModel_SetGrid(t *testing.T) {
	type args struct {
		numberOfElements int
		settings         GridSettings
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantSettings GridSettings
	}{
		{
			name: "fill grid with rows and columns",
			args: args{
				settings: GridSettings{
					Rows:    3,
					Columns: 4,
				},
				numberOfElements: 12,
			},
			wantErr: false,
			wantSettings: GridSettings{
				Rows:    3,
				Columns: 4,
			},
		},
		{
			name: "fill grid with rows and auto assign columns",
			args: args{
				settings: GridSettings{
					Rows: 10,
				},
				numberOfElements: 30,
			},
			wantErr: false,
			wantSettings: GridSettings{
				Rows:    10,
				Columns: 3,
			},
		},
		{
			name: "fill without specifying rows or columns",
			args: args{
				settings:         GridSettings{},
				numberOfElements: 30,
			},
			wantErr: false,
			wantSettings: GridSettings{
				Rows:    1,
				Columns: 30,
			},
		},
		{
			name: "fill grid with column and auto assign rows",
			args: args{
				settings: GridSettings{
					Columns: 50,
				},
				numberOfElements: 200,
			},
			wantErr: false,
			wantSettings: GridSettings{
				Rows:    4,
				Columns: 50,
			},
		},
		{
			name: "fails to fill grid with mathematically incorrect rows and columns",
			args: args{
				settings: GridSettings{
					Rows:    3,
					Columns: 30,
				},
				numberOfElements: 12,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cells []mockCell
			for i := 0; i < tt.args.numberOfElements; i++ {
				elements = append(elements, element.Model{})
			}

			m := &Model{}

			err := m.SetGrid(elements, tt.args.settings)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetGrid() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && len(m.grid) != tt.wantSettings.Rows {
				t.Errorf("SetGrid() grid rows = %v, want %v", len(m.grid), tt.args.settings.Rows)
			}

			if err == nil && len(m.grid[0]) != tt.wantSettings.Columns {
				t.Errorf("SetGrid() grid columns = %v, want %v", len(m.grid[0]), tt.args.settings.Columns)
			}
		})
	}
}
