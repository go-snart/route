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

// LinePrefix finds the first suitable prefix that matches the given line.
func (r *Route) LinePrefix(
	g discord.GuildID,
	me discord.User,
	mme *discord.Member,
	line string,
) (Prefix, bool) {
	line = strings.TrimSpace(line)

	used := func(pfx Prefix) bool {
		if strings.HasPrefix(line, pfx.Value) {
			return true
		}

		if strings.HasPrefix(line, pfx.Clean) {
			return true
		}

		return false
	}

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
	if ok && used(pfx) {
		return pfx, true
	}

	// member prefix
	if mme != nil {
		pfx = Prefix{
			Value: mme.Mention() + " ",
			Clean: "@" + me.Username + " ",
		}

		if mme.Nick != "" {
			pfx.Clean = "@" + mme.Nick + " "
		}

		if used(pfx) {
			return pfx, true
		}
	}

	// user prefix
	pfx = Prefix{
		Value: me.Mention(),
		Clean: "@" + me.Username,
	}
	if used(pfx) {
		return pfx, true
	}

	return Prefix{}, false
}
