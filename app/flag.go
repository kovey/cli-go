package app

import "flag"

type Type int

const (
	TYPE_INT    Type = 1
	TYPE_STRING Type = 2
	TYPE_BOOL   Type = 3
)

type Flag struct {
	Name    string
	Type    Type
	Default interface{}
	Comment string
	Value   interface{}
}

func (f *Flag) Parse() {
	switch f.Type {
	case TYPE_INT:
		val, ok := f.Default.(int)
		if !ok {
			break
		}

		f.Value = flag.Int(f.Name, val, f.Comment)
		break
	case TYPE_BOOL:
		val, ok := f.Default.(bool)
		if !ok {
			break
		}
		f.Value = flag.Bool(f.Name, val, f.Comment)
		break
	case TYPE_STRING:
		val, ok := f.Default.(string)
		if !ok {
			break
		}

		f.Value = flag.String(f.Name, val, f.Comment)
	}
}
