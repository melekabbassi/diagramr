<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# cmd

## Purpose
Container for the compiled binary entrypoint. Follows standard Go project layout — each subdirectory becomes its own binary.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `diagramr/` | Main binary — wires up the CLI and calls `os.Exit` on error (see `diagramr/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- Do not add application logic here; keep `main.go` files thin.
- New binaries get their own subdirectory (e.g., `cmd/diagramr-server/`).

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
