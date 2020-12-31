package route_test

import (
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

func TestHelp(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(nil, s)

	const emptyCat = "abc"

	r.Cats[emptyCat] = nil

	c := r.Cats["route"][0]

	const (
		guild   = 1234567890
		channel = 1234567890
	)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	m.Me(testMe)
	m.Member(guild, testMMe)
	m.SendEmbed(discord.Message{
		ChannelID: channel,
		Embeds: []discord.Embed{{
			Title:       "User Help",
			Description: "prefix: `//`",
			Footer: &discord.EmbedFooter{
				Text: "use the `-help` flag on a command for detailed help",
			},
			Fields: []discord.EmbedField{{
				Name:  emptyCat,
				Value: "*no commands*",
			}, {
				Name:  route.RouteCat,
				Value: "`//help`: *a help menu*",
			}},
		}},
	})

	err := c.Func(&route.Trigger{
		Route:   r,
		Command: c,

		Message: discord.Message{
			GuildID:   guild,
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
