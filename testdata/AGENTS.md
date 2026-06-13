<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# testdata

## Purpose
Fixture source files used as inputs by parser tests. Each subdirectory represents a language and contains small, self-contained Go (or future TypeScript) packages that exercise specific parsing scenarios.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `go/` | Go fixture packages (see `go/AGENTS.md`) |
| `typescript/` | TypeScript fixtures — empty, reserved for future TS parser |

## For AI Agents

### Working In This Directory
- Files here are **not** compiled into the binary — they are read at test time via filesystem paths.
- When adding a new parser test scenario, add a new subdirectory under the relevant language folder rather than modifying existing fixtures (tests may assert exact output).
- Do not add `_test.go` files here; tests live in the parser package alongside the fixtures they consume.

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
