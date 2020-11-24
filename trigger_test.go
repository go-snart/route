package route_test

import (
	"errors"
	"flag"
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/mavolin/dismock/pkg/dismock"

	"github.com/go-snart/route"
)

func TestTrigger(t *testing.T) {
	r := route.New(testDB(), nil)

	c, _ := testCmd()

	r.Add(c)

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

	r.Add(c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const line = "//"

	msg := discord.Message{
		Content: line,
	}

	_, err := r.Trigger(pfx, msg, line)
	if !errors.Is(err, route.ErrNoCmd) {
		t.Errorf("trigger %q %q", pfx.Clean, line)
	}
}

func TestTriggerErrNoTrigger(t *testing.T) {
	r := route.New(testDB(), nil)

	c, _ := testCmd()

	r.Add(c)

	pfx := &route.Prefix{
		Value: "//",
		Clean: "//",
	}

	const line = "//yeet"

	msg := discord.Message{
		Content: line,
	}

	_, err := r.Trigger(pfx, msg, line)
	if !errors.Is(err, route.ErrNoTrigger) {
		t.Errorf("trigger %q %q", pfx.Clean, line)
	}
}

func TestTriggerUsage(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, _ := testCmd()

	r.Add(c)

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
			Title:       "`cmd` Usage",
			Description: "lots of fun stuff",
			Fields: []discord.EmbedField{
				{
					Name:   "Flag `-run`",
					Value:  "run string\nDefault: `run`",
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
	c.Description = ""

	r.Add(c)

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
			Title:       "`cmd` Usage",
			Description: "*No description.*",
			Fields: []discord.EmbedField{
				{
					Name:   "Flag `-run`",
					Value:  "run string\nDefault: `run`",
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

	r.Add(c)

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

	r.Add(c)

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

	r.Add(c)

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

	r.Add(c)

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
			Title:       "`cmd` Usage",
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
