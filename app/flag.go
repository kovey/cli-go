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
}

func (f *Flag) IsInput() bool {
	return f.has
}

func (f *Flag) String() string {
	if !f.hasValue {
		panic(fmt.Sprintf("%s has not value", f.name))
	}

	if f.t != TYPE_STRING {
		return ""
	}

	if !f.has {
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
		return f.def.(bool)
	}

	tmp, err := strconv.ParseBool(f.value)
	if err != nil {
		panic(err)
	}

	return tmp
}

func (f *Flag) print(maxLen int) {
	name := fmt.Sprintf("-%s", f.name)
	if !f.isShort {
		name = fmt.Sprintf("--%s", f.name)
	}

	sub := maxLen - len(name)
	for i := 0; i < sub; i++ {
		name += " "
	}

	if f.hasValue {
		fmt.Printf(`
	%s  %s, default: %v
`, name, f.comment, f.def)
		return
	}
	fmt.Printf(`
	%s  %s
`, name, f.comment)
}
