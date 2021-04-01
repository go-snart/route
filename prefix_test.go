package route_test

import (
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

func TestLinePrefixGuild(t *testing.T) {
	t.Parallel()

	const (
		guild = 1234567890
		pfxv  = "%%"
	)

	c := testConfy()
	r := testRoute(t, nil, c)

	r.Prefix.Set(guild, pfxv)

	pfx, _ := r.LinePrefix(guild, testMe, nil, pfxv)
	expect := route.Prefix{
		Value: pfxv,
		Clean: pfxv,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestLinePrefixNull(t *testing.T) {
	t.Parallel()

	const pfxv = "test!"

	c := testConfy()
	r := testRoute(t, nil, c)

	const guild = 1234567890

	r.Prefix.Set(guild, "test!")

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
	t.Parallel()

	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

	pfx, _ := r.LinePrefix(discord.NullGuildID, testMe, nil, testMe.Mention())
	expect := route.Prefix{
		Value: testMe.Mention(),
		Clean: "@" + testMe.Username + " ",
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixMemberNick(t *testing.T) {
	t.Parallel()

	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

	const guild = 666

	pfx, _ := r.LinePrefix(guild, testMe, &testMMeNick, testMMeNick.Mention())
	expect := route.Prefix{
		Value: testMMeNick.Mention(),
		Clean: "@" + testMMeNick.Nick + " ",
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixNil(t *testing.T) {
	t.Parallel()

	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

	const guild = 666

	pfx, ok := r.LinePrefix(guild, testMe, nil, "")
	if ok {
		t.Errorf("should be !ok, got %#v", pfx)
	}

	m.Eval()
}
