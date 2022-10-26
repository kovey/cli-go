package gui

type Table struct {
	header *Row
	rows   []*Row
	foot   *Row
}

func NewTable() *Table {
	return &Table{rows: make([]*Row, 0), header: NewRepeatRow("-"), foot: NewRepeatRow("-")}
}

func (t *Table) Add(data string) {
	t.rows = append(t.rows, NewRow(data))
}

func (t *Table) rowMaxLen() int {
	var maxLen = 0
	for _, row := range t.rows {
		if row.Len > maxLen {
			maxLen = row.Len
		}
	}

	return maxLen + 2
}

func (t *Table) Show() {
	maxLen := t.rowMaxLen()
	t.header.Adjust(maxLen)
	t.header.Show()
	for _, row := range t.rows {
		row.Adjust(maxLen)
		row.Show()
	}
	t.foot.Adjust(maxLen)
	t.foot.Show()
}
