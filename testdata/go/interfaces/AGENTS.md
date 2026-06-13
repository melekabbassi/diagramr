<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# testdata/go/interfaces

## Purpose
Fixture for interface definitions and implementation inference. Tests that the parser emits `RelImplements` edges when a struct's method set satisfies an interface — purely structural, no explicit declaration.

## Key Files

| File | Description |
|------|-------------|
| `interfaces.go` | Package `interfaces` — interface type(s) and structs that implement them, used by `emitImplementsEdges` |

## For AI Agents

### Working In This Directory
- Implementation is inferred by method name matching only (not signatures) — keep test interfaces simple.
- An empty interface (`interface{}`) is never considered implemented per the parser logic (`len(i.Methods) == 0` guard).

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
