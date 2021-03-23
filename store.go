package route

import "github.com/diamondburned/arikawa/v2/discord"

// Store is a generic store for data that a Route may need.
//
// Implementations should be concurrent-safe.
type Store interface {
	// GetPrefix fetches a prefix for the given Guild.
	// A null Guild should be used for a global prefix.
	GetPrefix(g discord.GuildID) (string, bool)

	// SetPrefix stores a prefix for the given Guild.
	// A null Guild should be used for a global prefix.
	SetPrefix(g discord.GuildID, v string)
}
