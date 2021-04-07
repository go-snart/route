package route_test

import (
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

func TestForLineGuild(t *testing.T) {
	t.Parallel()

	const (
		guild = 1234567890
		pfxv  = "%%"
	)

	c := testConfy()
	r := testRoute(t, nil, c)

	r.Prefix.Set(guild, pfxv)

	pfx, _ := r.Prefix.ForLine(guild, testMe, nil, pfxv)
	expect := route.Prefix{
		Value: pfxv,
		Clean: pfxv,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestForLineNull(t *testing.T) {
	t.Parallel()

	const pfxv = "test!"

	c := testConfy()
	r := testRoute(t, nil, c)

	const guild = 1234567890

	r.Prefix.Set(guild, "test!")

	pfx, _ := r.Prefix.ForLine(guild, testMe, nil, "test!uwu")
	expect := route.Prefix{
		Value: pfxv,
		Clean: pfxv,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestForLineUser(t *testing.T) {
	t.Parallel()

	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

	pfx, _ := r.Prefix.ForLine(discord.NullGuildID, testMe, nil, testMe.Mention())
	expect := route.Prefix{
		Value: testMe.Mention(),
		Clean: "@" + testMe.Username + " ",
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestForLineMemberNick(t *testing.T) {
	t.Parallel()

	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

	const guild = 666

	pfx, _ := r.Prefix.ForLine(guild, testMe, &testMMeNick, testMMeNick.Mention())
	expect := route.Prefix{
		Value: testMMeNick.Mention(),
		Clean: "@" + testMMeNick.Nick + " ",
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestForLineNil(t *testing.T) {
	t.Parallel()

	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

	const guild = 666

	pfx, ok := r.Prefix.ForLine(guild, testMe, nil, "")
	if ok {
		t.Errorf("should be !ok, got %#v", pfx)
	}

	m.Eval()
}
