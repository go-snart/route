package route_test

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

func TestTrigger(t *testing.T) {
	r := route.New(testDB(), nil)

	c, _ := testCmd()

	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const line = "//cmd `-run=foo`"

	msg := discord.Message{
		Content: line,
	}

	tr, err := r.Trigger(pfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q", pfx.Clean, line)
	}

	expect := &route.Trigger{
		Route:   r,
		Message: msg,
		Prefix:  pfx,
		Command: c,
		FlagSet: tr.FlagSet, // probably shouldn't do this
		Args:    []string{},
		Flags: testFlags{
			Run: "foo",
		},
		Output: tr.Output, // probably shouldn't do this
	}

	if !reflect.DeepEqual(tr, expect) {
		t.Errorf("\nexpect %#v\ngot %#v", expect, tr)
	}
}

func TestTriggerErrNoCmd(t *testing.T) {
	r := route.New(testDB(), nil)

	c, _ := testCmd()

	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const line = "//"

	msg := discord.Message{
		Content: line,
	}

	_, err := r.Trigger(pfx, msg, line)
	if !errors.Is(err, route.ErrNoCommand) {
		t.Errorf("trigger %q %q: %w", pfx.Clean, line, err)
	}
}

func TestTriggerErrCommandNotFound(t *testing.T) {
	r := route.New(testDB(), nil)

	c, _ := testCmd()

	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const line = "//yeet"

	msg := discord.Message{
		Content: line,
	}

	_, err := r.Trigger(pfx, msg, line)
	if !errors.Is(err, route.ErrCommandNotFound) {
		t.Errorf("trigger %q %q: %w", pfx.Clean, line, err)
	}
}

func TestTriggerUsage(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, _ := testCmd()

	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const (
		channel = 1234567890
		line    = "//cmd `-help`"
	)

	msg := discord.Message{
		ChannelID: channel,
		Content:   line,
	}

	m.SendMessage(
		&discord.Embed{
			Title:       "`cmd` usage",
			Description: "lots of fun stuff",
			Fields: []discord.EmbedField{
				{
					Name:   "flag `-run`",
					Value:  "run string\ndefault: `run`",
					Inline: false,
				},
			},
		},
		discord.Message{
			ChannelID: channel,
			Content:   "",
		},
	)

	_, err := r.Trigger(pfx, msg, line)
	if !errors.Is(err, flag.ErrHelp) {
		t.Errorf("trigger %q %q: %s", pfx.Clean, line, err)
	}

	m.Eval()
}

func TestTriggerUsageNoDesc(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, _ := testCmd()
	c.Desc = ""

	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const (
		channel = 1234567890
		line    = "//cmd `-help`"
	)

	msg := discord.Message{
		ChannelID: channel,
		Content:   line,
	}

	m.SendMessage(
		&discord.Embed{
			Title:       "`cmd` usage",
			Description: "*no description*",
			Fields: []discord.EmbedField{
				{
					Name:   "flag `-run`",
					Value:  "run string\ndefault: `run`",
					Inline: false,
				},
			},
		},
		discord.Message{
			ChannelID: channel,
			Content:   "",
		},
	)

	_, err := r.Trigger(pfx, msg, line)
	if !errors.Is(err, flag.ErrHelp) {
		t.Errorf("trigger %q %q: %s", pfx.Clean, line, err)
	}

	m.Eval()
}

func TestTriggerBadFlags(t *testing.T) {
	r := route.New(testDB(), nil)

	c, _ := testCmd()
	c.Flags = (chan int)(nil)

	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const (
		channel = 1234567890
		line    = "//cmd `-help`"
	)

	msg := discord.Message{
		ChannelID: channel,
		Content:   line,
	}

	_, err := r.Trigger(pfx, msg, line)
	if err == nil {
		t.Errorf("trigger %q %q: %s", pfx.Clean, line, err)
	}
}

func TestReplySendErr(t *testing.T) {
	_, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, _ := testCmd()

	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const (
		channel = 1234567890
		line    = "//cmd"
	)

	msg := discord.Message{
		ChannelID: channel,
		Content:   line,
	}

	tr, err := r.Trigger(pfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q: %s", pfx.Clean, line, err)
	}

	rep := tr.Reply()

	err = rep.Send()
	if err == nil {
		t.Errorf("send: %s", err)
	}
}

