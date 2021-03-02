package route

import "log"

// DisplayName is the default returned by DisplayName.
const DisplayName = "Snart"

// DisplayName returns the State's display name for the given trigger.
func (t *Trigger) DisplayName() string {
	me, err := t.Router.State.Me()
	if err != nil {
		log.Printf("error: get me: %s", err)

		return DisplayName
	}

	if t.Message.GuildID.IsNull() {
		return me.Username
	}

	mme, err := t.Router.State.Member(t.Message.GuildID, me.ID)
	if err != nil {
		log.Printf("error: get mme: %s", err)

		return DisplayName
	}

	if mme.Nick != "" {
		return mme.Nick
	}

	return mme.User.Username
}
