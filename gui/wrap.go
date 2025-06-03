package gui

type Wrap struct {
	columns []*Column
	Weight  int
	pos     Position
	data    string
}

func NewWrap(data string, weight int, pos Position) *Wrap {
	return &Wrap{data: data, Weight: weight, pos: pos}
}

func (w *Wrap) init() {
	textLen := len(w.data)
	if textLen <= w.Weight {
		w.columns = append(w.columns, NewColumn(w.data, w.pos))
		return
	}

	count := textLen / w.Weight
	if textLen%w.Weight != 0 {
		count++
	}

	switch w.pos {
	case Position_Center:
		w.center(textLen, count)
	case Position_Down:
		w.down(textLen, count)
	case Position_Left:
		w.center(textLen, count)
	case Position_Left_Down:
		w.leftDown(textLen, count)
	case Position_Left_Up:
		w.leftUp(textLen, count)
	case Position_Right:
		w.center(textLen, count)
	case Position_Right_Down:
		w.leftDown(textLen, count)
	case Position_Right_Up:
		w.leftUp(textLen, count)
	case Position_Up:
		w.up(textLen, count)
	}
}

func (w *Wrap) leftUp(textLen, count int) {
	w.columns = append(w.columns, NewColumn(w.data[0:w.Weight], w.pos))
	for i := 1; i < count; i++ {
		start := i * w.Weight
		end := start + w.Weight
		if end > textLen {
			end = textLen
		}

		w.columns = append(w.columns, NewColumn(w.data[start:end], Position_Left))
	}
}

func (w *Wrap) up(textLen, count int) {
	w.columns = append(w.columns, NewColumn(w.data[0:w.Weight], w.pos))
	for i := 1; i < count; i++ {
		start := i * w.Weight
		end := start + w.Weight
		if end > textLen {
			end = textLen
		}

		w.columns = append(w.columns, NewColumn(w.data[start:end], Position_Center))
	}
}

func (w *Wrap) center(textLen, count int) {
	for i := 0; i < count; i++ {
		start := i * w.Weight
		end := start + w.Weight
		if end > textLen {
			end = textLen
		}

		w.columns = append(w.columns, NewColumn(w.data[start:end], w.pos))
	}
}

func (w *Wrap) leftDown(textLen, count int) {
	for i := 0; i < count-1; i++ {
		start := i * w.Weight
		end := start + w.Weight
		if end > textLen {
			end = textLen
		}

		w.columns = append(w.columns, NewColumn(w.data[start:end], Position_Left))
	}

	w.columns = append(w.columns, NewColumn(w.data[(count-1)*w.Weight:textLen], w.pos))
}

func (w *Wrap) down(textLen, count int) {
	for i := 0; i < count-1; i++ {
		start := i * w.Weight
		end := start + w.Weight
		if end > textLen {
			end = textLen
		}

		w.columns = append(w.columns, NewColumn(w.data[start:end], Position_Center))
	}

	w.columns = append(w.columns, NewColumn(w.data[(count-1)*w.Weight:textLen], w.pos))
}
