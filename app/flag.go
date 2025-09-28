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
	TYPE_FLOAT  Type = 4
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

func (f *Flag) _value(typ Type) any {
	if !f.hasValue {
		panic(fmt.Sprintf("%s has not value", f.name))
	}

	if f.t != typ {
		panic(fmt.Sprintf("%s type error", f.name))
	}

	if !f.has {
		if f.def == nil {
			panic(fmt.Sprintf("%s not input value", f.name))
		}

		return f.def
	}

	switch f.t {
	case TYPE_BOOL:
		tmp, err := strconv.ParseBool(f.value)
		if err != nil {
			panic(err)
		}

		return tmp
	case TYPE_FLOAT:
		tmp, err := strconv.ParseFloat(f.value, 64)
		if err != nil {
			panic(err)
		}

		return tmp
	case TYPE_INT:
		tmp, err := strconv.Atoi(f.value)
		if err != nil {
			panic(err)
		}

		return tmp
	default:
		return f.value
	}
}

func (f *Flag) String() string {
	if f.isArg {
		return f.name
	}

	return f._value(TYPE_STRING).(string)
}

func (f *Flag) Int() int {
	return f._value(TYPE_INT).(int)
}

func (f *Flag) Bool() bool {
	return f._value(TYPE_BOOL).(bool)
}

func (f *Flag) Float() float64 {
	return f._value(TYPE_FLOAT).(float64)
}
