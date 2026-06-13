<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# testdata/go/multifile

## Purpose
Fixture for types and methods split across multiple files in the same package. Tests that the parser correctly merges methods from all files into the same node when they share a receiver type.

## Key Files

| File | Description |
|------|-------------|
| `user.go` | Package `multifile` — struct definition for `User` |
| `user_methods.go` | Package `multifile` — method declarations on `User` (separate file, same package) |

## For AI Agents

### Working In This Directory
- The parser processes all files in a package together via `parser.ParseDir`, so cross-file method attachment is handled automatically.
- This fixture validates that `collectTypes` and `collectMethodsAndImports` are called over all files before edge emission.

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
