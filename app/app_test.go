package app

import (
	"strings"
	"testing"
)

// =============================================================================
// flag.go tests
// =============================================================================

func TestFlagString(t *testing.T) {
	f := &Flag{name: "test", t: TYPE_STRING, has: false, hasValue: true, def: "default_value"}
	if got := f.String(); got != "default_value" {
		t.Errorf("String() = %q, want %q", got, "default_value")
	}
}

func TestFlagInt(t *testing.T) {
	f := &Flag{name: "test", t: TYPE_INT, has: false, hasValue: true, def: 42}
	if got := f.Int(); got != 42 {
		t.Errorf("Int() = %d, want %d", got, 42)
	}
}

func TestFlagBool(t *testing.T) {
	f := &Flag{name: "test", t: TYPE_BOOL, has: false, hasValue: true, def: true}
	if got := f.Bool(); got != true {
		t.Errorf("Bool() = %v, want %v", got, true)
	}
}

func TestFlagFloat(t *testing.T) {
	f := &Flag{name: "test", t: TYPE_FLOAT, has: false, hasValue: true, def: 3.14}
	if got := f.Float(); got != 3.14 {
		t.Errorf("Float() = %v, want %v", got, 3.14)
	}
}

func TestFlagStringWithValue(t *testing.T) {
	f := &Flag{name: "test", t: TYPE_STRING, has: true, hasValue: true, def: "default", value: "provided_value"}
	if got := f.String(); got != "provided_value" {
		t.Errorf("String() = %q, want %q", got, "provided_value")
	}
}

func TestFlagIntWithValue(t *testing.T) {
	f := &Flag{name: "test", t: TYPE_INT, has: true, hasValue: true, def: 1, value: "99"}
	if got := f.Int(); got != 99 {
		t.Errorf("Int() = %d, want %d", got, 99)
	}
}

func TestFlagBoolWithValue(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_BOOL, has: true, hasValue: true, def: false, value: "true"}
		if got := f.Bool(); got != true {
			t.Errorf("Bool() = %v, want %v", got, true)
		}
	})
	t.Run("false", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_BOOL, has: true, hasValue: true, def: true, value: "false"}
		if got := f.Bool(); got != false {
			t.Errorf("Bool() = %v, want %v", got, false)
		}
	})
	t.Run("1", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_BOOL, has: true, hasValue: true, def: false, value: "1"}
		if got := f.Bool(); got != true {
			t.Errorf("Bool() = %v, want %v", got, true)
		}
	})
	t.Run("0", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_BOOL, has: true, hasValue: true, def: true, value: "0"}
		if got := f.Bool(); got != false {
			t.Errorf("Bool() = %v, want %v", got, false)
		}
	})
}

func TestFlagFloatWithValue(t *testing.T) {
	f := &Flag{name: "test", t: TYPE_FLOAT, has: true, hasValue: true, def: 1.0, value: "2.718"}
	if got := f.Float(); got != 2.718 {
		t.Errorf("Float() = %v, want %v", got, 2.718)
	}
}

func TestFlagIsInput(t *testing.T) {
	t.Run("has true", func(t *testing.T) {
		f := &Flag{has: true}
		if !f.IsInput() {
			t.Error("IsInput() should return true when has is true")
		}
	})
	t.Run("has false", func(t *testing.T) {
		f := &Flag{has: false}
		if f.IsInput() {
			t.Error("IsInput() should return false when has is false")
		}
	})
}

func TestFlagAddChild(t *testing.T) {
	parent := &Flag{name: "parent"}
	child := &Flag{name: "child", t: TYPE_STRING, hasValue: true, def: "child_val"}
	parent.AddChild(child)

	if parent.children == nil {
		t.Fatal("expected children to be non-nil after AddChild")
	}
	if !parent.children.Has("child") {
		t.Error("expected child flag to exist in parent's children")
	}

	got := parent.children.Get("child")
	if got == nil {
		t.Fatal("expected to get child from parent's children")
	}
	if got.String() != "child_val" {
		t.Errorf("child String() = %q, want %q", got.String(), "child_val")
	}
}

func TestFlagCheckDefault(t *testing.T) {
	t.Run("wrong bool default", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for wrong bool default type")
			}
		}()
		f := &Flag{name: "test", t: TYPE_BOOL, def: "not a bool"}
		f.checkDefault()
	})

	t.Run("wrong int default", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for wrong int default type")
			}
		}()
		f := &Flag{name: "test", t: TYPE_INT, def: "not an int"}
		f.checkDefault()
	})

	t.Run("wrong float default", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for wrong float default type")
			}
		}()
		f := &Flag{name: "test", t: TYPE_FLOAT, def: "not a float"}
		f.checkDefault()
	})

	t.Run("wrong string default", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for wrong string default type")
			}
		}()
		f := &Flag{name: "test", t: TYPE_STRING, def: 123}
		f.checkDefault()
	})

	t.Run("correct default no panic", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_STRING, def: "valid"}
		f.checkDefault()
	})

	t.Run("nil default no panic", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_INT, def: nil}
		f.checkDefault()
	})

	t.Run("has true skips check", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_INT, has: true, def: "would be wrong"}
		f.checkDefault()
	})
}

func TestFlagTypeMismatchPanics(t *testing.T) {
	t.Run("Int on string flag", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when calling Int() on a string flag")
			}
		}()
		f := &Flag{name: "test", t: TYPE_STRING, has: true, hasValue: true, def: "hello", value: "world"}
		f.Int()
	})

	t.Run("String on int flag", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when calling String() on an int flag")
			}
		}()
		f := &Flag{name: "test", t: TYPE_INT, has: true, hasValue: true, def: 1, value: "2"}
		_ = f.String()
	})

	t.Run("Bool on string flag", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when calling Bool() on a string flag")
			}
		}()
		f := &Flag{name: "test", t: TYPE_STRING, has: true, hasValue: true, def: "hello", value: "world"}
		f.Bool()
	})

	t.Run("Float on bool flag", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when calling Float() on a bool flag")
			}
		}()
		f := &Flag{name: "test", t: TYPE_BOOL, has: true, hasValue: true, def: true, value: "false"}
		f.Float()
	})
}

func TestFlagDefaultValueUsedWhenNotSet(t *testing.T) {
	t.Run("string default", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_STRING, has: false, hasValue: true, def: "default_str"}
		if got := f.String(); got != "default_str" {
			t.Errorf("String() = %q, want %q", got, "default_str")
		}
	})

	t.Run("int default", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_INT, has: false, hasValue: true, def: 100}
		if got := f.Int(); got != 100 {
			t.Errorf("Int() = %d, want %d", got, 100)
		}
	})

	t.Run("bool default", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_BOOL, has: false, hasValue: true, def: true}
		if got := f.Bool(); got != true {
			t.Errorf("Bool() = %v, want %v", got, true)
		}
	})

	t.Run("float default", func(t *testing.T) {
		f := &Flag{name: "test", t: TYPE_FLOAT, has: false, hasValue: true, def: 1.5}
		if got := f.Float(); got != 1.5 {
			t.Errorf("Float() = %v, want %v", got, 1.5)
		}
	})

	t.Run("nil default panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when default is nil and has is false")
			}
		}()
		f := &Flag{name: "test", t: TYPE_STRING, has: false, hasValue: true, def: nil}
		_ = f.String()
	})
}

