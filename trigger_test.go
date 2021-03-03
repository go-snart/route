package route_test

import (
	"errors"
	"flag"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/route"
)

var testPfx = route.Prefix{
	Value: "//",
	Clean: "//",
}

func TestTrigger(t *testing.T) {
	r := route.New(nil)

	c, _ := testCmd()

	r.AddCmds(c)

	const line = "//cmd `-run=foo`"

	msg := discord.Message{
		Content: line,
	}

	tr, err := r.Trigger(testPfx, msg, line)
	if err != nil {
		t.Errorf("trigger %q %q: %s", testPfx.Clean, line, err)
	}

	const expect = "foo"
	if flags, _ := tr.Flags.(testFlags); flags.Run != expect {
		t.Errorf("expect run %q, got %q", expect, flags.Run)
	}
}

func TestTriggerErrNoCmd(t *testing.T) {
	r := route.New(nil)

	c, _ := testCmd()

	r.AddCmds(c)

	const line = "//"

	msg := discord.Message{
		Content: line,
	}

	_, err := r.Trigger(testPfx, msg, line)
	if !errors.Is(err, route.ErrNoCmd) {
		t.Errorf("trigger %q %q: %w", testPfx.Clean, line, err)
	}
}

func TestTriggerErrCmdNotFound(t *testing.T) {
	r := route.New(nil)

	c, _ := testCmd()
	r.AddCmds(c)

	const line = "//yeet"

	msg := discord.Message{
		Content: line,
	}

	_, err := r.Trigger(testPfx, msg, line)
	if !errors.Is(err, route.ErrCmdNotFound) {
		t.Errorf("trigger %q %q: %w", testPfx.Clean, line, err)
	}
}

func TestTriggerUsage(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(s)

	c, _ := testCmd()
	r.AddCmds(c)

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

func TestReplySendErr(t *testing.T) {
	_, s := dismock.NewState(t)
	r := route.New(s)

	c, _ := testCmd()
	r.AddCmds(c)

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
	r := route.New(nil)

	c, run := testCmd()
	r.AddCmds(c)

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

	err = tr.Command.Func(tr)
	if err != nil {
		t.Errorf("run: %s", err)
	}

	if *run != erun {
		t.Errorf("expect %q\ngot %q", erun, *run)
	}
}

func TestTriggerFillError(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(s)

	c, _ := testCmd()
	c.Flags = (func())(nil)
	r.AddCmds(c)

	const channel = 1234567890

	line := testPfx.Value + c.Name

	_, err := r.Trigger(testPfx, discord.Message{
		ChannelID: channel,
	}, line)
	if err == nil {
		t.Errorf("expect err, got nil")
	}

	m.Eval()
}
