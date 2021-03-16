package route

import (
	"strings"
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/json"
)

// Prefix is a command prefix.
type Prefix struct {
	Value string
	Clean string
}

// PrefixStore is a concurrent-safe store of Prefix.Values for Guilds.
type PrefixStore struct {
	mu sync.RWMutex
	ma map[discord.GuildID]string
}

// GetPfx fetches a prefix for the given Guild.
// Returns a Prefix with Value and Clean set to the prefix value, for convenience.
func (p *PrefixStore) GetPfx(g discord.GuildID) (Prefix, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	v, ok := p.ma[g]

	return Prefix{
		Value: v,
		Clean: v,
	}, ok
}

// SetPfx stores a Prefix for the given Guild.
func (p *PrefixStore) SetPfx(g discord.GuildID, v string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.ma[g] = v
}

// UnmarshalJSON implements json.Unmarshaler.
func (p *PrefixStore) UnmarshalJSON(bs []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return json.Unmarshal(bs, &p.ma)
}

// MarshalJSON implements json.Unmarshaler.
func (p *PrefixStore) MarshalJSON() ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return json.Marshal(p.ma)
}

func (p *PrefixStore) findPrefix(
	g discord.GuildID,
	mme *discord.Member,
	me discord.User,
	fn func(Prefix) bool,
) (Prefix, bool) {
	// guild prefix
	pfx, ok := p.GetPfx(g)
	if !ok {
		// fallback to default prefix
		pfx, ok = p.GetPfx(discord.NullGuildID)
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
func (p *PrefixStore) LinePrefix(
	g discord.GuildID,
	me discord.User,
	mme *discord.Member,
	line string,
) (Prefix, bool) {
	line = strings.TrimSpace(line)

	return p.findPrefix(g, mme, me, func(pfx Prefix) bool {
		return strings.HasPrefix(line, pfx.Value)
	})
}
