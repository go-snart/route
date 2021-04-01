package route

import (
	"sort"
	"sync"
)

// CmdStore is a concurrent-safe store of Cmds.
type CmdStore struct {
	ma map[string]Cmd
	mu sync.RWMutex
}

// NewCmdStore creates a usable CmdStore.
func NewCmdStore() *CmdStore {
	return &CmdStore{
		ma: map[string]Cmd{},
		mu: sync.RWMutex{},
	}
}

// Add stores a Cmd, using its defined name.
func (c *CmdStore) Add(cmd Cmd) {
	c.mu.Lock()
	c.ma[cmd.Name] = cmd
	c.mu.Unlock()
}

// Get fetches a Cmd with the given name.
func (c *CmdStore) Get(name string) (Cmd, bool) {
	c.mu.RLock()
	cmd, ok := c.ma[name]
	c.mu.RUnlock()

	return cmd, ok
}

// Del removes a Cmd with the given name.
func (c *CmdStore) Del(name string) {
	c.mu.Lock()
	delete(c.ma, name)
	c.mu.Unlock()
}

// Cmd is a command.
type Cmd struct {
	Name  string
	Desc  string
	Cat   string
	Func  Func
	Hide  bool
	Flags interface{}
}

// ByCat creates a map of sorted Cmd categories, and a sorted list of category names.
//
// If hidden is true, Cmds with the Hide flag will be included.
func (c *CmdStore) ByCat(hidden bool) (map[string][]Cmd, []string) {
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
