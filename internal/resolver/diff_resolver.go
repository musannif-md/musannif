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

// applyAndGenerateDiff applies the given patches to the current text,
// generates the diff between the original text and the patched text,
// and returns the generated diff and any error.
func (s *DiffResolver) applyAndGenerateDiff(patches []diffmatchpatch.Patch) ([]diffmatchpatch.Diff, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	originalText := s.text // Store the original text for diffing

	newText, results := s.dmp.PatchApply(patches, s.text)
	allApplied := true
	for _, applied := range results {
		if !applied {
			allApplied = false
			break
		}
	}

	if !allApplied {
		return nil, fmt.Errorf("not all patches applied successfully")
	}

	s.text = newText
	s.version++

	diffs := s.dmp.DiffMain(originalText, s.text, false)
	diffs = s.dmp.DiffCleanupSemantic(diffs)

	return diffs, nil
}

func (s *DiffResolver) getText() (string, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.text, s.version
}
