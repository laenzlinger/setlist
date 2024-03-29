package testinginit

import (
	"os"
	"path"
	"runtime"
)

//nolint:gochecknoinits // required for testing
func init() {
	// assumption is that the program is started in the reqpertoire directory
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../test/Repertoire")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
