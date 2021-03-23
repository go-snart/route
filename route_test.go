package route_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/mavolin/dismock/v2/pkg/dismock"

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

func testFunc(run *string) route.Func {
	return func(t *route.Trigger) error {
		*run = t.Flags.(testFlags).Run

		return nil
	}
}

func TestNew(t *testing.T) {
	m, s := dismock.NewState(t)
	z := route.Store(nil)

	r := route.New(s, z)

	if r.State != s {
		t.Errorf("expect %v\ngot %v", s, r.State)
	}

	if r.Store != z {
		t.Errorf("expect %v\ngot %v", z, r.Store)
	}

	m.Eval()
}

func TestAdd(t *testing.T) {
	r := route.New(nil, nil)

	c, _ := testCmd()
	r.AddCmds(c)

	c2, _ := r.GetCmd(c.Name)

	if c2.Cat != c.Cat {
		t.Errorf("expect cat %#v, got %#v", c.Cat, c2.Cat)
	}

	if c2.Desc != c.Desc {
		t.Errorf("expect desc %#v, got %#v", c.Desc, c2.Desc)
	}

	if c2.Flags != c.Flags {
		t.Errorf("expect flags %#v, got %#v", c.Flags, c2.Flags)
	}

	if c2.Hide != c.Hide {
		t.Errorf("expect hide %#v, got %#v", c.Hide, c2.Hide)
	}

	if c2.Name != c.Name {
		t.Errorf("expect name %#v, got %#v", c.Name, c2.Name)
	}
}

func TestHandleIgnoreBot(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(s, nil)

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			Author: discord.User{
				Bot: true,
			},
		},
	})

	m.Eval()
}

func TestHandleIgnoreSelf(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(s, nil)

	m.Me(testMe)

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			Author: discord.User{
				ID: testMe.ID,
			},
		},
	})

	m.Eval()
}

func TestHandleNoPrefix(t *testing.T) {
	m, s := dismock.NewState(t)
	r := route.New(s, nil)

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
	r := route.New(s, nil)

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
	r := route.New(s, nil)

	c, _ := testCmd()
	c.Func = func(*route.Trigger) error {
		return io.EOF
	}

	r.AddCmds(c)

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
	r := route.New(s, nil)

	c, run := testCmd()

	r.AddCmds(c)

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
	r := route.New(s, nil)

	c, run := testCmd()

	r.AddCmds(c)

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