func TestFlagParseErrors(t *testing.T) {
	t.Run("invalid bool value panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid bool value")
			}
		}()
		f := &Flag{name: "test", t: TYPE_BOOL, has: true, hasValue: true, def: false, value: "notabool"}
		f.Bool()
	})

	t.Run("invalid int value panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid int value")
			}
		}()
		f := &Flag{name: "test", t: TYPE_INT, has: true, hasValue: true, def: 1, value: "notanint"}
		f.Int()
	})

	t.Run("invalid float value panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid float value")
			}
		}()
		f := &Flag{name: "test", t: TYPE_FLOAT, has: true, hasValue: true, def: 1.0, value: "notafloat"}
		f.Float()
	})
}

func TestFlagStringIsArg(t *testing.T) {
	f := &Flag{name: "command_name", t: TYPE_STRING, isArg: true, has: true, hasValue: true, def: "def", value: "val"}
	if got := f.String(); got != "command_name" {
		t.Errorf("isArg flag String() = %q, want %q", got, "command_name")
	}
}

// =============================================================================
// flags.go tests
// =============================================================================

func TestFlagsAdd(t *testing.T) {
	flags := NewFlags()
	f1 := &Flag{name: "flag1", t: TYPE_STRING, hasValue: true, def: "val1"}
	f2 := &Flag{name: "flag2", t: TYPE_INT, hasValue: true, def: 42}

	flags.Add(f1)
	flags.Add(f2)

	if !flags.Has("flag1") {
		t.Error("expected flag1 to exist")
	}
	if !flags.Has("flag2") {
		t.Error("expected flag2 to exist")
	}
}

func TestFlagsHas(t *testing.T) {
	flags := NewFlags()
	flags.Add(&Flag{name: "exists", t: TYPE_STRING, hasValue: true, def: "value"})

	t.Run("existing flag", func(t *testing.T) {
		if !flags.Has("exists") {
			t.Error("Has() should return true for existing flag")
		}
	})

	t.Run("non-existing flag", func(t *testing.T) {
		if flags.Has("nonexistent") {
			t.Error("Has() should return false for non-existing flag")
		}
	})
}

func TestFlagsGet(t *testing.T) {
	flags := NewFlags()
	expected := &Flag{name: "test", t: TYPE_STRING, hasValue: true, def: "hello"}
	flags.Add(expected)

	got := flags.Get("test")
	if got == nil {
		t.Fatal("Get() returned nil for existing flag")
	}
	if got.String() != "hello" {
		t.Errorf("Get() flag String() = %q, want %q", got.String(), "hello")
	}
}

func TestFlagsGetHierarchical(t *testing.T) {
	flags := NewFlags()

	parent := &Flag{name: "parent"}
	child := &Flag{name: "child", t: TYPE_STRING, hasValue: true, def: "child_value"}
	parent.AddChild(child)
	flags.Add(parent)

	got := flags.Get("parent", "child")
	if got == nil {
		t.Fatal("Get() returned nil for hierarchical path")
	}
	if got.String() != "child_value" {
		t.Errorf("Get() flag String() = %q, want %q", got.String(), "child_value")
	}

	parentGot := flags.Get("parent")
	if parentGot == nil {
		t.Fatal("Get() returned nil for parent flag")
	}
}

func TestFlagsGetNotFound(t *testing.T) {
	flags := NewFlags()
	flags.Add(&Flag{name: "exists", t: TYPE_STRING, hasValue: true, def: "val"})

	t.Run("missing top-level", func(t *testing.T) {
		if got := flags.Get("nonexistent"); got != nil {
			t.Error("Get() should return nil for non-existent flag")
		}
	})

	t.Run("missing child", func(t *testing.T) {
		if got := flags.Get("exists", "nochild"); got != nil {
			t.Error("Get() should return nil when child path doesn't exist")
		}
	})
}

func TestFlagsGetHelpSkip(t *testing.T) {
	flags := NewFlags()
	parent := &Flag{name: "create"}
	helpChild := &Flag{name: Ko_Command_Help, t: TYPE_STRING, hasValue: true, def: "help_val"}
	parent.AddChild(helpChild)
	flags.Add(parent)

	got := flags.Get("create", Ko_Command_Help, Ko_Command_Help)
	if got == nil {
		t.Error("Get should find help child when help is skipped in multi-element path")
	}

	got2 := flags.Get(Ko_Command_Help)
	if got2 != nil {
		t.Error("Get(help) with count=1 should try to look up help directly and return nil")
	}
}

func TestFlagsClean(t *testing.T) {
	flags := NewFlags()
	flags.Add(&Flag{name: "flag1", t: TYPE_STRING, hasValue: true, def: "val1"})
	flags.Add(&Flag{name: "flag2", t: TYPE_INT, hasValue: true, def: 42})

	if !flags.Has("flag1") || !flags.Has("flag2") {
		t.Fatal("flags should exist before Clean")
	}

	flags.Clean()

	if flags.Has("flag1") {
		t.Error("flag1 should not exist after Clean")
	}
	if flags.Has("flag2") {
		t.Error("flag2 should not exist after Clean")
	}
	if flags.Get("flag1") != nil {
		t.Error("Get should return nil after Clean")
	}
}

func TestFlagsDuplicateAdd(t *testing.T) {
	flags := NewFlags()
	first := &Flag{name: "dup", t: TYPE_STRING, hasValue: true, def: "first"}
	second := &Flag{name: "dup", t: TYPE_INT, hasValue: true, def: 99}

	flags.Add(first)
	flags.Add(second)

	got := flags.Get("dup")
	if got == nil {
		t.Fatal("Get() returned nil")
	}
	if got.t != TYPE_STRING {
		t.Errorf("expected flag type to remain TYPE_STRING, got %v", got.t)
	}
	if got.String() != "first" {
		t.Errorf("expected first flag's value %q, got %q", "first", got.String())
	}
}

// =============================================================================
// command_line.go tests
// =============================================================================

func TestCheckShort(t *testing.T) {
	cl := NewCommandLine()

	t.Run("valid single letter", func(t *testing.T) {
		if err := cl.checkShort("v"); err != nil {
			t.Errorf("checkShort(v) should be valid, got error: %v", err)
		}
	})

	t.Run("valid multi letter", func(t *testing.T) {
		if err := cl.checkShort("verbose"); err != nil {
			t.Errorf("checkShort(verbose) should be valid, got error: %v", err)
		}
	})

	t.Run("valid uppercase", func(t *testing.T) {
		if err := cl.checkShort("ABC"); err != nil {
			t.Errorf("checkShort(ABC) should be valid, got error: %v", err)
		}
	})

	t.Run("invalid with hyphen", func(t *testing.T) {
		if err := cl.checkShort("not-valid"); err == nil {
			t.Error("checkShort(not-valid) should return error")
		}
	})

	t.Run("invalid with digits", func(t *testing.T) {
		if err := cl.checkShort("abc123"); err == nil {
			t.Error("checkShort(abc123) should return error")
		}
	})

	t.Run("invalid empty", func(t *testing.T) {
		if err := cl.checkShort(""); err == nil {
			t.Error("checkShort('') should return error")
		}
	})
}

func TestCheckLong(t *testing.T) {
	cl := NewCommandLine()

	t.Run("valid simple", func(t *testing.T) {
		if err := cl.checkLong("name"); err != nil {
			t.Errorf("checkLong(name) should be valid, got error: %v", err)
		}
	})

	t.Run("valid with hyphens", func(t *testing.T) {
		if err := cl.checkLong("config-path"); err != nil {
			t.Errorf("checkLong(config-path) should be valid, got error: %v", err)
		}
	})

	t.Run("valid two letters", func(t *testing.T) {
		if err := cl.checkLong("ab"); err != nil {
			t.Errorf("checkLong(ab) should be valid, got error: %v", err)
		}
	})

	t.Run("invalid single letter", func(t *testing.T) {
		if err := cl.checkLong("v"); err == nil {
			t.Error("checkLong(v) should return error (needs at least 2 chars)")
		}
	})

	t.Run("invalid ends with hyphen", func(t *testing.T) {
		if err := cl.checkLong("bad-"); err == nil {
			t.Error("checkLong(bad-) should return error")
		}
	})

	t.Run("invalid only digits", func(t *testing.T) {
		if err := cl.checkLong("123"); err == nil {
			t.Error("checkLong(123) should return error")
		}
	})

	t.Run("invalid empty", func(t *testing.T) {
		if err := cl.checkLong(""); err == nil {
			t.Error("checkLong('') should return error")
		}
	})
}

