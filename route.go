// Package route contains a command router for Snart.
package route

import (
	"errors"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/arikawa/state"
	ff "github.com/itzg/go-flagsfiller"

	"github.com/go-snart/db"
	"github.com/go-snart/logs"
)

var (
	// ErrNoCmd occurs when no command is given after a prefix.
	ErrNoCmd = errors.New("no cmd")

	// ErrNoTrigger occurs when no suitable Command is found to create a Trigger.
	ErrNoTrigger = errors.New("no ctx found")
)

// Route handles storing and looking up routes.
type Route struct {
	*db.DB
	*state.State

	Filler   *ff.FlagSetFiller
	Commands []*Command
}

// New makes an empty Route from the given DB and Session.
func New(d *db.DB, s *state.State) *Route {
	return &Route{
		DB:    d,
		State: s,

		Filler:   ff.New(),
		Commands: nil,
	}
}

// Add adds Commands to the Route.
func (r *Route) Add(cmds ...*Command) {
	r.Commands = append(r.Commands, cmds...)
}

// Handle returns a MessageCreate handler function for the Route.
func (r *Route) Handle(m *gateway.MessageCreateEvent) {
	logs.Debug.Println("handling")

	if m.Message.Author.ID == r.State.Ready.User.ID {
		logs.Debug.Println("ignore self")

		return
	}

	if m.Message.Author.Bot {
		logs.Debug.Println("ignore bot")

		return
	}

	lines := strings.Split(m.Message.Content, "\n")
	logs.Debug.Printf("lines %#v", lines)

	for _, line := range lines {
		logs.Debug.Printf("line %q", line)

		pfx := r.LinePrefix(m.GuildID, line)
		if pfx == nil {
			continue
		}

		t, err := r.Trigger(pfx, m.Message, line)
		if err != nil {
			err = fmt.Errorf("get ctx: %w", err)
			logs.Warn.Println(err)

			continue
		}

		err = t.Run()
		if err != nil {
			err = fmt.Errorf("t run: %w", err)
			logs.Warn.Println(err)

			continue
		}
	}
}
