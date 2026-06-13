<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# cmd/diagramr

## Purpose
The `main` package for the `diagramr` binary. Delegates entirely to `internal/cli` — its only job is to call `cli.NewRootCmd().Execute()` and exit with a non-zero status on error.

## Key Files

| File | Description |
|------|-------------|
| `main.go` | Binary entrypoint — calls `cli.NewRootCmd().Execute()`, writes errors to stderr |

## For AI Agents

### Working In This Directory
- Keep `main.go` minimal. All logic belongs in `internal/`.
- Do not import any package other than `internal/cli` and standard library here.

### Testing Requirements
- No unit tests here; test via the `internal/cli` package or integration tests.

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
