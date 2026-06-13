<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# internal/parser

## Purpose
Defines the `Parser` interface, a `Registry` for looking up parsers by language name, and a file-extension matcher utility. Language-specific implementations live in subdirectories.

## Key Files

| File | Description |
|------|-------------|
| `parser.go` | `Parser` interface (`Parse`, `Language`, `Extensions`) and `Options` struct |
| `registry.go` | `Registry` — maps language name → `Parser`; `Get(lang)` returns error for unknown languages |
| `scanner.go` | `MatchExtension(path, exts)` — checks if a file path has one of the given extensions |

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `golang/` | Go AST-based parser implementation (see `golang/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- To add a new language parser: create `parser/<lang>/parser.go`, implement the `Parser` interface, then register it with `NewRegistry(golang.New(), newlang.New(), ...)`.
- `Options.IncludePrivate` — when false, unexported types/methods are omitted from the graph.
- `Options.Exclude` — glob patterns; scanner should skip matching paths.
- `Options.MaxDepth` — limits directory traversal depth (parser implementations are responsible for honoring this).

### Testing Requirements
- Parser tests use fixture files from `testdata/go/` — add corresponding fixtures when testing new scenarios.
- `go test ./internal/parser/... -race`

## Dependencies

### Internal
- `internal/graph` — parsers return `*graph.DiagramGraph`

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
