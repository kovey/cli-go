// Example: Basic CLI Application
//
// This example demonstrates the basic usage of the cli-go framework:
//   - Creating an app with NewApp
//   - Registering flags (short and long forms)
//   - Registering subcommands with hierarchical flags
//   - Accessing parsed flag values
//
// Usage (subcommands must come before flags):
//
//	# Root-level mode (no subcommand)
//	go run main.go --name=myapp --count=5 --verbose
//
//	# Show help
//	go run main.go --help
//
//	# Subcommand with its own flags
//	go run main.go create --path=/tmp/config --type=json
//
//	# Nested subcommand
//	go run main.go create build --to-user=admin
//
//	# Subcommand-specific help
//	go run main.go create help
package main

import (
	"fmt"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/debug-go/debug"
)

// MyService implements the ServInterface for handling the application lifecycle.
type MyService struct {
	*app.ServBase
}

// Flag registers custom command-line flags.
// This is called before argument parsing.
func (s *MyService) Flag(a app.AppInterface) error {
	// -- Top-level flags (used when no subcommand is specified) --
	a.FlagLong("name", "default-app", app.TYPE_STRING, "application name")
	a.FlagLong("count", 1, app.TYPE_INT, "worker count")
	a.FlagNonValueLong("verbose", "enable verbose logging")

	// -- Subcommand: create --
	a.FlagArg("create", "create a new resource")

	// Flags belonging to the "create" subcommand
	a.FlagLong("path", "/tmp/default", app.TYPE_STRING, "output path", "create")
	a.FlagLong("type", "json", app.TYPE_STRING, "config type: json|yaml|toml", "create")

	// -- Nested subcommand: create build --
	a.FlagArg("build", "build the created resource", "create")

	// Flags belonging to "create build"
	a.FlagLong("to-user", "admin", app.TYPE_STRING, "target user for build", "create", "build")

	return nil
}

// Init is called after flag parsing, before Run.
func (s *MyService) Init(a app.AppInterface) error {
	debug.Info("[%s] initializing...", a.Name())
	return nil
}

// Run is the main entry point after initialization.
func (s *MyService) Run(a app.AppInterface) error {
	debug.Info("[%s] running...", a.Name())

	// Access top-level flags (root mode)
	if name, err := a.Get("name"); err == nil {
		fmt.Printf("App name: %s\n", name.String())
	}

	if count, err := a.Get("count"); err == nil {
		fmt.Printf("Worker count: %d\n", count.Int())
	}

	if verbose, err := a.Get("verbose"); err == nil && verbose.IsInput() {
		fmt.Println("Verbose mode: enabled")
	}

	// Access positional arguments (subcommands)
	if arg, err := a.Arg(0, app.TYPE_STRING); err == nil {
		fmt.Printf("Command: %s\n", arg.String())
	}
	if arg, err := a.Arg(1, app.TYPE_STRING); err == nil {
		fmt.Printf("SubCommand: %s\n", arg.String())
	}

	// Access hierarchical flags under subcommands
	if path, err := a.Get("create", "path"); err == nil {
		fmt.Printf("Create path: %s\n", path.String())
	}
	if cfgType, err := a.Get("create", "type"); err == nil {
		fmt.Printf("Config type: %s\n", cfgType.String())
	}

	// Access nested subcommand flags
	if toUser, err := a.Get("create", "build", "to-user"); err == nil {
		fmt.Printf("Build target user: %s\n", toUser.String())
	}

	return nil
}

func main() {
	// Create the application
	cli := app.NewApp("myapp")

	// Configure debug output
	cli.SetDebugLevel(debug.Debug_Info)
	debug.SetFileLine(debug.File_Line_Show)

	// Attach the service handler
	cli.SetServ(&MyService{})

	// Run the application (parses args, calls Flag→Init→Run)
	if err := cli.Run(); err != nil {
		debug.Erro("application error: %s", err)
	}
}
