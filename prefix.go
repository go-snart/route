package route

import (
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/discord"
)

// Prefix is a command prefix.
type Prefix struct {
	Value string
	Clean string
}

// GuildPrefix finds the prefix for a given Guild.
func (r *Route) GuildPrefix(g discord.GuildID) (*Prefix, error) {
	key := "prefix"
	val := ""

	if g > 0 {
		key = fmt.Sprintf("prefix_%d", g)
	}

	err := r.DB.Get(key, &val)
	if err != nil {
		return nil, fmt.Errorf("get %q: %w", key, err)
	}

	return &Prefix{
		Value: val,
		Clean: val,
	}, nil
}

// UserPrefix makes the user mention prefix.
func (r *Route) UserPrefix() *Prefix {
	me := r.State.Ready.User

	return &Prefix{
		Value: me.Mention(),
		Clean: "@" + me.Username,
	}
}

// MemberPrefix finds the member mention prefix for a given Guild.
func (r *Route) MemberPrefix(g discord.GuildID) (*Prefix, error) {
	me := r.State.Ready.User

	mme, err := r.State.Member(g, me.ID)
	if err != nil {
		return nil, fmt.Errorf("member %d %d: %w", g, me.ID, err)
	}

	if mme.Nick != "" {
		return &Prefix{
			Value: mme.Mention(),
			Clean: "@" + mme.Nick,
		}, nil
	}

	return &Prefix{
		Value: mme.Mention(),
		Clean: "@" + mme.User.Username,
	}, nil
}

// FindPrefix finds a suitable prefix for the line.
// Checks guild prefix, default prefix, user mention, and member mention (in that order).
func (r *Route) FindPrefix(g discord.GuildID, line string) *Prefix {
	line = strings.TrimSpace(line)

	pfx, err := r.GuildPrefix(g)
	if err == nil && strings.HasPrefix(line, pfx.Value) {
		return pfx
	}

	pfx, err = r.GuildPrefix(0)
	if err == nil && strings.HasPrefix(line, pfx.Value) {
		return pfx
	}

	pfx = r.UserPrefix()
	if strings.HasPrefix(line, pfx.Value) {
		return pfx
	}

	pfx, err = r.MemberPrefix(g)
	if err == nil && strings.HasPrefix(line, pfx.Value) {
		return pfx
	}

	return nil
}
