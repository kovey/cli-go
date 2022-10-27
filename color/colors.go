package color

type Color int32

const (
	Color_None    Color = 0
	Color_Green   Color = 1
	Color_White   Color = 2
	Color_Yellow  Color = 3
	Color_Red     Color = 4
	Color_Blue    Color = 5
	Color_Magenta Color = 6
	Color_Cyan    Color = 7
	Color_Reset   Color = 8
)

const (
	Bg_None    Color = 0
	Bg_Green   Color = 11
	Bg_White   Color = 12
	Bg_Yellow  Color = 13
	Bg_Red     Color = 14
	Bg_Blue    Color = 15
	Bg_Megenta Color = 16
	Bg_Cyan    Color = 17
)

var colors = map[Color]string{
	Color_Green:   string([]byte{27, 91, 51, 50, 109}),
	Color_White:   string([]byte{27, 91, 51, 55, 109}),
	Color_Yellow:  string([]byte{27, 91, 51, 51, 109}),
	Color_Red:     string([]byte{27, 91, 51, 49, 109}),
	Color_Blue:    string([]byte{27, 91, 51, 52, 109}),
	Color_Magenta: string([]byte{27, 91, 51, 53, 109}),
	Color_Cyan:    string([]byte{27, 91, 51, 54, 109}),
	Color_Reset:   string([]byte{27, 91, 48, 109}),
}

var bgs = map[Color]string{
	Bg_Green:   string([]byte{27, 91, 57, 55, 59, 52, 50, 109}),
	Bg_White:   string([]byte{27, 91, 57, 48, 59, 52, 55, 109}),
	Bg_Yellow:  string([]byte{27, 91, 57, 48, 59, 52, 51, 109}),
	Bg_Red:     string([]byte{27, 91, 57, 55, 59, 52, 49, 109}),
	Bg_Blue:    string([]byte{27, 91, 57, 55, 59, 52, 52, 109}),
	Bg_Megenta: string([]byte{27, 91, 57, 55, 59, 52, 53, 109}),
	Bg_Cyan:    string([]byte{27, 91, 57, 55, 59, 52, 54, 109}),
}

func color(c Color) string {
	if cc, ok := colors[c]; ok {
		return cc
	}

	return ""
}

func bg(c Color) string {
	if bg, ok := bgs[c]; ok {
		return bg
	}

	return ""
}
