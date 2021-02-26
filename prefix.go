package route

import (
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
)

// Prefix is a command prefix.
type Prefix struct {
	Value string
	Clean string
}

// GuildPrefix returns the Prefix for the given Guild.
func (r *Route) GuildPrefix(g discord.GuildID) *Prefix {
	set, ok := r.Guilds[g]
	if !ok {
		return nil
	}

	return &Prefix{
		Value: set.Prefix,
		Clean: set.Prefix,
	}
}

// DefaultPrefix returns the default Prefix.
// This is the Prefix for discord.NullGuildID.
func (r *Route) DefaultPrefix() *Prefix {
	return r.GuildPrefix(discord.NullGuildID)
}

// MemberPrefix attempts to make a Prefix from a Member.
func (r *Route) MemberPrefix(mme *discord.Member) *Prefix {
	if mme == nil {
		return nil
	}

	if mme.Nick != "" {
		return &Prefix{
			Value: mme.Mention(),
			Clean: "@" + mme.Nick,
		}
	}

	return &Prefix{
		Value: mme.Mention(),
		Clean: "@" + mme.User.Username,
	}
}

// UserPrefix makes a Prefix from a User.
func (r *Route) UserPrefix(me discord.User) *Prefix {
	return &Prefix{
		Value: me.Mention(),
		Clean: "@" + me.Username,
	}
}

func (r *Route) findPrefix(
	g discord.GuildID,
	mme *discord.Member,
	me discord.User,
	fn func(*Prefix) bool,
) *Prefix {
	pfx := r.GuildPrefix(g)
	if pfx != nil && fn(pfx) {
		return pfx
	}

	pfx = r.DefaultPrefix()
	if pfx != nil && fn(pfx) {
		return pfx
	}

	pfx = r.MemberPrefix(mme)
	if pfx != nil && fn(pfx) {
		return pfx
	}

	pfx = r.UserPrefix(me)
	if pfx != nil && fn(pfx) {
		return pfx
	}

	return nil
}

// LinePrefix finds the first suitable prefix that matches the given line.
func (r *Route) LinePrefix(
	g discord.GuildID,
	me discord.User,
	mme *discord.Member,
	line string,
) *Prefix {
	line = strings.TrimSpace(line)

	return r.findPrefix(g, mme, me, func(pfx *Prefix) bool {
		return strings.HasPrefix(line, pfx.Value)
	})
}
