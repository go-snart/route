package route

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
)

// HelpFlags is flags for Help.
type HelpFlags struct {
	Help bool `default:"false" usage:"helpception"`
}

//nolint:gochecknoglobals // useful global
// HelpCmd is a Cmd providing a help menu.
var HelpCmd = Cmd{
	Name: "help",
	Desc: "a help menu",
	Cat:  "help",
	Func: HelpFunc,
	Hide: false,
	Flags: HelpFlags{
		Help: false,
	},
}

// HelpFunc is a Func that provides a help menu.
func HelpFunc(t *Trigger) error {
	if t.Flags.(HelpFlags).Help {
		rep := t.Reply()
		rep.Content = "helpception :thinking:"

		return rep.Send()
	}

	if len(t.Args) > 0 {
		for _, name := range t.Args {
			err := t.runHelp(name)
			if err != nil {
				return fmt.Errorf("run help %q: %w", name, err)
			}
		}

		return nil
	}

	rep := t.Reply()
	rep.Embed = &discord.Embed{
		Title:       fmt.Sprintf("%s Help", t.DisplayName()),
		Description: fmt.Sprintf("prefix: `%s`", t.Prefix.Clean),
	}

	cats, catNames := t.Route.cats()

	for _, catName := range catNames {
		cmds := cats[catName]
		helps := make([]string, 0, len(cmds))

		for _, cmdName := range cmds {
			cmd := t.Route.Cmds[cmdName]
			helps = append(helps, fmt.Sprintf(
				"`%s%s`: *%s*",
				t.Prefix.Clean, cmd.Name,
				strings.SplitN(cmd.Desc, "\n", 2)[0],
			))
		}

		rep.Embed.Fields = append(rep.Embed.Fields, discord.EmbedField{
			Name:  catName,
			Value: strings.Join(helps, "\n"),
		})
	}

	rep.Embed.Footer = &discord.EmbedFooter{
		Text: "use the `-help` flag on a command for detailed help",
	}

	return rep.Send()
}

func (r *Route) cats() (cats map[string][]string, catNames []string) {
	cats = make(map[string][]string)

	for name, c := range r.Cmds {
		if c.Hide {
			continue
		}

		cats[c.Cat] = append(cats[c.Cat], name)
	}

	catNames = make([]string, 0, len(cats))

	for name := range cats {
		sort.Strings(cats[name])
		catNames = append(catNames, name)
	}

	sort.Strings(catNames)

	return
}

func (t *Trigger) runHelp(name string) error {
	cmd, ok := t.Route.Cmds[name]
	if !ok {
		rep := t.Reply()
		rep.Content = fmt.Sprintf("command `%s` not known", name)

		return rep.Send()
	}

	ht := &Trigger{
		Route: t.Route,
		Message: discord.Message{
			ChannelID: t.Message.ChannelID,
		},
		Prefix:  t.Prefix,
		Command: cmd,
		FlagSet: flag.NewFlagSet(cmd.Name, flag.ContinueOnError),
		Args:    nil,
		Flags:   nil,
		Output:  &strings.Builder{},
	}
	if _, err := ht.fillFlagSet(); err != nil {
		return fmt.Errorf("fill flagset: %w", err)
	}

	ht.Usage()

	return nil
}
