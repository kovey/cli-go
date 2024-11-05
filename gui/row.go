package gui

import "strings"

type Row struct {
	columns []*Column
	index   int
}

func NewRow(index int) *Row {
	r := &Row{index: index}
	return r
}

func (r *Row) Add(data string) {
	switch r.index {
	case 0:
		if len(r.columns) == 0 {
			r.columns = append(r.columns, NewColumn(data, Position_Left_Up))
			return
		}

		r.columns = append(r.columns, NewColumn(data, Position_Up))
	default:
		if len(r.columns) == 0 {
			r.columns = append(r.columns, NewColumn(data, Position_Left))
			return
		}

		r.columns = append(r.columns, NewColumn(data, Position_Center))
	}
}

func (r *Row) Final(rowCount int) {
	length := len(r.columns)
	if length < 1 {
		return
	}

	if rowCount == 1 {
		r.columns[0].Add(Position_Left_Down)
		r.columns[length-1].Add(Position_Right_Down)
		for i := 1; i < length-1; i++ {
			r.columns[i].Add(Position_Down)
		}
		return
	}

	switch r.index {
	case 0:
		r.columns[length-1].Reset(Position_Right_Up)
	default:
		if rowCount-1 == r.index {
			r.columns[0].Reset(Position_Left_Down)
			for i := 1; i < length-1; i++ {
				r.columns[i].Reset(Position_Down)
			}
			r.columns[length-1].Reset(Position_Right_Down)
			return
		}

		r.columns[length-1].Reset(Position_Right)
	}
}

func (t *Row) ColumnLen(index int) int {
	if index >= len(t.columns) {
		return 0
	}

	return t.columns[index].Len
}

func (t *Row) UpBorder() string {
	var borders = make([]string, len(t.columns))
	for i, column := range t.columns {
		borders[i] = column.UpBorder()
	}

	return strings.Join(borders, "")
}

func (t *Row) DownBorder() string {
	var borders = make([]string, len(t.columns))
	for i, column := range t.columns {
		borders[i] = column.DownBorder()
	}

	return strings.Join(borders, "")
}

func (t *Row) Text() string {
	var texts = make([]string, len(t.columns))
	for i, column := range t.columns {
		texts[i] = column.Text()
	}

	return strings.Join(texts, "")
}

func (t *Row) Adjust(index, maxLen int) {
	if index >= len(t.columns) {
		return
	}

	t.columns[index].Adjust(maxLen)
}
