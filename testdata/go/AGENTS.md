<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# testdata/go

## Purpose
Go fixture packages used by `internal/parser/golang` tests. Each subdirectory is a minimal, self-contained Go package designed to exercise a specific AST parsing scenario.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `simple/` | A single struct with fields — baseline parsing scenario (see `simple/AGENTS.md`) |
| `embedded/` | Struct embedding — tests `RelEmbeds` edge extraction (see `embedded/AGENTS.md`) |
| `interfaces/` | Interface definitions — tests `RelImplements` inference (see `interfaces/AGENTS.md`) |
| `multifile/` | A type split across two files — tests multi-file package handling (see `multifile/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- Keep fixture packages small and focused — one scenario per directory.
- Do not add build constraints or external imports; fixtures must compile with `go build` using only the standard library.
- Fixture files are **not** excluded from `go vet` or `go build ./...` — keep them valid Go.

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
