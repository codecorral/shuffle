# Shuffle

## Project Overview
Shuffle provides declarative configuration for AgentDeck. Define agent-deck profiles, groups, session shells, and MCP attachments in static configuration files, then apply them as a "deal."

## Key Concepts
- **Deal** — A complete agent-deck configuration described declaratively
- **Groups** — Logical session groupings (e.g., Elaboration, Construction, Operations)
- **Session Shells** — Predefined session configurations with tools, prompts, and MCP servers
- **Profiles** — Reusable agent-deck profile definitions

## Key Dependencies
- [agent-deck](https://github.com/asheshgoplani/agent-deck) — The session manager Shuffle configures

## Project Structure
```
shuffle/
  openspec/          # OpenSpec change management
    changes/         # Active change proposals
    specs/           # Main specifications
```

## Development Guidelines
- Use OpenSpec for all design changes
- Configuration format should be declarative and human-readable
- Support agent-deck's full configuration surface (groups, sessions, shells, MCPs)
