package route_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/httputil"
	"github.com/mavolin/dismock/pkg/dismock"

	"github.com/go-snart/route"
)

func TestFindPrefixGuild(t *testing.T) {
	r := route.New(testDB(), nil)

	const (
		guild = 1234567890
		pfxv  = "//"
	)

	key := fmt.Sprintf("prefix_%d", guild)

	err := r.DB.Set(key, pfxv)
	if err != nil {
		t.Errorf("db set %q %q: %s", key, pfxv, err)
	}

	pfx := r.FindPrefix(guild, pfxv)
	expect := &route.Prefix{
		Value: pfxv,
		Clean: pfxv,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestFindPrefixDefault(t *testing.T) {
	r := route.New(testDB(), nil)

	const (
		pfxv = "//"
		key  = "prefix"
	)

	err := r.DB.Set(key, pfxv)
	if err != nil {
		t.Errorf("db set %q %q: %s", key, pfxv, err)
	}

	pfx := r.FindPrefix(123, pfxv)
	expect := &route.Prefix{
		Value: pfxv,
		Clean: pfxv,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestFindPrefixUser(t *testing.T) {
	_, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	me := r.State.Ready.User

	pfx := r.FindPrefix(0, me.Mention())
	expect := &route.Prefix{
		Value: me.Mention(),
		Clean: "@" + me.Username,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestFindPrefixMember(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
	}

	const guild = 666

	m.Member(guild, mme)

	pfx := r.FindPrefix(guild, mme.Mention())
	expect := &route.Prefix{
		Value: mme.Mention(),
		Clean: "@" + me.Username,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestFindPrefixMemberNick(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	const nick = "foo"

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
		Nick: nick,
	}

	const guild = 666

	m.Member(guild, mme)

	pfx := r.FindPrefix(guild, mme.Mention())
	expect := &route.Prefix{
		Value: mme.Mention(),
		Clean: "@" + nick,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestFindPrefixMemberErr(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
	}

	const guild = 666

	m.Error(
		"GET",
		fmt.Sprintf("/guilds/%d/members", guild),
		httputil.HTTPError{Status: 404},
	)

	pfx := r.FindPrefix(guild, mme.Mention())
	expect := (*route.Prefix)(nil)

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestFindPrefixNil(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
	}

	const guild = 666

	m.Member(guild, mme)

	pfx := r.FindPrefix(guild, "")
	expect := (*route.Prefix)(nil)

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}
