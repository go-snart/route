package route

import (
	"bytes"
	"log"
	"unicode"
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

func toAlNum(s string) string {
	b := bytes.NewBuffer(make([]byte, 0, len(s)))

	for _, r := range s {
		if unicode.IsLower(r) || unicode.IsNumber(r) {
			b.WriteRune(r)
		}
	}

	return b.String()
}

// AddCmd adds a Cmd to the Route.
func (r *Route) AddCmd(c Cmd) {
	name := toAlNum(c.Name)
	if name != c.Name {
		log.Printf("truncating cmd name %q to %q", c.Name, name)
	}

	c.Name = name

	if c.Name == "" {
		log.Panicf("cmd name is missing")
	}

	if c.Desc == "" {
		log.Panicf("cmd %q desc is missing", c.Name)
	}

	if c.Cat == "" {
		log.Panicf("cmd %q cat is missing", c.Name)
	}

	if c.Func == nil {
		log.Panicf("cmd %q func is missing", c.Name)
	}

	if c.Flags == nil {
		log.Panicf("cmd %q flags is missing", c.Name)
	}

	if _, ok := r.GetCmd(c.Name); ok {
		log.Panicf("cmd %q already exists", c.Name)
	}

	r.cmdMu.Lock()
	defer r.cmdMu.Unlock()

	r.cmdMap[c.Name] = c
}

// GetCmd fetches a Cmd by name from the Route.
// The bool indicates whether it was found.
func (r *Route) GetCmd(name string) (Cmd, bool) {
	r.cmdMu.RLock()
	defer r.cmdMu.RUnlock()

	c, ok := r.cmdMap[name]

	return c, ok
}

// DelCmd removes a Cmd from the Route.
func (r *Route) DelCmd(name string) {
	r.cmdMu.Lock()
	defer r.cmdMu.Unlock()

	delete(r.cmdMap, name)
}
