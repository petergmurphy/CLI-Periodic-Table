package table

import (
	"fmt"
	"periodic-table/ui/element"
)

var (
	noElementsError                 = fmt.Errorf("grid cannot be empty")
	elementsDoNotFitBoundariesError = fmt.Errorf("number of elements does not fit grid size")
)

type GridSettings struct {
	Rows    int
	Columns int
}

func (m *Model) SetGrid(elements []element.Model, settings GridSettings) error {
	if elements == nil {
		return fmt.Errorf("grid cannot be nil")
	}

	err := allocateColumnsAndRows(len(elements), &settings)
	if err != nil {
		return err
	}

	m.grid = fillGrid(elements, settings)

	return nil
}

func fillGrid(elements []element.Model, settings GridSettings) [][]element.Model {
	var grid [][]element.Model
	for i := 0; i < settings.Rows; i++ {
		var row []element.Model
		for j := 0; j < settings.Columns; j++ {
			row = append(row, elements[i*settings.Columns+j])
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
