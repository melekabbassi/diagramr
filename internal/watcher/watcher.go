package watcher

import (
	"context"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch calls onChange whenever files under path change.
// Only events whose file extension is in exts are considered (empty = all).
// Debounces rapid-fire events with a 300ms window.
// Blocks until ctx is cancelled or a watcher error occurs.
func Watch(ctx context.Context, path string, exts []string, onChange func()) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	if err := addDirs(w, path); err != nil {
		return err
	}

	debounce := time.NewTimer(0)
	if !debounce.Stop() {
		<-debounce.C
	}
	pending := false

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-w.Events:
			if !ok {
				return nil
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) ||
				event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				if !extAllowed(event.Name, exts) {
					continue
				}
				if !pending {
					pending = true
					debounce.Reset(300 * time.Millisecond)
				}
			}

		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			return err

		case <-debounce.C:
			pending = false
			onChange()
		}
	}
}

func extAllowed(name string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}
	got := filepath.Ext(name)
	for _, e := range exts {
		if got == e {
			return true
		}
	}
	return false
}

func addDirs(w *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return w.Add(p)
		}
		return nil
	})
}
