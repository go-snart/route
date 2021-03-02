package route

import (
	"github.com/diamondburned/arikawa/v2/discord"
)

// BaseGuild is the Guild ID for base configuration.
const BaseGuild = discord.NullGuildID

// Guild holds base or Guild-specific configuration.
type Guild struct {
	Prefix string
}

// GetGuild retrieves the given Guild.
func (r *Route) GetGuild(g discord.GuildID) (Guild, bool) {
	r.setMu.RLock()
	defer r.setMu.RUnlock()

	s, ok := r.setMap[g]

	return s, ok
}

// SetGuild stores the given Guild.
func (r *Route) SetGuild(g discord.GuildID, s Guild) {
	r.setMu.Lock()
	defer r.setMu.Unlock()

	r.setMap[g] = s
}

// ImportGuilds does a bulk import of a Guilds map.
func (r *Route) ImportGuilds(m map[discord.GuildID]Guild) {
	r.setMu.Lock()
	defer r.setMu.Unlock()

	for g, s := range m {
		r.setMap[g] = s
	}
}

// ExportGuilds does a bulk export of a Guilds map.
func (r *Route) ExportGuilds() map[discord.GuildID]Guild {
	r.setMu.RLock()
	defer r.setMu.RUnlock()

	m := make(map[discord.GuildID]Guild, len(r.setMap))

	for g, s := range r.setMap {
		m[g] = s
	}

	return m
}
