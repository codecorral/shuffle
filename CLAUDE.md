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
  openspec/          # OpenSpec change management
    changes/         # Active change proposals
    specs/           # Main specifications
```

## Development Guidelines
- Use OpenSpec for all design changes
- Configuration format should be declarative and human-readable
- Support agent-deck's full configuration surface (groups, sessions, shells, MCPs, skills, conductors)
- Build and test: `go build ./cmd/shuffle/`
