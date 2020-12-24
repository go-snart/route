package route_test

import (
	"fmt"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/go-snart/route"
	"github.com/mavolin/dismock/pkg/dismock"
)

func TestHelp(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(nil, s)

	r.State.Ready.User.Username = "username"

	r.Cats["abc"] = nil

	c := r.Cats["route"][0]

	const channel = 1234567890

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	m.SendEmbed(discord.Message{
		ChannelID: channel,
		Embeds: []discord.Embed{{
			Title:       r.State.Ready.User.Username + " Help",
			Description: fmt.Sprintf("prefix: `%s`", pfx.Clean),
			Footer: &discord.EmbedFooter{
				Text: "use the `-help` flag on a command for detailed help",
			},
			Fields: []discord.EmbedField{{
				Name:  route.RouteCat,
				Value: fmt.Sprintf("`%s%s`: *%s*", pfx.Clean, c.Name, c.Desc),
			}},
		}},
	})

	err := c.Func(&route.Trigger{
		Route:   r,
		Command: c,

		Message: discord.Message{
			ChannelID: channel,
		},
		Prefix: pfx,
		Flags: route.HelpFlags{
			Help: false,
		},
	})
	if err != nil {
		t.Errorf("help: %s", err)
	}

	m.Eval()
}

func TestHelpHelpception(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(nil, s)

	c := r.Cats["route"][0]

	const channel = 1234567890

	m.SendMessage(nil, discord.Message{
		ChannelID: channel,
		Content:   "helpception :thinking:",
	})

	err := c.Func(&route.Trigger{
		Route:   r,
		Command: c,

		Message: discord.Message{
			ChannelID: channel,
		},
		Flags: route.HelpFlags{
			Help: true,
		},
	})
	if err != nil {
		t.Errorf("help: %s", err)
	}

	m.Eval()
}

func TestHelpUsage(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(nil, s)

	c := r.Cats["route"][0]

	const channel = 1234567890

	m.SendMessage(&discord.Embed{
		Title:       "`help` usage",
		Type:        "rich",
		Description: "a help menu",
		Color:       3158064,
		Fields: []discord.EmbedField{{
			Name:   "flag `-help`",
			Value:  "helpception\ndefault: `false`",
			Inline: false,
		}},
	}, discord.Message{
		ChannelID: channel,
	})

	err := c.Func(&route.Trigger{
		Route:   r,
		Command: c,

		Message: discord.Message{
			ChannelID: channel,
		},
		Flags: route.HelpFlags{
			Help: false,
		},
		Args: []string{
			"help",
		},
	})
	if err != nil {
		t.Errorf("help: %s", err)
	}

	m.Eval()
}

func TestHelpUsageUnknown(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(nil, s)

	c := r.Cats["route"][0]

	const channel = 1234567890

	m.SendMessage(nil, discord.Message{
		ChannelID: channel,
		Content:   "command `abc` not known",
	})

	err := c.Func(&route.Trigger{
		Route:   r,
		Command: c,

		Message: discord.Message{
			ChannelID: channel,
		},
		Flags: route.HelpFlags{
			Help: false,
		},
		Args: []string{
			"abc",
		},
	})
	if err != nil {
		t.Errorf("help: %s", err)
	}

	m.Eval()
}
