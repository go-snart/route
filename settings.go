package route

import (
	"fmt"

	"github.com/diamondburned/arikawa/v2/discord"
)

// DefaultGuild is the GuildID used for global settings.
const DefaultGuild discord.GuildID = 0

// Settings holds global or Guild-specific settings.
type Settings struct {
	Prefix string
}

// Load retrieves Settings for the given Guild from the Route's DB.
func (r *Route) Load(g discord.GuildID) (Settings, error) {
	s := Settings{}

	err := r.DB.Get(r.key(g), &s)
	if err != nil {
		if g == DefaultGuild {
			return s, r.Save(g, s)
		}

		return r.Load(DefaultGuild)
	}

	return s, nil
}

// Save stores the given Settings for the given Guild into the Route's DB.
func (r *Route) Save(g discord.GuildID, s Settings) error {
	return r.DB.Set(r.key(g), s)
}

func (r *Route) key(g discord.GuildID) string {
	if g == DefaultGuild {
		return "settings"
	}

	return fmt.Sprintf("settings_%d", g)
}
