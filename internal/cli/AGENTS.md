<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# internal/cli

## Purpose
Cobra command definitions for all diagramr subcommands. Each file corresponds to one command. The root command wires them together; `main.go` calls `NewRootCmd().Execute()`.

## Key Files

| File | Description |
|------|-------------|
| `root.go` | `NewRootCmd()` — creates the root `diagramr` command and attaches all subcommands |
| `generate.go` | `generate [path]` — stub for the diagram generation command |
| `init.go` | `init` — writes a default `diagramr.config.yaml` to the current directory |
| `validate.go` | `validate` — loads config via `config.Load()` and prints a summary if valid |
| `version.go` | `version` — prints the binary version |

## For AI Agents

### Working In This Directory
- Add each new subcommand in its own file, named `<command>.go`.
- Wire it into `root.go` by calling `cmd.AddCommand(newXxxCmd())`.
- Commands write output via `cmd.OutOrStdout()` (not `fmt.Println`) so tests can capture output.
- The `init` command guards against overwriting an existing config file.

### Testing Requirements
- Use `cobra.Command.SetOut` / `SetArgs` in tests to capture output without spawning a subprocess.
- `go test ./internal/cli/... -race`

### Common Patterns
- All commands return `error`; errors are printed by the `Execute()` caller in `main.go`.
- Config is loaded lazily inside `RunE`, not at command construction time.

## Dependencies

### Internal
- `internal/config` — `config.Load()` and `config.Default()` used by `init` and `validate`

### External
- `github.com/spf13/cobra` — command/flag framework

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
