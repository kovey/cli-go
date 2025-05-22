package app

import (
	"fmt"
	"strconv"
)

type Type int

const (
	TYPE_INT    Type = 1
	TYPE_STRING Type = 2
	TYPE_BOOL   Type = 3
)

type Flag struct {
	name     string
	t        Type
	def      any
	comment  string
	value    string
	has      bool
	hasValue bool
	isShort  bool
	index    int
	isArg    bool
	children *Flags
}

func (f *Flag) AddChild(flag *Flag) {
	if f.children == nil {
		f.children = NewFlags()
	}

	f.children.Add(flag)
}

func (f *Flag) IsInput() bool {
	return f.has
}

func (f *Flag) String() string {
	if f.isArg {
		return f.name
	}

	if !f.hasValue {
		panic(fmt.Sprintf("%s has not value", f.name))
	}

	if f.t != TYPE_STRING {
		return ""
	}

	if !f.has {
		if f.def == nil {
			panic(fmt.Sprintf("%s not input value", f.name))
		}
		return f.def.(string)
	}

	return f.value
}

func (f *Flag) Int() int {
	if !f.hasValue {
		panic(fmt.Sprintf("%s has not value", f.name))
	}

	if f.t != TYPE_INT {
		return 0
	}

	if !f.has {
		if f.def == nil {
			panic(fmt.Sprintf("%s not input value", f.name))
		}

		return f.def.(int)
	}

	tmp, err := strconv.Atoi(f.value)
	if err != nil {
		panic(err)
	}

	return tmp
}

func (f *Flag) Bool() bool {
	if !f.hasValue {
		panic(fmt.Sprintf("%s has not value", f.name))
	}

	if f.t != TYPE_BOOL {
		return false
	}

	if !f.has {
		if f.def == nil {
			panic(fmt.Sprintf("%s not input value", f.name))
		}

		return f.def.(bool)
	}

	tmp, err := strconv.ParseBool(f.value)
	if err != nil {
		panic(err)
	}

	return tmp
}
