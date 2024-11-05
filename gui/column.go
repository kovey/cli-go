package gui

import (
	"fmt"
)

type Column struct {
	Data      string
	Len       int
	prev      string
	suff      string
	leftUp    string
	left      string
	leftDown  string
	up        *Border
	down      *Border
	rightUp   string
	right     string
	rightDown string
	p         Position
	upPrev    string
	downPrev  string
}

func NewColumn(data string, p Position) *Column {
	r := &Column{Data: data, p: p}
	r.init()
	return r
}

func (r *Column) Add(p Position) {
	r.p = p
	r.init()
}

func (r *Column) Reset(p Position) {
	r.leftUp = ""
	r.left = ""
	r.leftDown = ""
	r.up = nil
	r.down = nil
	r.right = ""
	r.rightDown = ""
	r.rightUp = ""
	r.Add(p)
}

func (r *Column) init() {
	switch r.p {
	case Position_Center:
		r.leftUp = Border_Center
		r.up = NewBorder(Border_Horizontal)
		r.left = Border_Vertical
	case Position_Down:
		r.leftUp = Border_Center
		r.up = NewBorder(Border_Horizontal)
		r.leftDown = Border_Un_T
		r.down = NewBorder(Border_Horizontal)
		r.left = Border_Vertical
	case Position_Left:
		r.leftUp = Border_Left_Center
		r.up = NewBorder(Border_Horizontal)
		r.left = Border_Vertical
	case Position_Left_Down:
		r.leftUp = Border_Left_Center
		r.leftDown = Border_Left_Bottom
		r.up = NewBorder(Border_Horizontal)
		r.down = NewBorder(Border_Horizontal)
		r.left = Border_Vertical
	case Position_Left_Up:
		r.leftUp = Border_Left_Up
		r.up = NewBorder(Border_Horizontal)
		r.left = Border_Vertical
	case Position_Right:
		r.leftUp = Border_Center
		r.up = NewBorder(Border_Horizontal)
		r.rightUp = Border_Right_Center
		r.left = Border_Vertical
		r.right = Border_Vertical
	case Position_Right_Down:
		r.leftUp = Border_Center
		r.up = NewBorder(Border_Horizontal)
		r.down = NewBorder(Border_Horizontal)
		r.rightUp = Border_Right_Center
		r.rightDown = Border_Right_Bottom
		r.leftDown = Border_Un_T
		r.left = Border_Vertical
		r.right = Border_Vertical
	case Position_Right_Up:
		r.leftUp = Border_T
		r.up = NewBorder(Border_Horizontal)
		r.rightUp = Border_Right_Up
		r.left = Border_Vertical
		r.right = Border_Vertical
	case Position_Up:
		r.leftUp = Border_T
		r.up = NewBorder(Border_Horizontal)
		r.left = Border_Vertical
	}
	texts := []rune(r.Data)
	other := 0
	for _, text := range texts {
		if IsChinese(string(text)) {
			other++
		}
	}
	r.Len = len(Border_Left_Up) + len(Border_Right_Up) + len(texts) + other
	r.upPrev = ""
	r.downPrev = ""
}

func (r *Column) Adjust(maxLen int) {
	if r.up != nil {
		r.up.Adjust(maxLen)
	}
	if r.down != nil {
		r.down.Adjust(maxLen)
	}

	if maxLen <= r.Len {
		return
	}

	sub := maxLen - r.Len
	r.prev += " "

	for i := 1; i < sub; i++ {
		r.suff += " "
	}
}

func (r *Column) UpBorder() string {
	return fmt.Sprintf("%s%s%s", r.leftUp, r.up.Text(), r.rightUp)
}

func (r *Column) DownBorder() string {
	return fmt.Sprintf("%s%s%s", r.leftDown, r.down.Text(), r.rightDown)
}

func (r *Column) Text() string {
	return fmt.Sprintf("%s%s%s%s%s", r.left, r.prev, r.Data, r.suff, r.right)
}
