package route

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/diamondburned/arikawa/v2/discord"
)

// Setup contains setup parameters for a Route.
type Setup struct {
	PrefixFile string                     `flag:"prefixf" usage:"file containing prefixes"`
	Prefixes   map[discord.GuildID]string `flag:""`
}

// LoadPrefixes loads a prefix map from the file named by PrefixFile.
// If s.Prefixes is set, that is returned without consulting PrefixFile.
func (s Setup) LoadPrefixes() (map[discord.GuildID]string, error) {
	if s.Prefixes != nil {
		return s.Prefixes, nil
	}

	f, err := os.Open(s.PrefixFile)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", s.PrefixFile, err)
	}
	defer f.Close()

	m := make(map[discord.GuildID]string)

	err = json.NewDecoder(f).Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("json decode m: %w", err)
	}

	return m, nil
}
