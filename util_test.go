package route_test

import (
	"fmt"

	"github.com/go-snart/db"
)

func testDB() *db.DB {
	const uri = "mem://"

	d, err := db.Open(uri)
	if err != nil {
		err = fmt.Errorf("open %q: %w", uri, err)
		panic(err)
	}

	return d
}
