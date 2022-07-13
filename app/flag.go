package app

import "flag"

type Type int

const (
	TYPE_INT    Type = 1
	TYPE_STRING Type = 2
	TYPE_BOOL   Type = 3
)

type Flag struct {
	name    string
	t       Type
	def     interface{}
	comment string
	value   interface{}
}

func (f *Flag) parse() {
	switch f.t {
	case TYPE_INT:
		val, ok := f.def.(int)
		if !ok {
			break
		}

		f.value = flag.Int(f.name, val, f.comment)
		break
	case TYPE_BOOL:
		val, ok := f.def.(bool)
		if !ok {
			break
		}
		f.value = flag.Bool(f.name, val, f.comment)
		break
	case TYPE_STRING:
		val, ok := f.def.(string)
		if !ok {
			break
		}

		f.value = flag.String(f.name, val, f.comment)
	}
}

func (f *Flag) String() string {
	if f.t != TYPE_STRING {
		return ""
	}

	return *(f.value.(*string))
}

func (f *Flag) Int() int {
	if f.t != TYPE_INT {
		return 0
	}

	return *(f.value.(*int))
}

func (f *Flag) Bool() bool {
	if f.t != TYPE_BOOL {
		return false
	}

	return *(f.value.(*bool))
}
