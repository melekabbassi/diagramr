<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# internal

## Purpose
All application logic for diagramr. Structured as independent packages that together implement the parse → graph → render pipeline. Nothing here is exported as a public library.

## Subdirectories

| Directory | Purpose |
|-----------|---------|
| `cli/` | Cobra command definitions for all CLI subcommands (see `cli/AGENTS.md`) |
| `config/` | Config schema and Viper-based loader (see `config/AGENTS.md`) |
| `graph/` | Intermediate graph model — nodes, edges, metadata (see `graph/AGENTS.md`) |
| `output/` | Output writers — JSON, Mermaid text, SVG, PNG (see `output/AGENTS.md`) |
| `parser/` | Parser interface, registry, and language implementations (see `parser/AGENTS.md`) |
| `renderer/` | Renderer interface and Mermaid class diagram renderer (see `renderer/AGENTS.md`) |
| `watcher/` | Filesystem watch loop for live-reload mode (see `watcher/AGENTS.md`) |

## For AI Agents

### Working In This Directory
- Packages are intentionally decoupled; prefer passing `graph.DiagramGraph` across boundaries rather than sharing state.
- Add new language parsers under `parser/<lang>/`, implement the `parser.Parser` interface, and register in the `Registry`.
- Add new output formats under `output/` and wire into `output.Write`.

### Common Patterns
- Data flows strictly: `parser` → `graph` → `renderer` → `output`.
- Visibility (public/private) is determined by Go's uppercase-first convention and reflected in `graph.Visibility`.

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
