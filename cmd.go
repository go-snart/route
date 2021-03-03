package route

import (
	"log"
	"sort"
	"sync"
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

// CmdStore is a concurrent-safe store of Cmds.
type CmdStore struct {
	mu sync.RWMutex
	ma map[string]Cmd
}

// AddCmds adds Cmds to the CmdStore.
// Duplicate names are skipped.
func (c *CmdStore) AddCmds(cs ...Cmd) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, cmd := range cs {
		if _, ok := c.ma[cmd.Name]; ok {
			log.Printf("cmd %q already exists (skipping)", cmd.Name)

			continue
		}

		c.ma[cmd.Name] = cmd
	}
}

// GetCmd fetches a Cmd from the CmdStore.
func (c *CmdStore) GetCmd(name string) (Cmd, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cmd, ok := c.ma[name]
	if !ok {
		return Cmd{}, false
	}

	return cmd, true
}

// DelCmd removes a Cmd from the CmdStore.
func (c *CmdStore) DelCmd(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.ma, name)
}

// CmdsByCat creates a map of Cmd categories and a list of category names.
// Each category, and category names, are sorted.
// If hidden is true, Cmds with Hide will be included.
func (c *CmdStore) CmdsByCat(hidden bool) (map[string][]Cmd, []string) {
	cats := make(map[string][]Cmd)

	c.mu.RLock()

	for _, cmd := range c.ma {
		if !cmd.Hide || hidden {
			cats[cmd.Cat] = append(cats[cmd.Cat], cmd)
		}
	}

	c.mu.RUnlock()

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
