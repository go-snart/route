package route

import (
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/v2/discord"
)

// CatBuiltin is the category for builtin commands.
const CatBuiltin = "builtin"

// HelpFlags is flags for Help.
type HelpFlags struct {
	Help bool `default:"false" usage:"helpception"`
}

// HelpCommand makes the Route's help menu Command.
func (r *Route) HelpCommand() Cmd {
	return Cmd{
		Name: "help",
		Desc: "a help menu",
		Cat:  CatBuiltin,
		Func: r.Help,
		Hide: false,
		Flags: HelpFlags{
			Help: false,
		},
	}
}

// Help is a Func that provides a help menu.
func (r *Route) Help(t *Trigger) error {
	if t.Flags.(HelpFlags).Help {
		rep := t.Reply()
		rep.Content = "helpception :thinking:"

		return rep.Send()
	}

	if len(t.Args) > 0 {
		for _, name := range t.Args {
			r.runHelp(t, name)
		}

		return nil
	}

	rep := t.Reply()
	rep.Embed = &discord.Embed{
		Title:       fmt.Sprintf("%s Help", t.DisplayName()),
		Description: fmt.Sprintf("prefix: `%s`", t.Prefix.Clean),
	}

	cats, catNames := r.cats()

	for _, catName := range catNames {
		cmds := cats[catName]
		helps := make([]string, 0, len(cmds))

		for _, cmdName := range cmds {
			cmd, _ := r.GetCmd(cmdName)
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

func (r *Route) runHelp(t *Trigger, name string) {
	cmd, ok := t.Router.GetCmd(name)
	if !ok {
		rep := t.Reply()
		rep.Content = fmt.Sprintf("command `%s` not known", name)
		_ = rep.Send()

		return
	}

	(&Trigger{
		Router: t.Router,
		Message: discord.Message{
			ChannelID: t.Message.ChannelID,
		},
		Prefix:  t.Prefix,
		Command: cmd,
		FlagSet: nil,
		Args:    nil,
		Flags:   nil,
		Output:  &strings.Builder{},
	}).Usage()
}
