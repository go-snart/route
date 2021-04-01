package route

import (
	"fmt"
	"strings"
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/superloach/confy"
)

// GuildIDGlobal is the GuildID used for global configurations.
const GuildIDGlobal = discord.NullGuildID

// PrefixStore describes a concurrent-safe store of Guild-specific command prefixes.
type PrefixStore struct {
	Confy *confy.Confy
	Key   string

	ma map[discord.GuildID]string
	mu sync.RWMutex
}

// OpenPrefixStore creates a usable PrefixStore and calls Load.
func OpenPrefixStore(c *confy.Confy, key string) (*PrefixStore, error) {
	pfxs := &PrefixStore{
		Confy: c,
		Key:   key,

		ma: map[discord.GuildID]string{},
		mu: sync.RWMutex{},
	}

	if err := pfxs.Load(); err != nil {
		return nil, fmt.Errorf("pfxs load: %w", err)
	}

	return pfxs, nil
}

// Load updates the PrefixStore with data from the Confy.
func (p *PrefixStore) Load() error {
	p.mu.Lock()
	err := p.Confy.Load(p.Key, &p.ma)
	p.mu.Unlock()

	if err != nil {
		return fmt.Errorf("confy load %q: %w", p.Key, err)
	}

	return nil
}

// Store updates the Confy with data from the PrefixStore.
func (p *PrefixStore) Store() error {
	p.mu.RLock()
	err := p.Confy.Store(p.Key, p.ma)
	p.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("confy store %q: %w", p.Key, err)
	}

	return nil
}

// Get allows looking up a Prefix by GuildID.
func (p *PrefixStore) Get(g discord.GuildID) (string, bool) {
	p.mu.RLock()
	pfxv, ok := p.ma[g]
	p.mu.RUnlock()

	return pfxv, ok
}

// Set allows storing a prefix for a given GuildID.
func (p *PrefixStore) Set(g discord.GuildID, pfxv string) {
	p.mu.Lock()
	p.ma[g] = pfxv
	p.mu.Unlock()
}

// Del removes the prefix for the given GuildID from the PrefixStore.
func (p *PrefixStore) Del(g discord.GuildID) {
	p.mu.Lock()
	delete(p.ma, g)
	p.mu.Unlock()
}

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
	pfxv, ok := r.Prefix.Get(g)
	if !ok {
		// fallback to default prefix
		pfxv, ok = r.Prefix.Get(GuildIDGlobal)
	}

	pfx := Prefix{pfxv, pfxv}
	if ok && used(pfx) {
		return pfx, true
	}

	// member prefix
	if mme != nil {
		pfx = Prefix{
			Value: mme.Mention(),
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
		Clean: "@" + me.Username + " ",
	}
	if used(pfx) {
		return pfx, true
	}

	return Prefix{"", ""}, false
}
