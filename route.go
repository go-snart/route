// Package route contains a command router for Snart.
package route

import (
	"fmt"
	"log"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state"
	ff "github.com/itzg/go-flagsfiller"

	"github.com/go-snart/db"
)

// RouteCat is the category name for internal Commands.
const RouteCat = "route"

// Route handles storing and looking up routes.
type Route struct {
	DB    *db.DB
	State *state.State
	Fill  *ff.FlagSetFiller
	Cats  map[string][]*Command
	Me    *discord.User
}

// New makes an empty Route from the given DB and Session.
func New(d *db.DB, s *state.State) *Route {
	r := &Route{
		DB:    d,
		State: s,
		Fill:  ff.New(),
		Cats:  make(map[string][]*Command),
		Me:    nil,
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

// GetMe wraps the Route's Me field, pulling from State if needed.
func (r *Route) GetMe() (*discord.User, error) {
	if r.Me != nil {
		return r.Me, nil
	}

	me, err := r.State.Me()
	if err != nil {
		return nil, fmt.Errorf("state me: %w", err)
	}

	r.Me = me

	return me, nil
}

// Handle returns a MessageCreate handler function for the Route.
func (r *Route) Handle(m *gateway.MessageCreateEvent) {
	me, err := r.GetMe()
	if err != nil {
		log.Printf("get me: %s", err)

		return
	}

	if m.Message.Author.ID == me.ID || m.Message.Author.Bot {
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
