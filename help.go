package route

import (
	"flag"
	"fmt"
	"log"
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
			t.runHelp(name)
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
			cmd, _ := t.Route.GetCmd(cmdName)
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
	r.cmdMu.RLock()
	defer r.cmdMu.RUnlock()

	cats = make(map[string][]string)

	for name, c := range r.cmdMap {
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

func (t *Trigger) runHelp(name string) {
	cmd, ok := t.Route.GetCmd(name)
	if !ok {
		rep := t.Reply()
		rep.Content = fmt.Sprintf("command `%s` not known", name)
		_ = rep.Send()

		return
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
		log.Printf("help trigger: fill flagset: %s", err)

		return
	}

	ht.Usage()
}
