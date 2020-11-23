package route

import (
	"errors"
	"fmt"
	"strings"

	"github.com/diamondburned/arikawa/discord"
)

var (
	// ErrPrefixUnset occurs when the prefix value isn't set in the Settings.
	ErrPrefixUnset = errors.New("prefix unset")

	// ErrNoPrefix occurs when no suitable prefix was found.
	ErrNoPrefix = errors.New("no suitable prefix")
)

// Prefix is a command prefix.
type Prefix struct {
	Value string
	Clean string
}

// GuildPrefix finds the prefix for a given Guild.
func (r *Route) GuildPrefix(g discord.GuildID) (*Prefix, error) {
	set, _ := r.LoadSettings(g)
	if set.Prefix == "" {
		return nil, ErrPrefixUnset
	}

	return &Prefix{
		Value: set.Prefix,
		Clean: set.Prefix,
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

// FindPrefix finds the first suitable prefix, where "suitable" is defined as a truthy return from fn.
func (r *Route) FindPrefix(g discord.GuildID, fn func(*Prefix) bool) *Prefix {
	pfx, err := r.GuildPrefix(g)
	if err == nil && fn(pfx) {
		return pfx
	}

	pfx = r.UserPrefix()
	if fn(pfx) {
		return pfx
	}

	pfx, err = r.MemberPrefix(g)
	if err == nil && fn(pfx) {
		return pfx
	}

	return nil
}

// LinePrefix finds the first suitable prefix that matches the given line.
func (r *Route) LinePrefix(g discord.GuildID, line string) *Prefix {
	line = strings.TrimSpace(line)

	return r.FindPrefix(g, func(pfx *Prefix) bool {
		return strings.HasPrefix(line, pfx.Value)
	})
}
