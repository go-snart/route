package route

import (
	"fmt"
	"sort"
	"strings"

	"github.com/diamondburned/arikawa/discord"
)

// HelpFlags is flags for Help.
type HelpFlags struct {
	Help bool `default:"false" usage:"helpception"`
}

// HelpCommand makes the Route's help menu Command.
func (r *Route) HelpCommand() *Command {
	return &Command{
		Name:  "help",
		Desc:  "a help menu",
		Func:  r.Help,
		Hide:  false,
		Flags: HelpFlags{},
	}
}

// Help is a Func that provides a help menu.
func (r *Route) Help(t *Trigger) error {
	flags := t.Flags.(HelpFlags)

	if flags.Help {
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
		Title:       fmt.Sprintf("%s Help", t.State.Ready.User.Username),
		Description: fmt.Sprintf("prefix: `%s`", t.Prefix.Clean),
	}

	cats := make([]string, 0)

	for cat, cmds := range r.Cats {
		if len(cmds) == 0 {
			continue
		}

		cats = append(cats, cat)
	}

	sort.Strings(cats)

	for _, name := range cats {
		helps := []string(nil)

		for _, c := range r.Cats[name] {
			helps = append(helps, fmt.Sprintf(
				"`%s%s`: *%s*",
				t.Prefix.Clean, c.Name,
				strings.SplitN(c.Desc, "\n", 2)[0],
			))
		}

		rep.Embed.Fields = append(rep.Embed.Fields, discord.EmbedField{
			Name:   name,
			Value:  strings.Join(helps, "\n"),
			Inline: false,
		})
	}

	rep.Embed.Footer = &discord.EmbedFooter{
		Text: "use the `-help` flag on a command for detailed help",
	}

	return rep.Send()
}

func (r *Route) runHelp(t *Trigger, name string) {
	cmd := (*Command)(nil)

	for _, cmds := range t.Cats {
		for _, c := range cmds {
			if c.Name == name {
				cmd = c

				break
			}
		}
	}

	if cmd == nil {
		rep := t.Reply()
		rep.Content = fmt.Sprintf("command `%s` not known", name)
		_ = rep.Send()

		return
	}

	(&Trigger{
		Route:   t.Route,
		Command: cmd,

		Message: discord.Message{
			ChannelID: t.Message.ChannelID,
		},
		Prefix: t.Prefix,
		Output: &strings.Builder{},
	}).Usage()
}
