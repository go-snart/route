package route_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/mavolin/dismock/v2/pkg/dismock"

	"github.com/go-snart/db"
	"github.com/go-snart/route"
)

const (
	testCat  = "test"
	testName = "cmd"
	testDesc = "lots of fun stuff"
)

var (
	testMe = discord.User{
		ID:       1234567890,
		Username: "User",
		Bot:      true,
	}

	testMMe = discord.Member{
		User: testMe,
	}

	testMMeNick = discord.Member{
		User: testMe,
		Nick: "Nick",
	}
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

	m.Me(testMe)

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			Author: testMe,
		},
	})

	m.Eval()
}

func TestHandleIgnoreBot(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	m.Me(testMe)

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

	m.Me(testMe)
	m.Member(guild, testMMe)

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

func TestHandleCommandNotFound(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	const guild = 123

	m.Me(testMe)
	m.Member(guild, testMMe)

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			GuildID: guild,
			Author: discord.User{
				ID: 999,
			},
			Content: testMMe.Mention(),
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

	m.Me(testMe)
	m.Member(guild, testMMe)

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			GuildID: guild,
			Author: discord.User{
				ID: 999,
			},
			Content: testMMe.Mention() + " " + c.Name,
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

	m.Me(testMe)
	m.Member(guild, testMMe)

	const testRun = "foo"

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			GuildID: guild,
			Author: discord.User{
				ID: 999,
			},
			Content: testMMe.Mention() + " " + c.Name + " -run=" + testRun,
		},
	})

	if *run != testRun {
		t.Errorf("expect %q\ngot %q", testRun, *run)
	}

	m.Eval()
}

func TestHandleMeError(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(testDB(), s)

	c, run := testCmd()

	r.Add(testCat, c)

	const guild = 123

	m.Error(
		http.MethodGet,
		"/users/@me",
		httputil.HTTPError{Status: 404},
	)

	const testRun = "foo"

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			GuildID: guild,
			Author: discord.User{
				ID: 999,
			},
			Content: testMMe.Mention() + " " + c.Name + " -run=" + testRun,
		},
	})

	if *run != "" {
		t.Errorf("expect %q\ngot %q", "", *run)
	}

	m.Eval()
}
