// Package route contains a command Route for Snart.
package route

import (
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state"
	ff "github.com/itzg/go-flagsfiller"

	"github.com/go-snart/lob"
)

// Route handles storing and looking up Commands.
type Route struct {
	State    *state.State
	Commands map[string]Command
	Guilds   map[discord.GuildID]Settings
	Filler   *ff.FlagSetFiller
}

// New makes an empty Route from the given Config and Session.
func New(base Settings, s *state.State) *Route {
	r := &Route{
		State:    s,
		Commands: map[string]Command{},
		Guilds:   map[discord.GuildID]Settings{discord.NullGuildID: base},
		Filler:   ff.New(),
	}

	r.Add(r.HelpCommand())

	return r
}

// Add adds Commands to the Route.
func (r *Route) Add(cmds ...Command) {
	for _, c := range cmds {
		c.Tidy()
		r.Commands[c.Name] = c
	}
}

// Handle returns a MessageCreate handler function for the Route.
func (r *Route) Handle(m *gateway.MessageCreateEvent) {
	if m.Author.Bot {
		return
	}

	me, err := r.State.Me()
	if err != nil {
		_ = lob.Std.Error("get me: %w", err)

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
			_ = lob.Std.Error("handle line %q: %w", line, err)
		}
	}
}

func (r *Route) handleLine(m *gateway.MessageCreateEvent, line string, me discord.User, mme *discord.Member) error {
	pfx := r.LinePrefix(m.GuildID, me, mme, line)
	if pfx == nil {
		return lob.Std.Error("no prefix")
	}

	t, err := r.Trigger(pfx, m.Message, line)
	if err != nil {
		return lob.Std.Error("get trigger: %w", err)
	}

	err = t.Run()
	if err != nil {
		return lob.Std.Error("run trigger: %w", err)
	}

	return nil
}
