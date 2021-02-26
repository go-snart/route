package route_test

import (
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

func TestLinePrefixGuild(t *testing.T) {
	r := route.New(testSettings, nil)

	const (
		guild = 1234567890
		pfxv  = "//"
	)

	r.Guilds[guild] = route.Settings{
		Prefix: pfxv,
	}

	pfx := r.LinePrefix(guild, testMe, nil, pfxv)

	expect := &route.Prefix{
		Value: pfxv,
		Clean: pfxv,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestLinePrefixDefault(t *testing.T) {
	r := route.New(testSettings, nil)

	const guild = 1234567890

	r.Guilds[discord.NullGuildID] = route.Settings{
		Prefix: "owo!",
	}

	pfx := r.LinePrefix(guild, testMe, nil, "owo!uwu")
	expect := &route.Prefix{
		Value: "owo!",
		Clean: "owo!",
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestLinePrefixUser(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	pfx := r.LinePrefix(discord.NullGuildID, testMe, nil, testMe.Mention())
	expect := &route.Prefix{
		Value: testMe.Mention(),
		Clean: "@User",
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixMemberNick(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	const guild = 666

	pfx := r.LinePrefix(guild, testMe, &testMMeNick, testMMeNick.Mention())
	expect := &route.Prefix{
		Value: testMMeNick.Mention(),
		Clean: "@" + testMMeNick.Nick,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixNil(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	const guild = 666

	pfx := r.LinePrefix(guild, testMe, nil, "")
	if pfx != nil {
		t.Errorf("should be nil, got %#v", pfx)
	}

	m.Eval()
}
