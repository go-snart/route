package route_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

func TestDisplayNameNullMeError(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	tr := &route.Trigger{
		Router: r,
		Message: discord.Message{
			GuildID: discord.NullGuildID,
		},
	}

	m.Error(http.MethodGet, "/users/@me", httputil.HTTPError{Status: 404})

	if dn := tr.DisplayName(); dn != route.DefaultDisplayName {
		t.Errorf("expect %q\ngot %q", route.DefaultDisplayName, dn)
	}

	m.Eval()
}

func TestDisplayNameNull(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	tr := &route.Trigger{
		Router: r,
		Message: discord.Message{
			GuildID: discord.NullGuildID,
		},
	}

	m.Me(testMe)

	if dn := tr.DisplayName(); dn != testMe.Username {
		t.Errorf("expect %q\ngot %q", testMe.Username, dn)
	}

	m.Eval()
}

func TestDisplayNameMMeError(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	const guild = 1234567890

	tr := &route.Trigger{
		Router: r,
		Message: discord.Message{
			GuildID: guild,
		},
	}

	m.Me(testMe)
	m.Error(
		http.MethodGet,
		fmt.Sprintf("/guilds/%d/members/%d", guild, testMe.ID),
		httputil.HTTPError{Status: 404},
	)

	if dn := tr.DisplayName(); dn != route.DefaultDisplayName {
		t.Errorf("expect %q\ngot %q", route.DefaultDisplayName, dn)
	}

	m.Eval()
}

func TestDisplayNameMMeNick(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	tr := &route.Trigger{
		Router: r,
		Message: discord.Message{
			GuildID: 123,
		},
	}

	m.Me(testMe)
	m.Member(tr.Message.GuildID, testMMeNick)

	if dn := tr.DisplayName(); dn != testMMeNick.Nick {
		t.Errorf("expect %q\ngot %q", testMMeNick.Nick, dn)
	}

	m.Eval()
}

func TestDisplayNameMMeUser(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	tr := &route.Trigger{
		Router: r,
		Message: discord.Message{
			GuildID: 123,
		},
	}

	m.Me(testMe)
	m.Member(tr.Message.GuildID, testMMe)

	if dn := tr.DisplayName(); dn != testMe.Username {
		t.Errorf("expect %q\ngot %q", testMe.Username, dn)
	}

	m.Eval()
}
