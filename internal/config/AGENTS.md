<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# internal/config

## Purpose
Defines the `Config` struct and loads it from `diagramr.config.yaml` (current directory) with fallback to `DIAGRAMR_*` environment variables. Provides `Default()` for zero-config operation.

## Key Files

| File | Description |
|------|-------------|
| `schema.go` | `Config` struct with JSON/mapstructure tags; `Default()` returning `{language:"auto", format:"mermaid", max_nodes:200}` |
| `loader.go` | `Load() (Config, error)` — Viper-based loader; silently succeeds if no config file is present |

## For AI Agents

### Working In This Directory
- Config file name is `diagramr.config` (Viper appends `.yaml`); located in the working directory.
- Env var prefix is `DIAGRAMR_` — e.g., `DIAGRAMR_MAX_NODES=500`.
- `Load()` does **not** error when the config file is absent; it only errors on malformed files.
- Add new fields to both `schema.go` (struct tag) and `loader.go` (`v.SetDefault`).

### Testing Requirements
- Test `Load()` by writing a temp config file and changing the working directory, or by setting env vars.

## Dependencies

### External
- `github.com/spf13/viper` — config loading with YAML + env var support

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
