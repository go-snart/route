// Package route contains a command Route for Snart.
package route

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state"
)

// ErrNoLinePrefix occurs when a line doesn't start with a valid prefix.
var ErrNoLinePrefix = errors.New("no prefix in line")

// Route handles storing and looking up Commands.
type Route struct {
	State *state.State

	cmdMu  *sync.RWMutex
	cmdMap map[string]Cmd

	setMu  *sync.RWMutex
	setMap map[discord.GuildID]Settings
}

// New makes an empty Route from the given Config and Session.
func New(base Settings, s *state.State) *Route {
	r := &Route{
		State: s,

		cmdMu:  &sync.RWMutex{},
		cmdMap: map[string]Cmd{},

		setMu: &sync.RWMutex{},
		setMap: map[discord.GuildID]Settings{
			BaseID: base,
		},
	}

	r.AddCmd(r.HelpCommand())

	return r
}

// Handle is a MessageCreate handler function for the Route.
func (r *Route) Handle(m *gateway.MessageCreateEvent) {
	if m.Author.Bot {
		return
	}

	me, err := r.State.Me()
	if err != nil {
		log.Printf("error: get me: %s", err)

		return
	}

	if m.Author.ID == me.ID {
		return
	}

	mme, _ := r.State.Member(m.GuildID, me.ID)

	lines := strings.Split(m.Message.Content, "\n")

	for _, line := range lines {
		err := r.handleLine(m, line, *me, mme)
		if err != nil {
			log.Printf("error: handle line %q: %s", line, err)
		}
	}
}

func (r *Route) handleLine(m *gateway.MessageCreateEvent, line string, me discord.User, mme *discord.Member) error {
	pfx := r.LinePrefix(m.GuildID, me, mme, line)
	if pfx == nil {
		return ErrNoLinePrefix
	}

	t, err := r.Trigger(*pfx, m.Message, line)
	if err != nil {
		return fmt.Errorf("get trigger: %w", err)
	}

	err = t.Command.Func(t)
	if err != nil {
		return fmt.Errorf("run trigger: %w", err)
	}

	return nil
}
