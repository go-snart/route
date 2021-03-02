package route

import (
	"errors"
	"log"

	"github.com/diamondburned/arikawa/v2/discord"
)

// ErrSetBase occurs when one attempts to set BaseID.
var ErrSetBase = errors.New("can't set base settings (use New)")

// BaseID is the Guild ID used for base settings.
const BaseID = discord.NullGuildID

// Settings holds base or Guild-specific settings.
type Settings struct {
	Prefix string
}

// GetSettings retrieves the settings for a given Guild.
// Pass BaseID to get base settings.
func (r *Route) GetSettings(g discord.GuildID) (Settings, bool) {
	r.setMu.RLock()
	defer r.setMu.RUnlock()

	s, ok := r.setMap[g]

	return s, ok
}

// SetSettings stores the settings for a given Guild.
// Passing BaseID is not allowed; base settings are specified in New.
func (r *Route) SetSettings(g discord.GuildID, s Settings) {
	r.setMu.Lock()
	defer r.setMu.Unlock()

	if g == BaseID {
		log.Panic(ErrSetBase)
	}

	r.setMap[g] = s
}

// ImportSettings does a bulk import of a Settings map.
// Passing BaseID is not allowed; base settings are specified in New.
func (r *Route) ImportSettings(m map[discord.GuildID]Settings) {
	r.setMu.Lock()
	defer r.setMu.Unlock()

	for g, s := range m {
		if g == BaseID {
			log.Panic(ErrSetBase)
		}

		r.setMap[g] = s
	}
}

// ExportSettings does a bulk export of a Settings map.
// Includes BaseID for base settings as specified in New.
func (r *Route) ExportSettings() map[discord.GuildID]Settings {
	r.setMu.RLock()
	defer r.setMu.RUnlock()

	m := make(map[discord.GuildID]Settings, len(r.setMap))

	for g, s := range r.setMap {
		m[g] = s
	}

	return m
}
