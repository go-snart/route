package route

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"
	re2 "github.com/dlclark/regexp2"
)

var splitter = re2.MustCompile(`(\x60+)(.*?)\1|(\S+)`, 0)

// Trigger holds the context that triggered a Command.
type Trigger struct {
	Route   *Route
	Message discord.Message
	Prefix  *Prefix
	Command *Command
	FlagSet *flag.FlagSet
	Args    []string
	Flags   interface{}
	Output  *strings.Builder
}

// Trigger gets a Trigger by finding an appropriate Command for a given prefix, session, message, etc.
func (r *Route) Trigger(pfx *Prefix, m discord.Message, line string) (*Trigger, error) {
	t := &Trigger{
		Route:   r,
		Message: m,
		Prefix:  pfx,
		Output:  &strings.Builder{},
	}

	line = strings.TrimSpace(strings.TrimPrefix(line, pfx.Value))

	args := split(line)

	if len(args) == 0 {
		return nil, ErrNoCmd
	}

	cmd := args[0]
	args = args[1:]

	for _, c := range r.Commands {
		if c.Name == cmd {
			t.Command = c

			break
		}
	}

	if t.Command == nil {
		return nil, ErrNoTrigger
	}

	t.FlagSet = flag.NewFlagSet(t.Command.Name, flag.ContinueOnError)
	t.FlagSet.SetOutput(t.Output)
	t.FlagSet.Usage = t.Usage

	if t.Command.Flags == nil {
		t.Command.Flags = struct{}{}
	}

	flags := reflect.New(reflect.TypeOf(t.Command.Flags))

	err := r.Filler.Fill(t.FlagSet, flags.Interface())
	if err != nil {
		return nil, fmt.Errorf("fill: %w", err)
	}

	err = t.FlagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	t.Flags = flags.Elem().Interface()
	t.Args = t.FlagSet.Args()

	return t, nil
}

// Usage is the help flag handler for the Trigger.
func (t *Trigger) Usage() {
	rep := t.Reply()

	if t.Output.Len() > 0 {
		rep.Content = t.Output.String()
	}

	desc := t.Command.Description
	if desc == "" {
		desc = "*no description*"
	}

	rep.Embed = &discord.Embed{
		Title:       "`" + t.Command.Name + "` usage",
		Description: desc,
	}

	t.FlagSet.VisitAll(func(f *flag.Flag) {
		rep.Embed.Fields = append(
			rep.Embed.Fields, discord.EmbedField{
				Name:   "flag `-" + f.Name + "`",
				Value:  f.Usage + "\ndefault: `" + f.DefValue + "`",
				Inline: false,
			},
		)
	})

	// fuck it, no error check
	_ = rep.Send()
}

// Run is a shortcut to t.Command.Func(t).
func (t *Trigger) Run() error {
	return t.Command.Func(t)
}

// Reply wraps a message to be sent to a given ChannelID using a given Session.
type Reply struct {
	Trigger *Trigger

	api.SendMessageData
}

// Reply gets a Reply for the Trigger.
func (t *Trigger) Reply() *Reply {
	return &Reply{
		Trigger: t,

		SendMessageData: api.SendMessageData{},
	}
}

// SendMsg sends the Reply.
func (r *Reply) SendMsg() (*discord.Message, error) {
	return r.Trigger.Route.State.SendMessageComplex(r.Trigger.Message.ChannelID, r.SendMessageData)
}

// Send is a shortcut for SendMsg that logs a warning on error and elides the resulting Message.
func (r *Reply) Send() error {
	_, err := r.SendMsg()
	if err != nil {
		return fmt.Errorf("send msg: %w", err)
	}

	return nil
}

func split(s string) []string {
	subj := []rune(s)
	args := []string(nil)

	for {
		// will only error if a timeout is set (it isn't)
		m, _ := splitter.FindRunesMatch(subj)
		if m == nil {
			break
		}

		gs := m.Groups()

		match := gs[3].Capture.String()
		if match == "" {
			match = gs[2].Capture.String()
		}

		args = append(args, match)

		l := gs[0].Capture.Length + 1
		if l > len(subj) {
			break
		}

		subj = subj[l:]
	}

	return args
}
