package app

type Flags struct {
	flags map[string]*Flag
}

func NewFlags() *Flags {
	return &Flags{flags: make(map[string]*Flag)}
}

func (f *Flags) Clean() {
	f.flags = make(map[string]*Flag)
}

func (f *Flags) Add(flag *Flag) {
	if _, ok := f.flags[flag.name]; !ok {
		f.flags[flag.name] = flag
	}
}

func (f *Flags) Has(name string) bool {
	_, ok := f.flags[name]
	return ok
}

func (f *Flags) Get(parents ...string) *Flag {
	flags := f.flags
	var last *Flag
	count := len(parents)
	for _, name := range parents {
		if name == Ko_Command_Help && count > 1 {
			continue
		}

		if flags == nil {
			return nil
		}

		flag, ok := flags[name]
		if !ok {
			return nil
		}

		last = flag
		if last.children == nil {
			flags = nil
			continue
		}

		flags = last.children.flags
	}

	return last
}
