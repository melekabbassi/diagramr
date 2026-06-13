<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# internal/watcher

## Purpose
Filesystem watch loop for live-reload mode. Monitors a directory for source file changes and triggers diagram regeneration. Currently a stub.

## Key Files

| File | Description |
|------|-------------|
| `watcher.go` | `Watch(path string) error` — stub; will use `fsnotify` to watch for file changes |

## For AI Agents

### Working In This Directory
- Implement using `github.com/fsnotify/fsnotify` (already in `go.mod`).
- `Watch` should accept a callback or channel so the caller can trigger re-parsing on change events.
- Debounce rapid file-system events (e.g., editor save bursts) before triggering a regeneration.

### Testing Requirements
- Integration tests should create temp files, trigger writes, and assert the watcher fires within a timeout.
- `go test ./internal/watcher/... -race`

## Dependencies

### External
- `github.com/fsnotify/fsnotify` — cross-platform filesystem event notifications

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
