package printdb

import (
	"github.com/boltdb/bolt"
	"io/ioutil"
	"os"
)

type testContext struct {
	db *bolt.DB
}

func setup() (*testContext, error) {
	path := tempfile()
	db, err := bolt.Open(path, 0600, nil)

	if err != nil {
		return nil, err
	}

	return &testContext{
		db: db,
	}, nil
}

func teardown(ctx *testContext) {
	defer os.Remove(ctx.db.Path())
	ctx.db.Close()
}

// tempfile returns a temporary file path.
func tempfile() string {
	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}
