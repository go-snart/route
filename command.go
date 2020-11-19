package route

// Command is a command.
type Command struct {
	Name        string
	Category    string
	Description string

	// Func is the actual logic of the Command.
	Func func(*Trigger) error

	// Hidden is simply a flag for tidying up help menus;
	// secure commands should still check for permissions.
	Hidden bool

	// Flags should be a struct value of the type used to parse flags.
	// See docs for github.com/itzg/go-flagsfiller.
	Flags interface{}
}
