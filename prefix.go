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
func (r *Route) GuildPrefix(g discord.GuildID) (Prefix, bool) {
	set, ok := r.GetGuild(g)
	if !ok {
		return Prefix{}, ok
	}

	return Prefix{
		Value: set.Prefix,
		Clean: set.Prefix,
	}, true
}

// DefaultPrefix returns the default Prefix.
// This is the Prefix for discord.NullGuildID.
func (r *Route) DefaultPrefix() (Prefix, bool) {
	return r.GuildPrefix(BaseGuild)
}

func memberPrefix(mme *discord.Member) Prefix {
	if mme.Nick != "" {
		return Prefix{
			Value: mme.Mention(),
			Clean: "@" + mme.Nick,
		}
	}

	return Prefix{
		Value: mme.Mention(),
		Clean: "@" + mme.User.Username,
	}
}

func userPrefix(me discord.User) Prefix {
	return Prefix{
		Value: me.Mention(),
		Clean: "@" + me.Username,
	}
}

func (r *Route) findPrefix(
	g discord.GuildID,
	mme *discord.Member,
	me discord.User,
	fn func(Prefix) bool,
) (Prefix, bool) {
	pfx, ok := r.GuildPrefix(g)
	if ok && fn(pfx) {
		return pfx, true
	}

	pfx, ok = r.DefaultPrefix()
	if ok && fn(pfx) {
		return pfx, true
	}

	if mme != nil {
		pfx = memberPrefix(mme)
		if fn(pfx) {
			return pfx, true
		}
	}

	pfx = userPrefix(me)
	if fn(pfx) {
		return pfx, true
	}

	return Prefix{}, false
}

// LinePrefix finds the first suitable prefix that matches the given line.
func (r *Route) LinePrefix(
	g discord.GuildID,
	me discord.User,
	mme *discord.Member,
	line string,
) (Prefix, bool) {
	line = strings.TrimSpace(line)

	return r.findPrefix(g, mme, me, func(pfx Prefix) bool {
		return strings.HasPrefix(line, pfx.Value)
	})
}
