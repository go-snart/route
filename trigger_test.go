package route_test

import (
	"errors"
	"flag"
	"strings"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

var testPfx = &route.Prefix{
	Value: "//",
	Clean: "//",
}

func TestTrigger(t *testing.T) {
	r := route.New(testSettings, nil)

	c, _ := testCmd()

	r.Add(c)

	const line = "//cmd `-run=foo`"

	msg := discord.Message{
		Content: line,
	}

	tr, err := r.Trigger(testPfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q: %s", testPfx.Clean, line, err)
	}

	const expect = "foo"

	if flags := tr.Flags.(testFlags); flags.Run != expect {
		t.Errorf("expect run %q, got %q", expect, flags.Run)
	}
}

func TestTriggerErrNoCmd(t *testing.T) {
	r := route.New(testSettings, nil)

	c, _ := testCmd()

	r.Add(c)

	const line = "//"

	msg := discord.Message{
		Content: line,
	}

	_, err := r.Trigger(testPfx, msg, line)
	if !errors.Is(err, route.ErrNoCommand) {
		t.Errorf("trigger %q %q: %w", testPfx.Clean, line, err)
	}
}

func TestTriggerErrCommandNotFound(t *testing.T) {
	r := route.New(testSettings, nil)

	c, _ := testCmd()
	r.Add(c)

	const line = "//yeet"

	msg := discord.Message{
		Content: line,
	}

	_, err := r.Trigger(testPfx, msg, line)
	if !errors.Is(err, route.ErrCommandNotFound) {
		t.Errorf("trigger %q %q: %w", testPfx.Clean, line, err)
	}
}

func TestTriggerUsage(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	c, _ := testCmd()
	r.Add(c)

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

	_, err := r.Trigger(testPfx, msg, line)
	if !errors.Is(err, flag.ErrHelp) {
		t.Errorf("trigger %q %q: %s", testPfx.Clean, line, err)
	}

	m.Eval()
}

func TestTriggerUsageNoDesc(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	c, _ := testCmd()
	c.Desc = ""
	r.Add(c)

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

	_, err := r.Trigger(testPfx, msg, line)
	if !errors.Is(err, flag.ErrHelp) {
		t.Errorf("trigger %q %q: %s", testPfx.Clean, line, err)
	}

	m.Eval()
}

func TestTriggerBadFlags(t *testing.T) {
	r := route.New(testSettings, nil)

	c, _ := testCmd()
	c.Flags = (chan int)(nil)
	r.Add(c)

	const (
		channel = 1234567890
		line    = "//cmd `-help`"
	)

	msg := discord.Message{
		ChannelID: channel,
		Content:   line,
	}

	_, err := r.Trigger(testPfx, msg, line)
	if err == nil {
		t.Errorf("trigger %q %q: %s", testPfx.Clean, line, err)
	}
}

func TestReplySendErr(t *testing.T) {
	_, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	c, _ := testCmd()
	r.Add(c)

	const (
		channel = 1234567890
		line    = "//cmd"
	)

	msg := discord.Message{
		ChannelID: channel,
		Content:   line,
	}

	tr, err := r.Trigger(testPfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q: %s", testPfx.Clean, line, err)
	}

	rep := tr.Reply()

	err = rep.Send()
	if err == nil {
		t.Errorf("send: %s", err)
	}
}

func TestTriggerRun(t *testing.T) {
	r := route.New(testSettings, nil)

	c, run := testCmd()
	r.Add(c)

	const (
		erun    = "foo"
		channel = 1234567890
		line    = "//cmd `-run=" + erun + "`"
	)

	msg := discord.Message{
		ChannelID: channel,
		Content:   line,
	}

	tr, err := r.Trigger(testPfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q: %s", testPfx.Clean, line, err)
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
	r := route.New(testSettings, s)

	c, _ := testCmd()
	c.Flags = nil
	r.Add(c)

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

	_, err := r.Trigger(testPfx, msg, line)
	if err == nil {
		t.Errorf("trigger %q %q: %#v", testPfx.Clean, line, err)
	}

	m.Eval()
}

func TestTriggerUsageFillError(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testSettings, s)

	c, _ := testCmd()
	c.Flags = (func())(nil)

	const channel = 1234567890

	(&route.Trigger{
		Router:  r,
		Command: c,

		Message: discord.Message{
			ChannelID: channel,
		},
		Prefix: testPfx,
		Output: &strings.Builder{},
	}).Usage()

	m.Eval()
}