func TestCommandLineParse_ShortFlag(t *testing.T) {
	cl := NewCommandLine()
	cl.Flag("v", "default", TYPE_STRING, "verbose flag")
	cl.Parse([]string{"-v", "hello"})

	flag := cl.Get("v")
	if flag == nil {
		t.Fatal("expected flag v to exist")
	}
	if !flag.IsInput() {
		t.Error("expected flag v to be set as input")
	}
	if flag.String() != "hello" {
		t.Errorf("flag value = %q, want %q", flag.String(), "hello")
	}
}

func TestCommandLineParse_LongFlag(t *testing.T) {
	cl := NewCommandLine()
	cl.FlagLong("name", "default_name", TYPE_STRING, "name flag")
	cl.Parse([]string{"--name=test_value"})

	flag := cl.Get("name")
	if flag == nil {
		t.Fatal("expected flag name to exist")
	}
	if !flag.IsInput() {
		t.Error("expected flag name to be set as input")
	}
	if flag.String() != "test_value" {
		t.Errorf("flag value = %q, want %q", flag.String(), "test_value")
	}
}

func TestCommandLineParse_LongFlagEquals_ValueWithEquals(t *testing.T) {
	cl := NewCommandLine()
	cl.FlagLong("config", "", TYPE_STRING, "config flag")
	cl.Parse([]string{"--config=a=b"})

	flag := cl.Get("config")
	if flag == nil {
		t.Fatal("expected flag config to exist")
	}
	if flag.String() != "a=b" {
		t.Errorf("flag value = %q, want %q (Bug 4: value after = should be preserved)", flag.String(), "a=b")
	}
}

func TestCommandLineParse_NonValueFlag(t *testing.T) {
	t.Run("short non-value flag", func(t *testing.T) {
		cl := NewCommandLine()
		cl.FlagNonValue("v", "verbose mode")
		cl.Parse([]string{"-v"})

		flag := cl.Get("v")
		if flag == nil {
			t.Fatal("expected flag v to exist")
		}
		if !flag.IsInput() {
			t.Error("expected non-value flag to be set as input")
		}
	})

	t.Run("long non-value flag", func(t *testing.T) {
		cl := NewCommandLine()
		cl.FlagNonValueLong("debug", "debug mode")
		cl.Parse([]string{"--debug"})

		flag := cl.Get("debug")
		if flag == nil {
			t.Fatal("expected flag debug to exist")
		}
		if !flag.IsInput() {
			t.Error("expected non-value long flag to be set as input")
		}
	})
}

func TestCommandLineParse_ArgsInSequence(t *testing.T) {
	cl := NewCommandLine()
	cl.FlagArg("create", "create command", 0)
	cl.FlagArg("config", "create config command", 0, "create")
	cl.Parse([]string{"create", "config"})

	arg0 := cl.Arg(0)
	if arg0 == nil {
		t.Fatal("expected arg 0 to exist")
	}
	if arg0.String() != "create" {
		t.Errorf("Arg(0) = %q, want %q", arg0.String(), "create")
	}

	arg1 := cl.Arg(1)
	if arg1 == nil {
		t.Fatal("expected arg 1 to exist")
	}
	if arg1.String() != "config" {
		t.Errorf("Arg(1) = %q, want %q", arg1.String(), "config")
	}

	if cl.Arg(2) != nil {
		t.Error("Arg(2) should return nil for out-of-bounds index")
	}

	args := cl.Args()
	if len(args) != 2 {
		t.Errorf("Args() length = %d, want 2", len(args))
	}
}

func TestCommandLineParse_UnknownFlag(t *testing.T) {
	cl := NewCommandLine()
	cl.args = []string{"-x"}
	_, err := cl.parseOne()
	if err == nil {
		t.Error("expected error when parsing unknown flag")
	}
}

func TestCommandLineParse_SingleDash(t *testing.T) {
	cl := NewCommandLine()
	cl.args = []string{"-"}
	_, err := cl.parseOne()
	if err == nil {
		t.Error("expected error when parsing single dash (Bug 2: should error, not panic)")
	}
}

func TestCommandLineParse_EmptyString(t *testing.T) {
	t.Run("empty args slice", func(t *testing.T) {
		cl := NewCommandLine()
		cl.args = []string{}
		end, err := cl.parseOne()
		if err != nil {
			t.Errorf("expected no error for empty args, got: %v", err)
		}
		if !end {
			t.Error("expected end=true for empty args")
		}
	})

	t.Run("Parse with empty args", func(t *testing.T) {
		cl := NewCommandLine()
		cl.Parse([]string{})
	})
}

func TestCommandLineParse_MultipleFlags(t *testing.T) {
	cl := NewCommandLine()
	cl.Flag("v", "default_v", TYPE_STRING, "verbose")
	cl.FlagNonValue("d", "debug mode")
	cl.FlagLong("name", "default_name", TYPE_STRING, "name flag")
	cl.Parse([]string{"-v", "hello", "-d", "--name=world"})

	flagV := cl.Get("v")
	if flagV == nil || !flagV.IsInput() {
		t.Error("expected flag v to be set")
	}
	if flagV.String() != "hello" {
		t.Errorf("flag v = %q, want %q", flagV.String(), "hello")
	}

	flagD := cl.Get("d")
	if flagD == nil || !flagD.IsInput() {
		t.Error("expected flag d to be set")
	}

	flagName := cl.Get("name")
	if flagName == nil || !flagName.IsInput() {
		t.Error("expected flag name to be set")
	}
	if flagName.String() != "world" {
		t.Errorf("flag name = %q, want %q", flagName.String(), "world")
	}
}

func TestCommandLineFlagRegistration(t *testing.T) {
	cl := NewCommandLine()
	cl.Flag("v", "default", TYPE_STRING, "short flag")
	cl.FlagLong("name", "default", TYPE_STRING, "long flag")
	cl.FlagNonValue("d", "non-value short")
	cl.FlagNonValueLong("debug", "non-value long")

	if cl.Get("v") == nil {
		t.Error("short flag v should be registered")
	}
	if cl.Get("name") == nil {
		t.Error("long flag name should be registered")
	}
	if cl.Get("d") == nil {
		t.Error("non-value short flag d should be registered")
	}
	if cl.Get("debug") == nil {
		t.Error("non-value long flag debug should be registered")
	}
}

func TestCommandLineFlagArg(t *testing.T) {
	cl := NewCommandLine()
	cl.FlagArg("create", "create command", 0)

	flag := cl.Get("create")
	if flag == nil {
		t.Fatal("expected command arg create to be registered")
	}
	if !flag.isArg {
		t.Error("expected flag to be an arg type")
	}
}

func TestCommandLineCleanDefaults(t *testing.T) {
	cl := NewCommandLine()
	cl.Flag("v", "default", TYPE_STRING, "verbose")
	cl.FlagLong("name", "default", TYPE_STRING, "name")

	if cl.Get("v") == nil || cl.Get("name") == nil {
		t.Fatal("flags should exist before CleanDefaults")
	}

	cl.CleanDefaults()

	if cl.Get("v") != nil {
		t.Error("flag v should not exist after CleanDefaults")
	}
	if cl.Get("name") != nil {
		t.Error("flag name should not exist after CleanDefaults")
	}
}

