package resolver

import (
	"fmt"
	"os"
	"sync"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type DiffResolver struct {
	mu      sync.Mutex
	fpath   string
	version int
	text    string
	dmp     *diffmatchpatch.DiffMatchPatch
}

func (s *DiffResolver) initialize() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.fpath)
	if err != nil {
		return fmt.Errorf("failed to read file during resolver init: %w", err)
	}

	s.text = string(data)
	s.dmp = diffmatchpatch.New()
	s.version = 0

	return nil
}

func (s *DiffResolver) cleanup() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Write contents to disk
	// TODO: eventually accomplish this via a ticker

	err := os.WriteFile(s.fpath, []byte(s.text), 0644)
	if err != nil {
		return fmt.Errorf("resolver failed to update local file: %w", err)
	}

	return nil
}

func (s *DiffResolver) resolve() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return nil
}
