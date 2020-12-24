package route_test

import (
	"reflect"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/go-snart/route"
	"github.com/mavolin/dismock/pkg/dismock"
)

func testCmd() (*route.Command, *string) {
	run := ""

	return &route.Command{
		Name:  testName,
		Desc:  testDesc,
		Func:  testFunc(&run),
		Hide:  false,
		Flags: testFlags{},
	}, &run
}

func TestTidyDesc(t *testing.T) {
	c, _ := testCmd()
	c.Desc = ""

	c.Tidy()

	if c.Desc != route.DefaultDesc {
		t.Errorf("expect %q, got %q", route.DefaultDesc, c.Desc)
	}
}

func TestTidyFunc(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(nil, s)

	c, _ := testCmd()
	c.Func = nil

	c.Tidy()

	if c.Func == nil {
		t.Error("shouldn't be nil")
	}

	const channel = 1234567890

	m.SendMessage(nil, discord.Message{
		ChannelID: channel,
		Content:   route.UndefinedMsg,
	})

	err := c.Func(&route.Trigger{
		Route: r,
		Message: discord.Message{
			ChannelID: channel,
		},
	})
	if err != nil {
		t.Errorf("run func: %s", err)
	}

	m.Eval()
}

func TestTidyFlags(t *testing.T) {
	c, _ := testCmd()
	c.Flags = nil

	c.Tidy()

	if !reflect.DeepEqual(c.Flags, struct{}{}) {
		t.Errorf("expect %#v, got %#v", struct{}{}, c.Flags)
	}
}