func TestCleanFlagBitmask(t *testing.T) {
	t.Run("constants", func(t *testing.T) {
		if Clean_Defalut != 0 {
			t.Errorf("Clean_Defalut = %d, want 0", Clean_Defalut)
		}
		if Clean_Help != 1 {
			t.Errorf("Clean_Help = %d, want 1", Clean_Help)
		}
		if Clean_Version != 2 {
			t.Errorf("Clean_Version = %d, want 2", Clean_Version)
		}
	})

	t.Run("bitwise OR combination", func(t *testing.T) {
		combined := Clean_Help | Clean_Version
		if combined != 3 {
			t.Errorf("Clean_Help | Clean_Version = %d, want 3", combined)
		}
	})

	t.Run("bitwise AND test", func(t *testing.T) {
		combined := Clean_Help | Clean_Version
		if combined&Clean_Help == 0 {
			t.Error("combined should have Clean_Help bit set")
		}
		if combined&Clean_Version == 0 {
			t.Error("combined should have Clean_Version bit set")
		}
		if combined&Clean_Defalut != 0 {
			t.Error("Clean_Defalut is 0, AND with anything should be 0")
		}
	})

	t.Run("individual flags", func(t *testing.T) {
		if Clean_Help&Clean_Help == 0 {
			t.Error("Clean_Help should have its own bit set")
		}
		if Clean_Version&Clean_Version == 0 {
			t.Error("Clean_Version should have its own bit set")
		}
		if Clean_Help&Clean_Version != 0 {
			t.Error("Clean_Help and Clean_Version should not overlap")
		}
	})
}

// =============================================================================
// command_line.go parse error path tests (via parseOne directly to avoid os.Exit)
// =============================================================================

func TestParseShortErrors(t *testing.T) {
	t.Run("value flag missing value", func(t *testing.T) {
		cl := NewCommandLine()
		cl.Flag("x", "default", TYPE_STRING, "flag with value")
		cl.args = []string{"-x"}
		_, err := cl.parseOne()
		if err == nil {
			t.Error("expected error when value flag has no value provided")
		}
	})
}

func TestParseLongErrors(t *testing.T) {
	t.Run("non-value flag given value", func(t *testing.T) {
		cl := NewCommandLine()
		cl.FlagNonValueLong("debug", "debug mode")
		cl.args = []string{"--debug=value"}
		_, err := cl.parseOne()
		if err == nil {
			t.Error("expected error when non-value long flag is given a value")
		}
	})

	t.Run("value flag given as non-value", func(t *testing.T) {
		cl := NewCommandLine()
		cl.FlagLong("name", "default", TYPE_STRING, "name")
		cl.args = []string{"--name"}
		_, err := cl.parseOne()
		if err == nil {
			t.Error("expected error when value long flag has no value")
		}
	})

	t.Run("long flag with invalid name format", func(t *testing.T) {
		cl := NewCommandLine()
		cl.args = []string{"--=value"}
		_, err := cl.parseOne()
		if err == nil {
			t.Error("expected error for --=value format")
		}
	})
}

func TestParseOneHasFlag(t *testing.T) {
	t.Run("non-dash arg after flag", func(t *testing.T) {
		cl := NewCommandLine()
		cl.Flag("x", "default", TYPE_STRING, "flag")
		cl.hasFlag = true
		cl.args = []string{"no-dash"}
		_, err := cl.parseOne()
		if err == nil {
			t.Error("expected error when arg after flag doesn't start with -")
		}
	})
}

func TestCommandLineHelp(t *testing.T) {
	t.Run("no args shows top-level help", func(t *testing.T) {
		cl := NewCommandLine()
		cl.help.AppName = "test"
		cl.Help()
	})

	t.Run("with command args", func(t *testing.T) {
		cl := NewCommandLine()
		cl.FlagArg("create", "create command", 0)
		cl.help.AppName = "test"
		cl.Parse([]string{"create"})
		cl.Help()
	})
}

func TestCommandLineHasHelp(t *testing.T) {
	t.Run("no help in others", func(t *testing.T) {
		cl := NewCommandLine()
		if cl.hasHelp() {
			t.Error("hasHelp should return false with no others")
		}
	})

	t.Run("help in others", func(t *testing.T) {
		cl := NewCommandLine()
		cl.others = append(cl.others, &Flag{name: Ko_Command_Help, value: Ko_Command_Help, isArg: true, has: true})
		if !cl.hasHelp() {
			t.Error("hasHelp should return true when help is in others")
		}
	})
}

func TestCommandLineFlagInvalidNames(t *testing.T) {
	t.Run("Flag with invalid short name", func(t *testing.T) {
		cl := NewCommandLine()
		cl.Flag("invalid-name", "default", TYPE_STRING, "bad name")
		if cl.Get("invalid-name") != nil {
			t.Error("flag with invalid name should not be registered")
		}
	})

	t.Run("FlagLong with invalid long name", func(t *testing.T) {
		cl := NewCommandLine()
		cl.FlagLong("", "default", TYPE_STRING, "bad name")
		if cl.Get("") != nil {
			t.Error("flag with empty long name should not be registered")
		}
	})

	t.Run("FlagArg with invalid name", func(t *testing.T) {
		cl := NewCommandLine()
		cl.FlagArg("bad-", "bad name", 0)
		if cl.Get("bad-") != nil {
			t.Error("flag arg with invalid name should not be registered")
		}
	})

	t.Run("FlagNonValue with invalid name", func(t *testing.T) {
		cl := NewCommandLine()
		cl.FlagNonValue("123", "bad name")
		if cl.Get("123") != nil {
			t.Error("flag non-value with invalid name should not be registered")
		}
	})

	t.Run("FlagNonValueLong with invalid name", func(t *testing.T) {
		cl := NewCommandLine()
		cl.FlagNonValueLong("bad-", "bad name")
		if cl.Get("bad-") != nil {
			t.Error("flag non-value long with invalid name should not be registered")
		}
	})
}

func TestCommandLinePrintDefaults(t *testing.T) {
	cl := NewCommandLine()
	cl.Flag("v", "default", TYPE_STRING, "verbose")
	cl.PrintDefaults()
}

// =============================================================================
// help.go tests
// =============================================================================

func TestHelpShow(t *testing.T) {
	h := NewHelp("testapp")
	h.Title = "Test application"
	h.Args.Add("daemon", "run with daemon mode", false, false)
	h.Commands.AddCommand("create", "create something")
	h.Show()
}

func TestHelpCommands(t *testing.T) {
	h := NewHelp("testapp")

	cmd := h.Commands.AddCommand("create", "create module")
	if cmd == nil {
		t.Fatal("expected command to be created")
	}
	if cmd.Name != "create" {
		t.Errorf("command name = %q, want %q", cmd.Name, "create")
	}
	if cmd.Comment != "create module" {
		t.Errorf("command comment = %q, want %q", cmd.Comment, "create module")
	}

	sub := h.Commands.AddSubCommand("create", "config", "create config")
	if sub == nil {
		t.Fatal("expected subcommand to be created")
	}
	if sub.Name != "config" {
		t.Errorf("subcommand name = %q, want %q", sub.Name, "config")
	}

	nilSub := h.Commands.AddSubCommand("nonexistent", "sub", "sub cmd")
	if nilSub != nil {
		t.Error("expected nil for AddSubCommand with non-existent parent")
	}

	h.Commands.AddArg("create", "path", "config path", false, true)
}

