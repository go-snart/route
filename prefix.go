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

// FindPrefix finds the first prefix that satisfies fn.
func (r *Route) FindPrefix(
	g discord.GuildID,
	mme *discord.Member,
	me discord.User,
	fn func(Prefix) bool,
) (Prefix, bool) {
	// guild prefix
	pfxv, ok := r.GetPrefix(g)
	if !ok {
		// fallback to default prefix
		pfxv, ok = r.GetPrefix(discord.NullGuildID)
	}

	pfx := Prefix{
		Value: pfxv,
		Clean: pfxv,
	}
	if ok && fn(pfx) {
		return pfx, true
	}

	// member prefix
	if mme != nil {
		pfx = Prefix{
			Value: mme.Mention(),
			Clean: me.Username,
		}

		if mme.Nick != "" {
			pfx.Clean = "@" + mme.Nick
		}

		if fn(pfx) {
			return pfx, true
		}
	}

	// user prefix
	pfx = Prefix{
		Value: me.Mention(),
		Clean: "@" + me.Username,
	}
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

	return r.FindPrefix(g, mme, me, func(pfx Prefix) bool {
		return strings.HasPrefix(line, pfx.Value)
	})
}
