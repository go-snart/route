package route_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/state"
	"github.com/diamondburned/arikawa/v2/utils/httputil"
	"github.com/mavolin/dismock/v2/pkg/dismock"
	"github.com/superloach/confy"

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

func testConfy() confy.Confy {
	c := confy.NewMem()

	err := c.Set(route.KeyPrefix, map[discord.GuildID]string{
		discord.NullGuildID: testPfx.Value,
	})
	if err != nil {
		// shouldn't error
		panic(err)
	}

	return c
}

func testRoute(t *testing.T, s *state.State, c confy.Confy) *route.Route {
	t.Helper()

	r, err := route.New(s, c)
	if err != nil {
		t.Errorf("new route: %s", err)
	}

	return r
}

func TestNew(t *testing.T) {
	t.Parallel()

	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

	if r.State != s {
		t.Errorf("expect %v\ngot %v", s, r.State)
	}

	if r.Confy != c {
		t.Errorf("expect %v\ngot %v", c, r.Confy)
	}

	m.Eval()
}

func TestAdd(t *testing.T) {
	t.Parallel()
	r := testRoute(t, nil, testConfy())

	cmd, _ := testCmd()
	r.Cmd.Add(cmd)

	cmd2, _ := r.Cmd.Get(cmd.Name)

	if cmd2.Cat != cmd.Cat {
		t.Errorf("expect cat %#v, got %#v", cmd.Cat, cmd2.Cat)
	}

	if cmd2.Desc != cmd.Desc {
		t.Errorf("expect desc %#v, got %#v", cmd.Desc, cmd2.Desc)
	}

	if cmd2.Flags != cmd.Flags {
		t.Errorf("expect flags %#v, got %#v", cmd.Flags, cmd2.Flags)
	}

	if cmd2.Hide != cmd.Hide {
		t.Errorf("expect hide %#v, got %#v", cmd.Hide, cmd2.Hide)
	}

	if cmd2.Name != cmd.Name {
		t.Errorf("expect name %#v, got %#v", cmd.Name, cmd2.Name)
	}
}

func TestHandleIgnoreBot(t *testing.T) {
	t.Parallel()
	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

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
	t.Parallel()
	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

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
	t.Parallel()
	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

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
	t.Parallel()
	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

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
	t.Parallel()
	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

	cmd, _ := testCmd()
	cmd.Func = func(*route.Trigger) error {
		return io.EOF
	}

	r.Cmd.Add(cmd)

	const guild = 123

	m.Me(testMe)
	m.Member(guild, testMMe)

	r.Handle(&gateway.MessageCreateEvent{
		Message: discord.Message{
			GuildID: guild,
			Author: discord.User{
				ID: 999,
			},
			Content: testMMe.Mention() + " " + cmd.Name,
		},
	})

	m.Eval()
}

func TestHandle(t *testing.T) {
	t.Parallel()
	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

	cmd, run := testCmd()
	r.Cmd.Add(cmd)

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
			Content: testMMe.Mention() + " " + cmd.Name + " -run=" + testRun,
		},
	})

	if *run != testRun {
		t.Errorf("expect %q\ngot %q", testRun, *run)
	}

	m.Eval()
}

func TestHandleMeError(t *testing.T) {
	t.Parallel()
	m, s := dismock.NewState(t)
	c := testConfy()
	r := testRoute(t, s, c)

	cmd, run := testCmd()
	r.Cmd.Add(cmd)

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
			Content: testMMe.Mention() + " " + cmd.Name + " -run=" + testRun,
		},
	})

	if *run != "" {
		t.Errorf("expect %q\ngot %q", "", *run)
	}

	m.Eval()
}
