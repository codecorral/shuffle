## Why

Shuffle needs a declarative configuration format so users can describe their entire agent-deck setup — profiles, groups, sessions, conductors, MCPs — in a single file and apply it idempotently. Currently agent-deck requires manual CLI commands to build up a session layout, which is tedious, error-prone, and not reproducible.

## What Changes

- Define the `.deck.yaml` format: the schema for declaratively describing an agent-deck profile layout
- Implement the `shuffle deal` CLI command that applies a deck file to agent-deck (idempotent, additive, non-destructive)
- Implement `shuffle validate` for deck file validation
- Implement `shuffle diff` for dry-run change preview
- Support reusable session shell templates within the deck
- Support conductors nested under groups with inline or file-referenced CLAUDE.md and POLICY.md
- Support arbitrary commands (e.g., `ralph-tui`) as session types alongside agent names

## Capabilities

### New Capabilities
- `deck-format`: The YAML deck file schema — top-level sections (name, profile, mcps, shells, groups), session and conductor declarations, shell template resolution, file-ref vs inline markdown fields
- `deal-cli`: The `shuffle deal`, `shuffle validate`, and `shuffle diff` CLI commands — parsing deck files, diffing against current agent-deck state, and applying changes idempotently

### Modified Capabilities

## Impact

- New CLI binary/entrypoint for `shuffle` with subcommands `deal`, `validate`, `diff`
- Depends on `agent-deck` CLI being available on PATH (profile, group, session, conductor, mcp commands)
- Reads/writes `~/.agent-deck/config.toml` for MCP definitions
- Creates files in conductor directories (`~/.agent-deck/conductor/<name>/`) for inline CLAUDE.md/POLICY.md content
