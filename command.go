package route

import (
	"log"
)

// Cmd is a command.
type Cmd struct {
	Name  string
	Desc  string
	Cat   string
	Func  Func
	Hide  bool
	Flags interface{}
}

// Add adds Cmds to the Route.
// Duplicating names is not allowed.
func (r *Route) Add(cs ...Cmd) {
	r.Lock()
	defer r.Unlock()

	for _, cmd := range cs {
		if _, ok := r.Cmds[cmd.Name]; ok {
			log.Printf("cmd %q already exists (skipping)", cmd.Name)

			continue
		}

		r.Cmds[cmd.Name] = cmd
	}
}
