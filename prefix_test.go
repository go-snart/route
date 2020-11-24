package route_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/utils/httputil"
	"github.com/mavolin/dismock/pkg/dismock"

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
	_, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	me := r.State.Ready.User

	pfx := r.LinePrefix(0, me.Mention())
	expect := &route.Prefix{
		Value: me.Mention(),
		Clean: "@" + me.Username,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}
}

func TestLinePrefixMember(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
	}

	const guild = 666

	m.Member(guild, mme)

	pfx := r.LinePrefix(guild, mme.Mention())
	expect := &route.Prefix{
		Value: mme.Mention(),
		Clean: "@" + me.Username,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixMemberNick(t *testing.T) {
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

	pfx := r.LinePrefix(guild, mme.Mention())
	expect := &route.Prefix{
		Value: mme.Mention(),
		Clean: "@" + nick,
	}

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixMemberErr(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
	}

	const guild = 666

	m.Error(
		http.MethodGet,
		fmt.Sprintf("/guilds/%d/members", guild),
		httputil.HTTPError{Status: 404},
	)

	pfx := r.LinePrefix(guild, mme.Mention())
	expect := (*route.Prefix)(nil)

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}

func TestLinePrefixNil(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
	}

	const guild = 666

	m.Member(guild, mme)

	pfx := r.LinePrefix(guild, "")
	expect := (*route.Prefix)(nil)

	if !reflect.DeepEqual(pfx, expect) {
		t.Errorf("expect: %v\ngot: %v", expect, pfx)
	}

	m.Eval()
}
