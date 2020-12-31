package route_test

import (
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

func TestDisplayNameNullMeError(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	tr := &route.Trigger{
		Route: r,
		Message: discord.Message{
			GuildID: discord.NullGuildID,
		},
	}

	m.Error(http.MethodGet, "/users/@me", httputil.HTTPError{Status: 404})

	dn := tr.DisplayName()
	if dn != route.DefaultDisplayName {
		t.Errorf("expect %q\ngot %q", route.DefaultDisplayName, dn)
	}

	m.Eval()
}

func TestDisplayNameNull(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	tr := &route.Trigger{
		Route: r,
		Message: discord.Message{
			GuildID: discord.NullGuildID,
		},
	}

	m.Me(testMe)

	dn := tr.DisplayName()
	if dn != testMe.Username {
		t.Errorf("expect %q\ngot %q", testMe.Username, dn)
	}

	m.Eval()
}

func TestDisplayNameMMeError(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	tr := &route.Trigger{
		Route: r,
	}

	m.Error(http.MethodGet, "/users/@me", httputil.HTTPError{Status: 404})

	dn := tr.DisplayName()
	if dn != route.DefaultDisplayName {
		t.Errorf("expect %q\ngot %q", route.DefaultDisplayName, dn)
	}

	m.Eval()
}

func TestDisplayNameMMeNick(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	tr := &route.Trigger{
		Route: r,
		Message: discord.Message{
			GuildID: 123,
		},
	}

	m.Me(testMe)
	m.Member(tr.Message.GuildID, testMMeNick)

	dn := tr.DisplayName()
	if dn != testMMeNick.Nick {
		t.Errorf("expect %q\ngot %q", testMMeNick.Nick, dn)
	}

	m.Eval()
}

func TestDisplayNameMMeUser(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	tr := &route.Trigger{
		Route: r,
		Message: discord.Message{
			GuildID: 123,
		},
	}

	m.Me(testMe)
	m.Member(tr.Message.GuildID, testMMe)

	dn := tr.DisplayName()
	if dn != testMe.Username {
		t.Errorf("expect %q\ngot %q", testMe.Username, dn)
	}

	m.Eval()
}
