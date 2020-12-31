package route_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

func TestLinePrefixGuild(t *testing.T) {
	r := route.New(testDB(), nil)

	const (
		guild = 1234567890
		pfxv  = "//"
	)

	set := route.Settings{
		Prefix: pfxv,
	}

	err := r.Save(guild, set)
	if err != nil {
		t.Errorf("save set %d %v: %s", guild, set, err)
	}

	pfx := r.LinePrefix(guild, pfxv)
	expect := &route.Prefix{
		Value: pfxv,
		Clean: pfxv,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestLinePrefixUser(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	m.Me(testMe)

	pfx := r.LinePrefix(discord.NullGuildID, testMe.Mention())
	expect := &route.Prefix{
		Value: testMe.Mention(),
		Clean: "@User",
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixMember(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	const guild = 666

	m.Me(testMe)
	m.Member(guild, testMMe)

	pfx := r.LinePrefix(guild, testMMe.Mention())
	expect := &route.Prefix{
		Value: testMMe.Mention(),
		Clean: "@User",
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixMemberNick(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	const guild = 666

	m.Me(testMe)
	m.Member(guild, testMMeNick)

	pfx := r.LinePrefix(guild, testMMeNick.Mention())
	expect := &route.Prefix{
		Value: testMMeNick.Mention(),
		Clean: "@" + testMMeNick.Nick,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixMemberErr(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	const guild = 666

	m.Me(testMe)
	m.Error(
		http.MethodGet,
		fmt.Sprintf("/guilds/%d/members/%d", guild, testMe.ID),
		httputil.HTTPError{Status: 404},
	)

	pfx := r.LinePrefix(guild, testMMe.Mention())
	expect := (*route.Prefix)(nil)

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixNil(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	const guild = 666

	m.Me(testMe)
	m.Member(guild, testMMe)

	pfx := r.LinePrefix(guild, "")
	expect := (*route.Prefix)(nil)

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestFindPrefixUserError(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	const guild = 666

	m.Error(http.MethodGet, "/users/@me", httputil.HTTPError{Status: 404})
	m.Error(http.MethodGet, "/users/@me", httputil.HTTPError{Status: 404})

	pfx := r.FindPrefix(guild, func(*route.Prefix) bool {
		return true
	})
	if pfx != nil {
		t.Errorf("expect nil\ngot: %v", pfx)
	}

	m.Eval()
}
