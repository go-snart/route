package route_test

import (
	"io"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/dismock/pkg/dismock"

	"github.com/go-snart/route"
)

type testFlags struct {
	Run string `default:"run" usage:"run string"`
}

func testCmd() (*route.Command, *string) {
	run := ""

	return &route.Command{
		Name:        "cmd",
		Category:    "misc",
		Description: "lots of fun stuff",

		Func: func(t *route.Trigger) error {
			run = t.Flags.(testFlags).Run

			return nil
		},

		Hidden: false,

		Flags: testFlags{},
	}, &run
}

func TestNew(t *testing.T) {
	m, s := dismock.NewState(t)
	d := testDB()
	r := route.New(d, s)

	if r.State != s {
		t.Errorf("expect %v\ngot %v", s, r.State)
	}

	if r.DB != d {
		t.Errorf("expect %v\ngot %v", d, r.DB)
	}

	m.Eval()
}

func TestAdd(t *testing.T) {
	r := route.New(testDB(), nil)

	c, _ := testCmd()

	r.Add(c)

	if len(r.Commands) != 1 {
		t.Errorf("expect: 1\ngot %d", len(r.Commands))
	}

	if r.Commands[0] != c {
		t.Errorf("expect %v\ngot %v", c, r.Commands[0])
	}
}

func TestHandleIgnoreSelf(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	me := r.State.Ready.User

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			Author: me,
		},
	})

	m.Eval()
}

func TestHandleIgnoreBot(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			Author: discord.User{
				ID:  999,
				Bot: true,
			},
		},
	})

	m.Eval()
}

func TestHandleNoPrefix(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	const guild = 123

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
	}

	m.Member(guild, mme)

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			GuildID: guild,
			Author: discord.User{
				ID: 999,
			},
			Content: "yeet",
		},
	})

	m.Eval()
}

func TestHandleNoTrigger(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	const guild = 123

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
	}

	m.Member(guild, mme)

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			GuildID: guild,
			Author: discord.User{
				ID: 999,
			},
			Content: mme.Mention(),
		},
	})

	m.Eval()
}

func TestHandleRunError(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, _ := testCmd()
	c.Func = func(*route.Trigger) error {
		return io.EOF
	}

	r.Add(c)

	const guild = 123

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
	}

	m.Member(guild, mme)

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			GuildID: guild,
			Author: discord.User{
				ID: 999,
			},
			Content: mme.Mention() + " " + c.Name,
		},
	})

	m.Eval()
}

func TestHandle(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, run := testCmd()

	r.Add(c)

	const guild = 123

	me := r.State.Ready.User
	mme := discord.Member{
		User: me,
	}

	m.Member(guild, mme)

	const testRun = "foo"

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			GuildID: guild,
			Author: discord.User{
				ID: 999,
			},
			Content: mme.Mention() + " " + c.Name + " -run=" + testRun,
		},
	})

	if *run != testRun {
		t.Errorf("expect %q\ngot %q", testRun, *run)
	}

	m.Eval()
}
