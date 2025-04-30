package connection

import (
	"os"
	_ "sync"
)

type DiffSolver struct {
	// mu sync.Mutex
	f *os.File
}
