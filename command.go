package route

import (
	"log"
	"reflect"
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

// AddCmd adds a Cmd to the Route.
func (r *Route) AddCmd(c Cmd) {
	for _, r := range c.Name {
		if !unicode.IsLower(r) && !unicode.IsNumber(r) {
			log.Panicf("cmd name %q not alphanumeric", c.Name)
		}
	}

	val := reflect.ValueOf(c)
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if f.IsZero() && f.Type().Kind() != reflect.Bool {
			log.Panicf("cmd %q missing field %s", c.Name, val.Type().Field(i).Name)
		}
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