func TestHelpArgs(t *testing.T) {
	h := NewHelp("testapp")

	h.Args.Add("path", "config path", false, nil)
	h.Args.Add("daemon", "daemon mode", false, false)

	t.Run("HasRequired", func(t *testing.T) {
		if !h.Args.HasRequired() {
			t.Error("expected HasRequired() to return true")
		}
	})

	t.Run("HasOptions", func(t *testing.T) {
		if !h.Args.HasOptions() {
			t.Error("expected HasOptions() to return true")
		}
	})

	t.Run("HelpTitle includes required", func(t *testing.T) {
		title := h.Args.HelpTitle()
		if title == "" {
			t.Error("expected non-empty HelpTitle")
		}
	})

	t.Run("MaxLen", func(t *testing.T) {
		if h.Args.MaxLen() == 0 {
			t.Error("expected MaxLen > 0 with args registered")
		}
	})

	t.Run("subMaxLen required", func(t *testing.T) {
		if h.Args.subMaxLen(true) == 0 {
			t.Error("expected subMaxLen(true) > 0")
		}
	})

	t.Run("subMaxLen options", func(t *testing.T) {
		if h.Args.subMaxLen(false) == 0 {
			t.Error("expected subMaxLen(false) > 0")
		}
	})
}

func TestArgsFormat(t *testing.T) {
	h := NewHelp("testapp")
	h.Args.Add("daemon", "run with daemon", false, false)
	h.Args.Add("d", "debug mode", true, false)

	formatted := h.Args.Format("    ")
	if formatted == "" {
		t.Error("Format should not return empty string")
	}

	formattedSub := h.Args.FormatSub("    ", false)
	if formattedSub == "" {
		t.Error("FormatSub should not return empty string")
	}

	h2 := NewHelp("testapp")
	h2.Args.Add("path", "config path", false, nil)
	formattedReq := h2.Args.FormatSub("    ", true)
	if formattedReq == "" {
		t.Error("FormatSub for required should not return empty")
	}
}

func TestCommandFormatArg(t *testing.T) {
	h := NewHelp("testapp")
	cmd := h.Commands.AddCommand("create", "create module")
	cmd.AddArg("path", "config path", false, nil)
	cmd.AddArg("type", "config type", false, "json")

	formatted := cmd.FormatArg("    ")
	if formatted == "" {
		t.Error("FormatArg should not return empty string")
	}
}

func TestArgFormat(t *testing.T) {
	t.Run("short arg with comment and default", func(t *testing.T) {
		a := &Arg{Name: "v", Comment: "verbose", IsShort: true, Default: "false"}
		s := a.Format("    ", "  ")
		if !strings.Contains(s, "-v") {
			t.Errorf("short arg format should contain -v, got %q", s)
		}
	})

	t.Run("long arg with only default", func(t *testing.T) {
		a := &Arg{Name: "daemon", Comment: "", IsShort: false, Default: "true"}
		s := a.Format("    ", "  ")
		if !strings.Contains(s, "--daemon") {
			t.Errorf("long arg format should contain --daemon, got %q", s)
		}
	})

	t.Run("long arg no comment no default", func(t *testing.T) {
		a := &Arg{Name: "name", Comment: "", IsShort: false, Default: ""}
		s := a.Format("    ", "  ")
		if !strings.Contains(s, "--name") {
			t.Errorf("long arg format should contain --name, got %q", s)
		}
	})
}

func TestHelpGet(t *testing.T) {
	h := NewHelp("testapp")
	h.Commands.AddCommand("create", "create module")
	h.Commands.AddSubCommand("create", "config", "create config")

	t.Run("get existing command", func(t *testing.T) {
		cmd := h.Get("create")
		if cmd == nil {
			t.Fatal("expected to get create command")
		}
		if cmd.Name != "create" {
			t.Errorf("command name = %q, want %q", cmd.Name, "create")
		}
	})

	t.Run("get subcommand", func(t *testing.T) {
		cmd := h.Get("create", "config")
		if cmd == nil {
			t.Fatal("expected to get config subcommand")
		}
		if cmd.Name != "config" {
			t.Errorf("command name = %q, want %q", cmd.Name, "config")
		}
	})

	t.Run("get non-existent", func(t *testing.T) {
		if cmd := h.Get("nonexistent"); cmd != nil {
			t.Error("expected nil for non-existent command")
		}
	})

	t.Run("get non-existent subcommand", func(t *testing.T) {
		if cmd := h.Get("create", "nonexistent"); cmd != nil {
			t.Error("expected nil for non-existent subcommand")
		}
	})
}

func TestCommandHelp(t *testing.T) {
	h := NewHelp("testapp")
	cmd := h.Commands.AddCommand("create", "create module")
	cmd.AddArg("path", "config path", false, nil)
	cmd.AddArg("type", "config type", false, "json")
	cmd.AddArg("c", "short option", true, false)
	sub := cmd.AddCommand("config", "create config")
	sub.AddArg("o", "open src", true, nil)

	t.Run("Help top-level command", func(t *testing.T) {
		cmd.Help("testapp", "testapp")
	})

	t.Run("Help subcommand", func(t *testing.T) {
		cmd.HelpSub("testapp", "config")
	})

	t.Run("Help with prefix", func(t *testing.T) {
		cmd.Help("testapp create", "testapp")
	})
}

func TestHelpHelp(t *testing.T) {
	h := NewHelp("testapp")
	h.Commands.AddCommand("create", "create module")
	h.Commands.AddSubCommand("create", "config", "create config")

	t.Run("help for top-level command", func(t *testing.T) {
		h.Help("create")
	})

	t.Run("help for nested command", func(t *testing.T) {
		h.Help("create", "config")
	})

	t.Run("help for non-existent command", func(t *testing.T) {
		h.Help("nonexistent")
	})
}

// =============================================================================
// run_time.go tests
// =============================================================================

func TestStartTime(t *testing.T) {
	val := StartTime()
	_ = val
}

func TestGetRunTime(t *testing.T) {
	val := GetRunTime()
	if val < 0 {
		t.Errorf("GetRunTime() = %d, expected >= 0", val)
	}
}

func TestStartTimestamp(t *testing.T) {
	ts := StartTimestamp()
	if ts == "" {
		t.Error("StartTimestamp() should not return empty string")
	}
	if len(ts) != 19 {
		t.Errorf("StartTimestamp() length = %d, want 19 (yyyy-mm-dd HH:MM:SS)", len(ts))
	}
}

func TestGetFormatRunTime(t *testing.T) {
	formatted := GetFormatRunTime()
	if formatted == "" {
		t.Error("GetFormatRunTime() should not return empty string")
	}
	for _, word := range []string{"days", "hours", "minutes", "seconds"} {
		if !strings.Contains(formatted, word) {
			t.Errorf("GetFormatRunTime() = %q, should contain %q", formatted, word)
		}
	}
}

// =============================================================================
// app.go tests
// =============================================================================

func TestNewApp(t *testing.T) {
	app := NewApp("testapp")
	if app == nil {
		t.Fatal("NewApp should not return nil")
	}
}

func TestAppName(t *testing.T) {
	app := NewApp("myapp")
	if app.Name() != "myapp" {
		t.Errorf("Name() = %q, want %q", app.Name(), "myapp")
	}
}

func TestCleanFlagConstants(t *testing.T) {
	if Clean_Defalut != 0 {
		t.Errorf("Clean_Defalut = %d, want 0", Clean_Defalut)
	}
	if Clean_Help != 1 {
		t.Errorf("Clean_Help = %d, want 1", Clean_Help)
	}
	if Clean_Version != 2 {
		t.Errorf("Clean_Version = %d, want 2", Clean_Version)
	}
}

