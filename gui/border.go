package gui

import "regexp"

const (
	Border_Left_Up      = "┌"
	Border_Left_Bottom  = "└"
	Border_Right_Up     = "┐"
	Border_Right_Bottom = "┘"
	Border_Horizontal   = "─"
	Border_Vertical     = "│"

	Reg_Chinese      = "^[\u4E00-\u9FA5]+$"
	Reg_Chinese_Sign = "^[\u3002|\uff1f|\uff01|\uff0c|\u3001|\uff1b|\uff1a|\u201c|\u201d|\u2018|\u2019|\uff08|\uff09|\u300a|\u300b|\u3010|\u3011|\u007e]+$"
)

func IsChinese(text string) bool {
	if ok, _ := regexp.Match(Reg_Chinese, []byte(text)); ok {
		return ok
	}

	ok, _ := regexp.Match(Reg_Chinese_Sign, []byte(text))
	return ok
}
