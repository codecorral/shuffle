# Shuffle

## Project Overview
Shuffle provides declarative configuration for AgentDeck. Define agent-deck profiles, groups, session shells, and MCP attachments in static configuration files, then apply them as a "deal."

## Key Concepts
- **Deck** — A `.deck.yaml` file describing the desired agent-deck layout
- **Shells** — Reusable session templates with agents, MCPs, worktree config, and prompts
- **Groups** — Logical groupings of sessions and conductors
- **Conductors** — Meta-agent orchestrators declared within groups
- **Profiles** — Agent-deck profile targeting

## Key Dependencies
- [agent-deck](https://github.com/asheshgoplani/agent-deck) — The session manager Shuffle configures

## Project Structure
```
shuffle/
  cmd/shuffle/       # CLI entry point (validate, diff, deal commands)
  internal/
    deck/            # Deck file parsing, validation, shell resolution
    state/           # Agent-deck state discovery (profiles, groups, sessions, conductors, MCPs)
    diff/            # Diff engine — computes actions from desired vs current state
    apply/           # Action executor — runs agent-deck CLI commands
  examples/          # Example .deck.yaml files (simple, medium, complex)
    conductors/      # Conductor prompt templates
    prompts/         # Session prompt templates
  openspec/          # OpenSpec change management
    changes/         # Active change proposals
    specs/           # Main specifications
```

## Development Guidelines
- Use OpenSpec for all design changes
- Configuration format should be declarative and human-readable
- Support agent-deck's full configuration surface (groups, sessions, shells, MCPs, skills, conductors)
- Build: `go build ./cmd/shuffle/`
- Build (Nix): `nix build`
- Test: `go test ./...`
- Requires Go 1.18+ (uses generics)

## CLI Usage
```
shuffle validate [--profile <name>] <deck.yaml>   # Parse, validate, check references
shuffle diff [--profile <name>] <deck.yaml>       # Show planned actions (desired vs current)
shuffle deal [--profile <name>] [--warn-only] <deck.yaml>  # Apply the deck
shuffle version                                    # Show version
```

## Runtime Requirements
- `agent-deck` binary must be on PATH (used for state discovery and action execution)
- Config file: `~/.agent-deck/config.toml` (read for MCP discovery, appended for new MCPs)

## Gotchas
- **Shell merging**: Sessions inherit from shell templates — MCPs/skills are union-merged, other fields are overridden by session
- **Markdown fields**: `claude_md`, `policy_md`, `prompt` can be inline strings or file paths (contains `/` or ends in `.md`); paths resolve relative to the deck file
- **Best-effort apply**: `apply.Execute()` logs errors but continues — it won't halt on partial failures
- **Config append**: MCPs are appended as raw TOML to config.toml to preserve existing comments/formatting (not re-encoded)
- **No deletions**: Diff engine only creates missing entities, never deletes existing ones
