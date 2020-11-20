package route

import (
	"errors"
	"fmt"

	"github.com/diamondburned/arikawa/discord"

	"github.com/go-snart/db/impl"
)

type Settings struct {
	Prefix string `json:"prefix,omitempty"`
}

func (r *Route) LoadSettings(g discord.GuildID) (Settings, error) {
	key := r.settingsKey(g)
	s := Settings{}

	err := r.DB.Get(key, &s)
	if errors.As(err, &impl.NoKeyError{}) {
		if g == 0 {
			s = Settings{}
			return s, r.SaveSettings(g, s)
		}

		return r.LoadSettings(0)
	}

	return s, err
}

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
