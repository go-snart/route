package route

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	re2 "github.com/dlclark/regexp2"
	ff "github.com/itzg/go-flagsfiller"
)

//nolint:gochecknoglobals // pre-compiling regexp
var splitter = re2.MustCompile(`(\x60+)(.*?)\1|(\S+)`, 0)

//nolint:gochecknoglobals // only needs to be set up once
var filler = ff.New()

var (
	// ErrCmdNotFound occurs when an unknown command is called.
	ErrCmdNotFound = errors.New("command not found")

	// ErrNoCmd occurs when there is no command after the prefix.
	ErrNoCmd = errors.New("no command")
)

// Func is a handler for a Trigger.
type Func = func(*Trigger) error

// Trigger holds the context that triggered a Command.
type Trigger struct {
	Route   *Route
	Message discord.Message
	Prefix  Prefix
	Command Cmd
	FlagSet *flag.FlagSet
	Args    []string
	Flags   interface{}
	Output  *strings.Builder
}

// Trigger gets a Trigger by finding an appropriate Command for a given prefix, message, and line.
func (r *Route) Trigger(pfx Prefix, m discord.Message, line string) (*Trigger, error) {
	//nolint:exhaustivestruct // this stuff will be filled
	t := &Trigger{
		Route:   r,
		Message: m,
		Prefix:  pfx,
		Output:  &strings.Builder{},
	}

	line = strings.TrimSpace(strings.TrimPrefix(line, pfx.Value))
	if len(line) == 0 {
		return nil, ErrNoCmd
	}

	name, args := split(line)

	cmd, ok := r.Cmd.Get(name)
	if !ok {
		return nil, ErrCmdNotFound
	}

	t.Command = cmd

	flags, err := t.fillFlagSet()
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

func (t *Trigger) fillFlagSet() (reflect.Value, error) {
	t.FlagSet = flag.NewFlagSet(t.Command.Name, flag.ContinueOnError)
	t.FlagSet.SetOutput(t.Output)
	t.FlagSet.Usage = t.Usage

	flags := reflect.New(reflect.TypeOf(t.Command.Flags))

	err := filler.Fill(t.FlagSet, flags.Interface())
	if err != nil {
		return reflect.ValueOf(nil), fmt.Errorf("fill flags: %w", err)
	}

	return flags, nil
}

// Usage is the help flag handler for the Trigger.
func (t *Trigger) Usage() {
	rep := t.Reply()
	rep.Content = t.Output.String()

	//nolint:exhaustivestruct
	// discord types are excessive
	// Fields can be nil for append
	rep.Embed = &discord.Embed{
		Title:       fmt.Sprintf("`%s` usage", t.Command.Name),
		Description: t.Command.Desc,
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

// Reply wraps a message to be sent to a given ChannelID using a given Session.
type Reply struct {
	api.SendMessageData

	Trigger *Trigger
}

// Reply gets a Reply for the Trigger.
func (t *Trigger) Reply() *Reply {
	return &Reply{
		//nolint:exhaustivestruct
		// discord types are excessive
		SendMessageData: api.SendMessageData{},

		Trigger: t,
	}
}

// SendMsg sends the Reply.
func (r *Reply) SendMsg() (*discord.Message, error) {
	msg, err := r.Trigger.Route.State.SendMessageComplex(r.Trigger.Message.ChannelID, r.SendMessageData)
	if err != nil {
		return nil, fmt.Errorf("reply send: %w", err)
	}

	return msg, nil
}

// Send is a shortcut for SendMsg that elides the resulting message.
func (r *Reply) Send() error {
	_, err := r.SendMsg()

	return err
}

func split(s string) (string, []string) {
	subj := []rune(s)
	args := []string(nil)

	for {
		// will only error if a timeout is set (it isn't)
		m, _ := splitter.FindRunesMatch(subj)
		gs := m.Groups()

		match := gs[3].Capture.String()
		if len(match) == 0 {
			match = gs[2].Capture.String()
		}

		args = append(args, match)

		l := gs[0].Capture.Length + 1
		if l > len(subj) {
			break
		}

		subj = subj[l:]
	}

	return args[0], args[1:]
}
