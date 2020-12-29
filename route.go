// Package route contains a command router for Snart.
package route

import (
	"errors"
	"log"
	"strings"

	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/arikawa/state"
	ff "github.com/itzg/go-flagsfiller"

	"github.com/go-snart/db"
)

// RouteCat is the category name for internal Commands.
const RouteCat = "route"

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

	Fill *ff.FlagSetFiller
	Cats map[string][]*Command
}

// New makes an empty Route from the given DB and Session.
func New(d *db.DB, s *state.State) *Route {
	r := &Route{
		DB:    d,
		State: s,

		Fill: ff.New(),
		Cats: make(map[string][]*Command),
	}

	r.Add(RouteCat,
		r.HelpCommand(),
	)

	return r
}

// Add adds Commands to the Route.
func (r *Route) Add(cat string, cmds ...*Command) {
	for _, c := range cmds {
		c.Tidy()
	}

	r.Cats[cat] = append(r.Cats[cat], cmds...)
}

// Handle returns a MessageCreate handler function for the Route.
func (r *Route) Handle(m *gateway.MessageCreateEvent) {
	if m.Message.Author.ID == r.State.Ready.User.ID || m.Message.Author.Bot {
		return
	}

	lines := strings.Split(m.Message.Content, "\n")

	for _, line := range lines {
		pfx := r.LinePrefix(m.GuildID, line)
		if pfx == nil {
			continue
		}

		t, err := r.Trigger(pfx, m.Message, line)
		if err != nil {
			log.Printf("get trigger: %s", err)

			continue
		}

		err = t.Run()
		if err != nil {
			log.Printf("run trigger: %s", err)

			continue
		}
	}
}
