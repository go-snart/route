package route

import (
	"log"
	"sort"
)

// Cmd is a command.
type Cmd struct {
	Name  string
	Desc  string
	Cat   string
	Func  Func
	Hide  bool
	Flags interface{}
}

// AddCmds adds Cmds to the Route.
// Duplicate names are skipped.
func (r *Route) AddCmds(cs ...Cmd) {
	r.cmdMu.Lock()
	defer r.cmdMu.Unlock()

	for _, cmd := range cs {
		if _, ok := r.cmdMa[cmd.Name]; ok {
			log.Printf("cmd %q already exists (skipping)", cmd.Name)

			continue
		}

		r.cmdMa[cmd.Name] = cmd
	}
}

// GetCmd fetches a Cmd from the Route.
func (r *Route) GetCmd(name string) (Cmd, bool) {
	r.cmdMu.RLock()
	defer r.cmdMu.RUnlock()

	cmd, ok := r.cmdMa[name]
	if !ok {
		return Cmd{}, false
	}

	return cmd, true
}

// DelCmd removes a Cmd from the Route.
func (r *Route) DelCmd(name string) {
	r.cmdMu.Lock()
	defer r.cmdMu.Unlock()

	delete(r.cmdMa, name)
}

// CmdsByCat creates a map of Cmd categories and a list of category names.
//
// Each category, and category names, are sorted.
//
// If hidden is true, Cmds with Hide will be included.
func (r *Route) CmdsByCat(hidden bool) (map[string][]Cmd, []string) {
	cats := make(map[string][]Cmd)

	r.cmdMu.RLock()

	for _, cmd := range r.cmdMa {
		if !cmd.Hide || hidden {
			cats[cmd.Cat] = append(cats[cmd.Cat], cmd)
		}
	}

	r.cmdMu.RUnlock()

	catNames := make([]string, 0, len(cats))

	for name := range cats {
		catNames = append(catNames, name)

		cat := cats[name]
		sort.Slice(cat, func(i, j int) bool {
			return cat[i].Name < cat[j].Name
		})
	}

	sort.Strings(catNames)

	return cats, catNames
}
