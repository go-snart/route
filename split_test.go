package route_test

import (
	"reflect"
	"testing"

	"github.com/go-snart/route"
)

func TestSplit(t *testing.T) {
	const str = "foo `bar` ``baz`` ```qux``` zo`op"

	expect := []string{
		"foo",
		"bar",
		"baz",
		"qux",
		"zo`op",
	}

	split := route.Split(str)
	if !reflect.DeepEqual(split, expect) {
		t.Errorf("expect %v\ngot %v", expect, split)
	}
}
