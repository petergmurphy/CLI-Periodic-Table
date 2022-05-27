package table

type Stack [][]Cell

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(cells []Cell) {
	*s = append(*s, cells)
}

func (s *Stack) Pop() ([]Cell, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1
		cells := (*s)[index]
		*s = (*s)[:index]
		return cells, true
	}
}