// =============================================================================
// server.go tests (ServBase trivial methods)
// =============================================================================

func TestServBaseVersion(t *testing.T) {
	s := &ServBase{}
	if v := s.Version(); v != "0.0.0" {
		t.Errorf("Version() = %q, want %q", v, "0.0.0")
	}
}

func TestServBaseAuthor(t *testing.T) {
	s := &ServBase{}
	if a := s.Author(); a != "kovey" {
		t.Errorf("Author() = %q, want %q", a, "kovey")
	}
}

func TestServBaseFlag(t *testing.T) {
	s := &ServBase{}
	if err := s.Flag(nil); err != nil {
		t.Errorf("Flag() should return nil, got %v", err)
	}
}

func TestServBaseReload(t *testing.T) {
	s := &ServBase{}
	if err := s.Reload(nil); err != nil {
		t.Errorf("Reload() should return nil, got %v", err)
	}
}

func TestServBasePanic(t *testing.T) {
	s := &ServBase{}
	// Panic() expects a valid AppInterface; pass a real App
	app := NewApp("test-panic")
	s.Panic(app)
}

func TestServBaseUsage(t *testing.T) {
	s := &ServBase{}
	// Usage calls PrintDefaults which uses global _commanLine
	// Should not panic
	s.Usage()
}

// =============================================================================
// global functions test
// =============================================================================

func TestGlobalPrintDefaults(t *testing.T) {
	// Should not panic
	PrintDefaults()
}

func TestGlobalGetHelp(t *testing.T) {
	h := GetHelp()
	if h == nil {
		t.Error("GetHelp() should not return nil")
	}
}

// =============================================================================
// additional parseOne coverage tests
// =============================================================================

// TestParseOneHelpCommand tests parseOne when parsing the help command
func TestParseOneHelpCommand(t *testing.T) {
	cl := NewCommandLine()
	// Register help as a command arg
	cl.FlagArg(Ko_Command_Help, "show help", 0)
	cl.args = []string{Ko_Command_Help}
	end, err := cl.parseOne()
	if err != nil {
		t.Errorf("expected no error parsing help command, got: %v", err)
	}
	if !end {
		t.Error("expected end=true after parsing help command")
	}
	if !cl.hasHelp() {
		t.Error("expected hasHelp to return true after parsing help command")
	}
}

// TestParseOneHasFlagLongDash tests parseOne when hasFlag is true and arg has long dash
func TestParseOneHasFlagLongDash(t *testing.T) {
	cl := NewCommandLine()
	cl.Flag("x", "default", TYPE_STRING, "flag")
	cl.hasFlag = true
	cl.args = []string{"---bad"}
	_, err := cl.parseOne()
	if err == nil {
		t.Error("expected error for long dash arg after flag")
	}
}

// TestParseOneFlagArgShortDash tests parseOne when hasFlag is true with short arg after flag
func TestParseOneHasFlagShortArg(t *testing.T) {
	cl := NewCommandLine()
	cl.Flag("x", "default", TYPE_STRING, "flag")
	cl.FlagNonValue("y", "another flag")
	cl.hasFlag = true
	cl.args = []string{"-y"}
	end, err := cl.parseOne()
	if err != nil {
		t.Errorf("expected no error for valid short flag after flag, got: %v", err)
	}
	if !end {
		t.Error("expected end=true")
	}
	flagY := cl.Get("y")
	if flagY == nil || !flagY.IsInput() {
		t.Error("expected flag y to be set")
	}
}

// TestParseShortWithNonValueAndValueProvided tests error when value provided to non-value flag
func TestParseShortWithNonValueAndValueProvided(t *testing.T) {
	cl := NewCommandLine()
	cl.FlagNonValue("x", "non-value flag")
	cl.args = []string{"-x", "somevalue"}
	_, err := cl.parseOne()
	if err == nil {
		t.Error("expected error when providing value to non-value short flag")
	}
}

// TestParseLongNonValueFlagGivenValue tests long non-value flag given value
func TestParseLongNonValueFlagGivenValue(t *testing.T) {
	cl := NewCommandLine()
	cl.FlagNonValueLong("quiet", "quiet mode")
	cl.args = []string{"--quiet=yes"}
	_, err := cl.parseOne()
	if err == nil {
		t.Error("expected error when non-value long flag is given a value via =")
	}
}

// TestParseLongNoEqualButCheckLongFails tests parseLong when checkLong fails
func TestParseLongNoEqualCheckLongFails(t *testing.T) {
	cl := NewCommandLine()
	cl.args = []string{"--123"}
	_, err := cl.parseOne()
	if err == nil {
		t.Error("expected error for --123 (checkLong fails)")
	}
}

// TestPrintDefaultsOnInstance tests the instance-level PrintDefaults
func TestCommandLineInstancePrintDefaults(t *testing.T) {
	cl := NewCommandLine()
	cl.help.AppName = "test"
	cl.help.Title = "Test CLI"
	cl.Flag("x", "", TYPE_STRING, "a flag")
	cl.PrintDefaults()
}

// TestCommandLineFlagHierarchical tests registering flags with parent hierarchy
func TestCommandLineFlagHierarchical(t *testing.T) {
	cl := NewCommandLine()
	cl.FlagArg("create", "create command", 0)
	cl.Flag("x", "default", TYPE_STRING, "flag under create", "create")

	flag := cl.Get("create", "x")
	if flag == nil {
		t.Fatal("expected hierarchical flag create->x to exist")
	}

	// Verify the flag works
	cl.Parse([]string{"create", "-x", "hello"})
	flagParsed := cl.Get("create", "x")
	if flagParsed == nil || !flagParsed.IsInput() {
		t.Error("expected hierarchical flag to be set after parsing")
	}
	if flagParsed.String() != "hello" {
		t.Errorf("hierarchical flag value = %q, want %q", flagParsed.String(), "hello")
	}
}

// TestArgNameLen tests Arg.NameLen for both short and long args
func TestArgNameLen(t *testing.T) {
	t.Run("short arg name len", func(t *testing.T) {
		a := &Arg{Name: "v", IsShort: true}
		if n := a.NameLen(); n != 2 {
			t.Errorf("NameLen() for short arg = %d, want 2 (name + 1 for -)", n)
		}
	})

	t.Run("long arg name len", func(t *testing.T) {
		a := &Arg{Name: "verbose", IsShort: false}
		if n := a.NameLen(); n != 9 {
			t.Errorf("NameLen() for long arg = %d, want 9 (name + 2 for --)", n)
		}
	})
}

// TestCommandsMaxLen tests Commands.MaxLen
func TestCommandsMaxLen(t *testing.T) {
	cmds := NewCommands()
	cmds.AddCommand("short", "short cmd")
	cmds.AddCommand("verylongcommand", "long cmd")

	maxLen := cmds.MaxLen()
	if maxLen != 15 {
		t.Errorf("MaxLen() = %d, want 15", maxLen)
	}
}

// TestArgsHasRequiredWithNoRequired tests HasRequired returns false when no required args
func TestArgsHasRequiredFalse(t *testing.T) {
	a := NewArgs()
	a.Add("opt", "optional", false, "default")
	if a.HasRequired() {
		t.Error("HasRequired() should return false when all args are optional")
	}
	if !a.HasOptions() {
		t.Error("HasOptions() should return true when args have defaults")
	}
}

