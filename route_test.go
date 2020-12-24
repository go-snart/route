package route_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/mavolin/dismock/pkg/dismock"

	"github.com/go-snart/db"
	"github.com/go-snart/route"
)

const (
	testCat  = "test"
	testName = "cmd"
	testDesc = "lots of fun stuff"
)

type testFlags struct {
	Run string `default:"run" usage:"run string"`
}

func testDB() *db.DB {
	const uri = "mem://"

	d, err := db.Open(uri)
	if err != nil {
		err = fmt.Errorf("open %q: %w", uri, err)
		panic(err)
	}

	return d
}

func testFunc(run *string) route.Func {
	return func(t *route.Trigger) error {
		*run = t.Flags.(testFlags).Run

		return nil
	}
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

	r.Add(testCat, c)

	cmds := r.Cats[testCat]

	if len(cmds) != 1 {
		t.Errorf("expect: 1\ngot %d", len(cmds))
	}

	if cmds[0] != c {
		t.Errorf("expect %v\ngot %v", c, cmds[0])
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

	r.Add(testCat, c)

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

	r.Add(testCat, c)

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
