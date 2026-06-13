<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# diagramr

## Purpose
A CLI tool that parses source code (currently Go, TypeScript planned) and generates architectural diagrams in multiple output formats (Mermaid, JSON, SVG, PNG). It walks a directory tree, extracts types, interfaces, and relationships via AST analysis, builds an intermediate graph model, then renders it to the chosen format.

## Key Files

| File | Description |
|------|-------------|
| `go.mod` | Module definition — `github.com/melekabbassi/diagramr`, Go 1.26.2 |
| `go.sum` | Dependency checksums |
| `.goreleaser.yml` | GoReleaser config for multi-platform binary releases |
| `.gitignore` | VCS ignore rules |
| `README.md` | Project overview and usage |
| `DIAGRAMR_GUIDE.md` | Developer guide |
| `DIAGRAMR_PLAN.md` | Implementation roadmap |
| `diagramr.md` | Architecture notes |

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `cmd/` | Binary entrypoint (see `cmd/AGENTS.md`) |
| `internal/` | All application logic — parsers, graph model, renderers, CLI (see `internal/AGENTS.md`) |
| `testdata/` | Fixture source files used by parser tests (see `testdata/AGENTS.md`) |
| `.github/` | GitHub Actions CI workflows (see `.github/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- Module path is `github.com/melekabbassi/diagramr` — use this prefix for all internal import paths.
- Run `go build ./...` and `go test ./... -race` before claiming any change complete.
- CI runs on ubuntu, macos, and windows — avoid OS-specific path assumptions.

### Testing Requirements
```
go test ./... -race -coverprofile=coverage.out
go vet ./...
```

### Common Patterns
- Core pipeline: source files → parser → `graph.DiagramGraph` → renderer → output format.
- All application logic lives under `internal/`; no exported library surface.
- Config is loaded from `diagramr.config.yaml` (current dir) or `DIAGRAMR_*` env vars via Viper.

## Dependencies

### External
- `github.com/spf13/cobra` — CLI command tree
- `github.com/spf13/viper` — config loading (YAML + env vars)
- `github.com/fsnotify/fsnotify` — filesystem watcher (used by `watcher` package)
- `github.com/smacker/go-tree-sitter` — tree-sitter bindings (future TypeScript parser)
- `github.com/stretchr/testify` — test assertions

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