// TestArgsHasOptionsFalse tests HasOptions returns false when all args are required
func TestArgsHasOptionsFalse(t *testing.T) {
	a := NewArgs()
	a.Add("path", "required path", false, nil)
	if !a.HasRequired() {
		t.Error("HasRequired() should return true when arg has nil default")
	}
	if a.HasOptions() {
		t.Error("HasOptions() should return false when all args are required")
	}
}

// TestFlagHasValueFalsePanics tests that accessing a non-value flag's value panics
func TestFlagHasValueFalsePanics(t *testing.T) {
	t.Run("String on non-value flag", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when getting value of non-value flag")
			}
		}()
		f := &Flag{name: "test", t: TYPE_STRING, has: true, hasValue: false}
		_ = f.String()
	})

	t.Run("Int on non-value flag", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when getting int value of non-value flag")
			}
		}()
		f := &Flag{name: "test", t: TYPE_INT, has: true, hasValue: false}
		f.Int()
	})
}

// TestArgHelpName tests Arg.HelpName for both short and long
func TestArgHelpName(t *testing.T) {
	t.Run("short help name", func(t *testing.T) {
		a := &Arg{Name: "v", IsShort: true}
		if n := a.HelpName(); n != "-v" {
			t.Errorf("HelpName() = %q, want %q", n, "-v")
		}
	})

	t.Run("long help name", func(t *testing.T) {
		a := &Arg{Name: "verbose", IsShort: false}
		if n := a.HelpName(); n != "--verbose" {
			t.Errorf("HelpName() = %q, want %q", n, "--verbose")
		}
	})
}

// TestCommandFormat tests Command.Format
func TestCommandFormat(t *testing.T) {
	cmd := &Command{Name: "create", Comment: "create a module"}
	formatted := cmd.Format("    ", "  ")
	if formatted == "" {
		t.Error("Format should not return empty string")
	}
	if !strings.Contains(formatted, "create") || !strings.Contains(formatted, "create a module") {
		t.Errorf("Format() = %q, should contain name and comment", formatted)
	}
}

// TestCommandsFormat tests Commands.Format
func TestCommandsFormat(t *testing.T) {
	cmds := NewCommands()
	cmds.AddCommand("create", "create module")
	cmds.AddCommand("delete", "delete module")

	formatted := cmds.Format(cmds.MaxLen())
	if formatted == "" {
		t.Error("Format should not return empty string")
	}
}

// TestHelpHelpWithNestedPrefix tests help with multi-level nested prefix
func TestHelpHelpWithNestedPrefix(t *testing.T) {
	h := NewHelp("testapp")
	h.Commands.AddCommand("svc", "service commands")
	h.Commands.AddSubCommand("svc", "create", "create service")

	// help "svc" "create" -> nested path
	h.Help("svc", "create")
}

// =============================================================================
// daemon.go safe delegate method tests
// =============================================================================

func TestAppDaemonPid(t *testing.T) {
	app := NewApp("test-pid")
	pid := app.Pid()
	// Pid is set by os.Getpid() during Run(), not during NewApp.
	// For NewApp, it should be 0 (not initialized).
	_ = pid
}

func TestAppDaemonPidString(t *testing.T) {
	app := NewApp("test-pidstr")
	ps := app.PidString()
	if ps == "" {
		t.Error("PidString should not be empty")
	}
}

func TestAppDaemonSetDebugLevel(t *testing.T) {
	app := NewApp("test-debug")
	// Should not panic
	app.SetDebugLevel("4")
}

func TestAppDaemonContext(t *testing.T) {
	app := NewApp("test-ctx")
	ctx := app.Context()
	if ctx == nil {
		t.Error("Context should not return nil")
	}
}

func TestAppDaemonUsageWhenErr(t *testing.T) {
	app := NewApp("test-usage")
	// Should not panic
	app.UsageWhenErr()
}

func TestAppDaemonCleanCommandLine(t *testing.T) {
	app := NewApp("test-clean")
	// Should not panic
	app.CleanCommandLine(true)
}

func TestAppDaemonCleanCommandLineWith(t *testing.T) {
	app := NewApp("test-cleanwith")
	// Clean with Clean_Version flag
	app.CleanCommandLineWith(Clean_Version)
	// Clean with Clean_Help flag
	app.CleanCommandLineWith(Clean_Help)
	// Clean with both
	app.CleanCommandLineWith(Clean_Help | Clean_Version)
}

func TestAppDaemonGet(t *testing.T) {
	app := NewApp("test-get")
	// Get returns flag from global _commanLine
	f, err := app.Get("start")
	if err != nil {
		t.Errorf("Get(start) should succeed, got: %v", err)
	}
	if f == nil {
		t.Error("Get(start) should return a flag")
	}
}

func TestAppDaemonFlagMethods(t *testing.T) {
	app := NewApp("test-flags")

	// These methods delegate to _commanLine
	app.Flag("z", "default", TYPE_STRING, "test short flag")
	app.FlagLong("testlong", "default", TYPE_STRING, "test long flag")
	app.FlagNonValue("y", "test non-value")
	app.FlagNonValueLong("testnonvalue", "test non-value long")
	app.FlagArg("testarg", "test command arg")

	// Verify they were registered on global _commanLine
	f1, _ := app.Get("z")
	if f1 == nil {
		t.Error("Flag z should be registered")
	}
	f2, _ := app.Get("testlong")
	if f2 == nil {
		t.Error("FlagLong testlong should be registered")
	}
	f3, _ := app.Get("y")
	if f3 == nil {
		t.Error("FlagNonValue y should be registered")
	}
	f4, _ := app.Get("testnonvalue")
	if f4 == nil {
		t.Error("FlagNonValueLong testnonvalue should be registered")
	}
	f5, _ := app.Get("testarg")
	if f5 == nil {
		t.Error("FlagArg testarg should be registered")
	}
}

func TestAppDaemonArg(t *testing.T) {
	app := NewApp("test-daemon-arg")
	// No args parsed yet, so Arg(0) should return error
	_, err := app.Arg(0, TYPE_STRING)
	if err == nil {
		t.Error("Arg(0) should return error when no args parsed")
	}
}

func TestAppDaemonFlagNameTooShort(t *testing.T) {
	app := NewApp("test-short")
	// Daemon.FlagArg with name < 2 characters should log warning and not register
	app.FlagArg("x", "too short")
	_, err := app.Get("x")
	// Should not have registered since name is too short
	if err == nil {
		t.Error("FlagArg with short name should not be registered")
	}
}

func TestAppDaemonFlagLongNameTooShort(t *testing.T) {
	app := NewApp("test-longshort")
	// Daemon.FlagLong with name < 2 characters
	app.FlagLong("x", "default", TYPE_STRING, "too short")
	_, err := app.Get("x")
	if err == nil {
		t.Error("FlagLong with short name should not be registered")
	}
}

func TestAppDaemonFlagNonValueLongTooShort(t *testing.T) {
	app := NewApp("test-nvlong-short")
	// Daemon.FlagNonValueLong with name < 2 characters
	app.FlagNonValueLong("x", "too short")
	_, err := app.Get("x")
	if err == nil {
		t.Error("FlagNonValueLong with short name should not be registered")
	}
}

func TestAppDaemonFlagArgTooShort(t *testing.T) {
	app := NewApp("test-argshort")
	// Daemon.FlagArg with name < 2 characters
	app.FlagArg("x", "too short")
	_, err := app.Get("x")
	if err == nil {
		t.Error("FlagArg with short name should not be registered")
	}
}

func TestAppDaemonTerm(t *testing.T) {
	app := NewApp("test-term")
	// term() cancels the context
	app.term()
	ctx := app.Context()
	if ctx.Err() == nil {
		t.Error("context should be cancelled after term()")
	}
}

