package route_test

import (
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

func TestLinePrefixGuild(t *testing.T) {
	r := route.New(testGuild, nil)

	const (
		guild = 1234567890
		pfxv  = "//"
	)

	r.SetGuild(guild, route.Guild{
		Prefix: pfxv,
	})

	pfx, _ := r.LinePrefix(guild, testMe, nil, pfxv)

	expect := route.Prefix{
		Value: pfxv,
		Clean: pfxv,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestLinePrefixBase(t *testing.T) {
	const pfxv = "test!"

	r := route.New(route.Guild{
		Prefix: pfxv,
	}, nil)

	const guild = 1234567890

	pfx, _ := r.LinePrefix(guild, testMe, nil, "test!uwu")
	expect := route.Prefix{
		Value: pfxv,
		Clean: pfxv,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestLinePrefixUser(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testGuild, s)

	pfx, _ := r.LinePrefix(discord.NullGuildID, testMe, nil, testMe.Mention())
	expect := route.Prefix{
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
	r := route.New(testGuild, s)

	const guild = 666

	pfx, _ := r.LinePrefix(guild, testMe, &testMMeNick, testMMeNick.Mention())
	expect := route.Prefix{
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
	r := route.New(testGuild, s)

	const guild = 666

	pfx, ok := r.LinePrefix(guild, testMe, nil, "")
	if ok {
		t.Errorf("should be !ok, got %#v", pfx)
	}

	m.Eval()
}
