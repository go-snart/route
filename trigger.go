package route

import (
	"flag"
	"fmt"
	"reflect"
	"strings"

	"github.com/diamondburned/arikawa/api"
	"github.com/diamondburned/arikawa/discord"

	"github.com/go-snart/logs"
)

// Trigger holds the context that triggered a Command.
type Trigger struct {
	Route   *Route
	Message discord.Message
	Prefix  *Prefix
	Command *Command
	FlagSet *flag.FlagSet
	Args    []string
	Flags   interface{}
}

// Trigger gets a Trigger by finding an appropriate Command for a given prefix, session, message, etc.
func (r *Route) Trigger(pfx *Prefix, m discord.Message, line string) (*Trigger, error) {
	t := &Trigger{
		Route:   r,
		Message: m,
		Prefix:  pfx,
		Command: nil,
		Args:    nil,
		Flags:   nil,
	}

	logs.Debug.Println("line", line)

	line = strings.TrimSpace(strings.TrimPrefix(line, pfx.Value))

	logs.Debug.Println("line", line)

	args := Split(line)

	logs.Debug.Println("args", args)

	if len(args) == 0 {
		logs.Debug.Println("0 args")

		return nil, ErrNoCmd
	}

	cmd := args[0]
	logs.Debug.Println("cmd", cmd)

	args = args[1:]
	logs.Debug.Println("args", args)

	for _, c := range r.Commands {
		if c.Name == cmd {
			t.Command = c

			break
		}
	}

	logs.Debug.Println("cmd", t.Command)

	if t.Command == nil {
		return nil, ErrNoTrigger
	}

	t.FlagSet = flag.NewFlagSet(t.Command.Name, flag.ContinueOnError)
	t.FlagSet.Usage = t.Usage

	flags := reflect.New(reflect.TypeOf(t.Command.Flags))

	err := r.Filler.Fill(t.FlagSet, flags.Interface())
	if err != nil {
		return nil, fmt.Errorf("fill: %w", err)
	}

	err = t.FlagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	t.Args = t.FlagSet.Args()
	t.Flags = flags.Elem().Interface()

	return t, nil
}

// Usage is the help flag handler for the Trigger.
func (t *Trigger) Usage() {
	rep := t.Reply()

	desc := t.Command.Description
	if desc == "" {
		desc = "*No description.*"
	}

	rep.Embed = &discord.Embed{
		Title:       "`" + t.Command.Name + "` Usage",
		Description: desc,
	}

	t.FlagSet.VisitAll(func(f *flag.Flag) {
		rep.Embed.Fields = append(
			rep.Embed.Fields, discord.EmbedField{
				Name:   "Flag `-" + f.Name + "`",
				Value:  f.Usage + "\nDefault: `" + f.DefValue + "`",
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
		err = fmt.Errorf("send msg: %w", err)
		logs.Warn.Println(err)

		return err
	}

	return nil
}
