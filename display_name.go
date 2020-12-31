package route

import "log"

// DefaultDisplayName is used when DisplayName can't find a username or nickname.
const DefaultDisplayName = "Snart"

// DisplayName returns the State's display name for the given trigger.
//
// If Me errors, it simply returns DefaultDisplayName.
// If Me returns a Member, it checks for a Nick and returns that.
// Otherwise, it uses the User's Username.
func (t *Trigger) DisplayName() string {
	if t.Message.GuildID.IsNull() {
		me, err := t.Route.GetMe()
		if err != nil {
			log.Println("get me:", err)

			return DefaultDisplayName
		}

		return me.Username
	}

	mme, err := t.GetMMe()
	if err != nil {
		log.Println("get mme:", err)

		return DefaultDisplayName
	}

	if mme.Nick != "" {
		return mme.Nick
	}

	return mme.User.Username
}
