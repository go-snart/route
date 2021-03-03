package route_test

import (
	"log"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

func TestHelp(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(s)

	c, _ := r.GetCmd("help")

	const (
		guild   = 1234567890
		channel = 1234567890
	)

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
				Name:  "help",
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
		Prefix: testPfx,
		Flags:  route.HelpFlags{},
	})
	if err != nil {
		t.Errorf("help: %s", err)
	}

	m.Eval()
}

func TestHelpHide(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(s)

	c, _ := r.GetCmd("help")
	c.Hide = true

	r.DelCmd(c.Name)
	r.AddCmds(c)

	const (
		guild   = 1234567890
		channel = 1234567890
	)

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
			Fields: nil,
		}},
	})

	err := c.Func(&route.Trigger{
		Route:   r,
		Command: c,

		Message: discord.Message{
			GuildID:   guild,
			ChannelID: channel,
		},
		Prefix: testPfx,
		Flags:  route.HelpFlags{},
	})
	if err != nil {
		t.Errorf("help: %s", err)
	}

	m.Eval()
}

func TestHelpHelpception(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(s)

	c, _ := r.GetCmd("help")

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
	r := route.New(s)

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

	const line = ".help help"

	tr, err := r.Trigger(
		route.Prefix{
			Value: ".",
		},
		discord.Message{
			ChannelID: channel,
		},
		line,
	)
	if err != nil {
		t.Errorf("get trigger %q: %s", line, err)
	}

	log.Printf("%#v", tr)
	log.Printf("%#v", tr.FlagSet)

	err = tr.Command.Func(tr)
	if err != nil {
		t.Errorf("tr run: %s", err)
	}

	m.Eval()
}

func TestHelpUsageUnknown(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(s)

	c, _ := r.GetCmd("help")

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
