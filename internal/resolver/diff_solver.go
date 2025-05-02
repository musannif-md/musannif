package resolver

import (
	"os"
	"sync"
)

type DiffSolver struct {
	mu    sync.Mutex
	f     *os.File
	fpath string
}

type diff struct {
}

func (s *DiffSolver) initialize() error {
	// TODO: open file and shit here; depends on diffing lib
	return nil
}

func (s *DiffSolver) cleanup() error {
	// TODO: close file and shit here; depends on diffing lib
	return nil
}

func (s *DiffSolver) resolve() (diff, error) {
	resolved := diff{}
	// TODO: implement diff algorithm, lock for a couple of seconds to let diffs accumulate
	return resolved, nil
}
