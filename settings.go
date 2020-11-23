package route

import (
	"fmt"

	"github.com/diamondburned/arikawa/discord"
)

// Settings holds global or Guild-specific settings.
type Settings struct {
	Prefix string
}

// LoadSettings retrieves Settings for the given Guild from the Route's DB.
func (r *Route) LoadSettings(g discord.GuildID) (Settings, error) {
	key := r.settingsKey(g)
	s := Settings{}

	err := r.DB.Get(key, &s)
	if err != nil {
		if g == 0 {
			return s, r.SaveSettings(g, s)
		}

		return r.LoadSettings(0)
	}

	return s, nil
}

// SaveSettings stores the given Settings for the given Guild into the Route's DB.
func (r *Route) SaveSettings(g discord.GuildID, s Settings) error {
	key := r.settingsKey(g)

	return r.DB.Set(key, s)
}

func (r *Route) settingsKey(g discord.GuildID) string {
	if g > 0 {
		return fmt.Sprintf("settings_%d", g)
	}

	return "settings"
}
