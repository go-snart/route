package route

import (
	"fmt"
	"strings"
	"sync"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/superloach/confy"
)

// KeyPrefix is the Confy key used to load/store prefixes.
const KeyPrefix = "prefix"

// GlobalGuildID is the GuildID used for global configurations.
const GlobalGuildID = discord.NullGuildID

// PrefixStore describes a concurrent-safe store of Guild-specific command prefixes.
type PrefixStore struct {
	Confy confy.Confy

	ma map[discord.GuildID]string
	mu sync.RWMutex
}

// OpenPrefixStore creates a usable PrefixStore and calls Load.
func OpenPrefixStore(c confy.Confy) (*PrefixStore, error) {
	pfxs := &PrefixStore{
		Confy: c,

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
	err := p.Confy.Get(KeyPrefix, &p.ma)
	p.mu.Unlock()

	if err != nil {
		return fmt.Errorf("confy load %q: %w", KeyPrefix, err)
	}

	return nil
}

// Store updates the Confy with data from the PrefixStore.
func (p *PrefixStore) Store() error {
	p.mu.RLock()
	err := p.Confy.Set(KeyPrefix, p.ma)
	p.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("confy store %q: %w", KeyPrefix, err)
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

// ForLine finds the first suitable prefix that matches the given line.
func (p *PrefixStore) ForLine(
	g discord.GuildID,
	me discord.User,
	mme *discord.Member,
	line string,
) (Prefix, bool) {
	line = strings.TrimSpace(line)

	// guild prefix
	pfxv, ok := p.Get(g)
	if !ok {
		// fallback to default prefix
		pfxv, ok = p.Get(GlobalGuildID)
	}

	if ok && strings.HasPrefix(line, pfxv) {
		return Prefix{pfxv, pfxv}, true
	}

	// member prefix
	if mme != nil && strings.HasPrefix(line, mme.Mention()) {
		pfx := Prefix{
			Value: mme.Mention(),
			Clean: "@" + me.Username + " ",
		}

		if mme.Nick != "" {
			pfx.Clean = "@" + mme.Nick + " "
		}

		return pfx, true
	}

	// user prefix
	if strings.HasPrefix(line, me.Mention()) {
		return Prefix{
			Value: me.Mention(),
			Clean: "@" + me.Username + " ",
		}, true
	}

	return Prefix{"", ""}, false
}
