package route_test

import (
	"github.com/go-snart/route"
)

func testCmd() (route.Cmd, *string) {
	run := ""

	return route.Cmd{
		Name:  testName,
		Desc:  testDesc,
		Cat:   testCat,
		Func:  testFunc(&run),
		Hide:  false,
		Flags: testFlags{},
	}, &run
}
