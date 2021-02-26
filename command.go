package route

// DefaultDesc is the default description used by Tidy.
const DefaultDesc = "*no description*"

// UndefinedMsg is the message content sent by Undefined.
const UndefinedMsg = "the behavior of this command is not yet defined."

// Func is a handler for a Trigger.
type Func = func(*Trigger) error

// Command is a command.
type Command struct {
	Name  string
	Desc  string
	Cat   string
	Func  Func
	Hide  bool
	Flags interface{}
}

// Tidy fixes common issues with Command contents.
func (c *Command) Tidy() {
	if c.Desc == "" {
		c.Desc = DefaultDesc
	}

	if c.Func == nil {
		c.Func = Undefined
	}

	if c.Flags == nil {
		c.Flags = struct{}{}
	}
}

// Undefined is a Func used when a Command's Func isn't defined.
func Undefined(t *Trigger) error {
	rep := t.Reply()
	rep.Content = UndefinedMsg

	return rep.Send()
}