func TestTriggerRun(t *testing.T) {
	r := route.New(testDB(), nil)

	c, run := testCmd()

	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const (
		erun    = "foo"
		channel = 1234567890
		line    = "//cmd `-run=" + erun + "`"
	)

	msg := discord.Message{
		ChannelID: channel,
		Content:   line,
	}

	tr, err := r.Trigger(pfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q: %s", pfx.Clean, line, err)
	}

	err = tr.Run()
	if err != nil {
		t.Errorf("run: %s", err)
	}

	if *run != erun {
		t.Errorf("expect %q\ngot %q", erun, *run)
	}
}

func TestTriggerNilFlags(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, _ := testCmd()
	c.Flags = nil

	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const (
		channel = 123456790
		line    = "//cmd `-run=foo`"
	)

	msg := discord.Message{
		ChannelID: channel,
		Content:   line,
	}

	m.SendMessage(
		&discord.Embed{
			Title:       "`cmd` usage",
			Description: "lots of fun stuff",
		},
		discord.Message{
			ChannelID: channel,
			Content:   "flag provided but not defined: -run\n",
		},
	)

	_, err := r.Trigger(pfx, msg, line)
	if err == nil {
		t.Errorf("trigger %q %q: %#v", pfx.Clean, line, err)
	}

	m.Eval()
}

func TestTriggerUsageFillError(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(nil, s)

	c, _ := testCmd()
	c.Flags = (func())(nil)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const channel = 1234567890

	(&route.Trigger{
		Route:   r,
		Command: c,

		Message: discord.Message{
			ChannelID: channel,
		},
		Prefix: pfx,
		Output: &strings.Builder{},
	}).Usage()

	m.Eval()
}

func TestGetMMe(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, _ := testCmd()
	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const (
		guild   = 1234567890
		channel = 1234567890
		line    = "//cmd"
	)

	m.Me(testMe)
	m.Member(guild, testMMe)

	msg := discord.Message{
		GuildID:   guild,
		ChannelID: channel,
		Content:   line,
	}

	tr, err := r.Trigger(pfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q: %s", pfx.Clean, line, err)
	}

	mme, err := tr.GetMMe()
	if err != nil {
		t.Errorf("get mme: %s", err)
	}

	if !reflect.DeepEqual(testMMe, *mme) {
		t.Errorf("expect %#v\ngot %#v", testMMe, *mme)
	}

	m.Eval()
}

func TestGetMMeExists(t *testing.T) {
	tr := &route.Trigger{
		MMe: &testMMe,
	}

	mme, err := tr.GetMMe()
	if err != nil {
		t.Errorf("get mme: %s", err)
	}

	if !reflect.DeepEqual(testMMe, *mme) {
		t.Errorf("expect %#v\ngot %#v", testMMe, *mme)
	}
}

func TestGetMMeNullGuild(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, _ := testCmd()
	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const (
		guild   = discord.NullGuildID
		channel = 1234567890
		line    = "//cmd"
	)

	msg := discord.Message{
		GuildID:   guild,
		ChannelID: channel,
		Content:   line,
	}

	tr, err := r.Trigger(pfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q: %s", pfx.Clean, line, err)
	}

	_, err = tr.GetMMe()
	if !errors.Is(err, route.ErrNullGuild) {
		t.Errorf("get mme: %s", err)
	}

	m.Eval()
}

func TestGetMMeErrorMe(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, _ := testCmd()
	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const (
		guild   = 1234567890
		channel = 1234567890
		line    = "//cmd"
	)

	m.Error(http.MethodGet, "/users/@me", httputil.HTTPError{Status: 404})

	msg := discord.Message{
		GuildID:   guild,
		ChannelID: channel,
		Content:   line,
	}

	tr, err := r.Trigger(pfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q: %s", pfx.Clean, line, err)
	}

	_, err = tr.GetMMe()
	if err == nil {
		t.Errorf("get mme: %s", err)
	}

	m.Eval()
}

func TestGetMMeErrorMember(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, _ := testCmd()
	r.Add(testCat, c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const (
		guild   = 1234567890
		channel = 1234567890
		line    = "//cmd"
	)

	m.Me(testMe)
	m.Error(http.MethodGet,
		fmt.Sprintf("/guilds/%d/members/%d", guild, testMe.ID),
		httputil.HTTPError{Status: 404})

	msg := discord.Message{
		GuildID:   guild,
		ChannelID: channel,
		Content:   line,
	}

	tr, err := r.Trigger(pfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q: %s", pfx.Clean, line, err)
	}

	_, err = tr.GetMMe()
	if err == nil {
		t.Errorf("get mme: %s", err)
	}

	m.Eval()
}
