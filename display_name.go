package route

import "github.com/go-snart/lob"

// DefaultDisplayName is the default returned by DisplayName.
const DefaultDisplayName = "Snart"

// DisplayName returns the State's display name for the given trigger.
func (t *Trigger) DisplayName() string {
	me, err := t.Router.State.Me()
	if err != nil {
		_ = lob.Std.Error("get me: %w", err)

		return DefaultDisplayName
	}

	if t.Message.GuildID.IsNull() {
		return me.Username
	}

	mme, err := t.Router.State.Member(t.Message.GuildID, me.ID)
	if err != nil {
		_ = lob.Std.Error("get mme: %w", err)

		return DefaultDisplayName
	}

	if mme.Nick != "" {
		return mme.Nick
	}

	return mme.User.Username
}