func TestCommandLineFlagNonExistentParent(t *testing.T) {
	cl := NewCommandLine()
	// Flag under non-existent parent should log warning but not panic
	cl.Flag("x", "default", TYPE_STRING, "flag", "nonexistent")
	// Should not be registered since parent doesn't exist
	if cl.Get("nonexistent", "x") != nil {
		t.Error("flag should not be registered under non-existent parent")
	}
}

func TestCommandLineFlagDuplicateChild(t *testing.T) {
	cl := NewCommandLine()
	cl.FlagArg("create", "create command", 0)
	// Register first child
	cl.Flag("x", "first", TYPE_STRING, "first flag", "create")
	// Register duplicate child with same name - should log warning, not overwrite
	cl.Flag("x", "second", TYPE_STRING, "second flag", "create")

	flag := cl.Get("create", "x")
	if flag == nil {
		t.Fatal("expected flag to exist")
	}
	// Should still be the first one
	if flag.def != "first" {
		t.Errorf("duplicate child should not overwrite, got def=%v", flag.def)
	}
}

func TestSkipHelpInGetPathCount(t *testing.T) {
	flags := NewFlags()
	parent := &Flag{name: "create"}
	child := &Flag{name: "child", t: TYPE_STRING, hasValue: true, def: "val"}
	parent.AddChild(child)
	flags.Add(parent)

	// Multi-element path without help: should find child normally
	got := flags.Get("create", "child")
	if got == nil {
		t.Error("Get(create, child) should find the child")
	}

	// Multi-element path with help in middle (count > 1): help should be skipped
	got2 := flags.Get("create", Ko_Command_Help, "child")
	// After skipping help, looks for child -> should find it
	if got2 == nil {
		t.Error("Get(create, help, child) should skip help and find child")
	}
}

// =============================================================================
// mock ServInterface for testing daemon methods that require a service
// =============================================================================

type testServ struct {
	ServBase
}

func (t *testServ) Init(AppInterface) error          { return nil }
func (t *testServ) Run(AppInterface) error           { return nil }
func (t *testServ) AsyncLogClose(AppInterface)       {} // override to prevent async.Close panic

func TestAppDaemonSetServ(t *testing.T) {
	app := NewApp("test-setserv")
	s := &testServ{}
	app.SetServ(s)
}

func TestAppDaemonGetPidFileNotFound(t *testing.T) {
	app := NewApp("test-getpid")
	s := &testServ{}
	app.SetServ(s)
	pid := app.getPid()
	if pid != -1 {
		t.Errorf("getPid() with no pid file = %d, want -1", pid)
	}
}

func TestAppDaemonGetPidAndChildPidNotFound(t *testing.T) {
	app := NewApp("test-getpidchild")
	s := &testServ{}
	app.SetServ(s)
	pids := app.getPidAndChildPid()
	if pids != nil {
		t.Errorf("getPidAndChildPid() with no pid file = %v, want nil", pids)
	}
}

func TestAppDaemonCleanCommandLineBothPaths(t *testing.T) {
	app := NewApp("test-clean2")
	app.CleanCommandLine(true)
	app.CleanCommandLine(false)
}

func TestServBaseAsyncLogEarlyReturn(t *testing.T) {
	s := &ServBase{}
	s.AsyncLog(nil)
}

func TestServBasePidFileDefault(t *testing.T) {
	s := &ServBase{}
	app := NewApp("test-pidfile")
	path := s.PidFile(app)
	if path == "" {
		t.Error("PidFile should not return empty string")
	}
	if !strings.Contains(path, app.Name()) {
		t.Errorf("PidFile = %q, should contain app name %q", path, app.Name())
	}
}

func TestAppDaemonShowVersion(t *testing.T) {
	app := NewApp("test-showversion")
	s := &testServ{}
	app.SetServ(s)
	// _showVersion should not panic with a valid serv set
	app._showVersion()
}

func TestAppDaemonShowVersionNilServ(t *testing.T) {
	app := NewApp("test-showversion-nil")
	// _showVersion should return early when serv is nil
	app._showVersion()
}

func TestAppDaemonRunCommandHelp(t *testing.T) {
	app := NewApp("test-cmd-help")
	// _runCommand with Ko_Command_Help should call _commanLine.Help()
	err := app._runCommand(Ko_Command_Help)
	if err != nil {
		t.Errorf("_runCommand(help) should not error, got: %v", err)
	}
}

func TestAppDaemonRunCommandReload(t *testing.T) {
	app := NewApp("test-cmd-reload")
	s := &testServ{}
	app.SetServ(s)
	// _runCommand with Ko_Command_Reload calls _reload()
	// getPid returns -1 (no pid file), so _reload returns early
	err := app._runCommand(Ko_Command_Reload)
	if err != nil {
		t.Errorf("_runCommand(reload) should not error, got: %v", err)
	}
}

func TestAppDaemonRunCommandStop(t *testing.T) {
	app := NewApp("test-cmd-stop")
	s := &testServ{}
	app.SetServ(s)
	// _runCommand with Ko_Command_Stop calls _stop()
	// getPid returns -1 (no pid file), so _stop returns early
	err := app._runCommand(Ko_Command_Stop)
	if err != nil {
		t.Errorf("_runCommand(stop) should not error, got: %v", err)
	}
}

func TestAppDaemonRunCommandKill(t *testing.T) {
	app := NewApp("test-cmd-kill")
	s := &testServ{}
	app.SetServ(s)
	// _runCommand with Ko_Command_Kill calls _kill()
	// getPidAndChildPid returns nil (no pid file), so _kill returns early
	err := app._runCommand(Ko_Command_Kill)
	if err != nil {
		t.Errorf("_runCommand(kill) should not error, got: %v", err)
	}
}

func TestAppDaemonReloadNoPid(t *testing.T) {
	app := NewApp("test-reload-nopid")
	s := &testServ{}
	app.SetServ(s)
	// _reload with no pid file should return early
	err := app._reload()
	if err != nil {
		t.Errorf("_reload() should not error, got: %v", err)
	}
}

func TestAppDaemonStopNoPid(t *testing.T) {
	app := NewApp("test-stop-nopid")
	s := &testServ{}
	app.SetServ(s)
	// _stop with no pid file should return early
	err := app._stop()
	if err != nil {
		t.Errorf("_stop() should not error, got: %v", err)
	}
}

func TestAppDaemonKillNoPid(t *testing.T) {
	app := NewApp("test-kill-nopid")
	s := &testServ{}
	app.SetServ(s)
	// _kill with no pid file should return early
	err := app._kill()
	if err != nil {
		t.Errorf("_kill() should not error, got: %v", err)
	}
}

func TestAppDaemonRunCommandStart(t *testing.T) {
	app := NewApp("test-cmd-start")
	s := &testServ{}
	app.SetServ(s)
	// _runCommand with Ko_Command_Start calls _run("start")
	// With isBackground=false, listen() returns immediately
	// _runApp runs in goroutine with mock serv
	err := app._runCommand(Ko_Command_Start)
	if err != nil {
		t.Errorf("_runCommand(start) should not error, got: %v", err)
	}
}

func TestAppDaemonRunCommandRestart(t *testing.T) {
	app := NewApp("test-cmd-restart")
	s := &testServ{}
	app.SetServ(s)
	// _runCommand with Ko_Command_Restart calls _restart()
	// _restart calls _stop() (safe) then _run(Ko_Command_Restart)
	// _run calls runApp with mock serv
	err := app._runCommand(Ko_Command_Restart)
	if err != nil {
		t.Errorf("_runCommand(restart) should not error, got: %v", err)
	}
}