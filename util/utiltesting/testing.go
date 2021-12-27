package utiltesting

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
)

// RepositoryRoot is the absolute path to the root of the repository.
var RepositoryRoot string

// IsTest returns whether it's running within the context of a test.
func IsTest() bool {
	// `go test` recompiles test files into a binary with the '.test' extension, before running it.
	// We can thus check if we're running within a test by checking the binary extension.
	// See https://github.com/golang/go/blob/master/src/cmd/go/internal/test/test.go
	// Related: https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test
	return strings.HasSuffix(os.Args[0], ".test")
}

func init() {
	_, b, _, _ := runtime.Caller(0)
	RepositoryRoot = filepath.Join(filepath.Dir(b), "../../..") // root is 3 levels above this package

	// Test environment should only be loaded when in the context of a test.
	// The testing package may be used by non test code so this check is necessary to guarantee the above.
	if !IsTest() {
		return
	}

	// Ignoring error because .env.testing is only present locally. On CI env variable are used instead.
	_ = godotenv.Load(filepath.Join(RepositoryRoot, ".env.testing"))
}

// AbsolutePath returns the absolute path to the given folder,
// regardless of where this function is called from.
func AbsolutePath(dir string) string {
	return filepath.Join(RepositoryRoot, dir)
}
