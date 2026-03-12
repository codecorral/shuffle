# Shuffle

Declarative configuration for [AgentDeck](https://github.com/asheshgoplani/agent-deck). Define your entire agent-deck setup — profiles, groups, sessions, conductors, MCPs — in a single YAML file and apply it idempotently.

## Overview

Shuffle reads a `.deck.yaml` file — a declarative description of your agent-deck layout — and applies it idempotently. Instead of manually running dozens of CLI commands to set up profiles, groups, sessions, and conductors, define the desired state in a deck file and run `shuffle deal` to make it so.

## Concepts

- **Deck** — A `.deck.yaml` file describing the desired agent-deck layout
- **Shells** — Reusable session templates with agents, MCPs, worktree config, and prompts
- **Groups** — Logical groupings of sessions and conductors (e.g., elaboration, construction, operations)
- **Conductors** — Meta-agent orchestrators declared within groups
- **Profiles** — Agent-deck profile targeting (which profile the deck applies to)

## Installation

```bash
go install github.com/afterthought/shuffle/cmd/shuffle@latest
```

Or build from source:

```bash
git clone https://github.com/afterthought/shuffle.git
cd shuffle
go build -o shuffle ./cmd/shuffle/
```

Or with Nix:

```bash
nix build github:codecorral/shuffle
# or add as a flake input:
# shuffle.url = "github:codecorral/shuffle";
```

Requires [agent-deck](https://github.com/asheshgoplani/agent-deck) on PATH.

## CLI Usage

```bash
# Validate a deck file
shuffle validate my-project.deck.yaml

# Preview what changes would be made
shuffle diff my-project.deck.yaml

# Apply the deck to agent-deck
shuffle deal my-project.deck.yaml

# Override the target profile (takes priority over deck YAML profile.name)
shuffle deal --profile work my-project.deck.yaml
```

## Deck Format Reference

A `.deck.yaml` file has these top-level sections:

```yaml
name: my-project          # Required: deck name

profile:                   # Optional: agent-deck profile
  name: my-profile

mcps:                      # Optional: MCP server definitions
  github:
    command: npx
    args: ["-y", "@anthropic/mcp-github"]
    env:
      GITHUB_TOKEN: "${GITHUB_TOKEN}"
    description: GitHub MCP

shells:                    # Optional: reusable session templates
  research:
    agent: claude
    worktree: subdirectory
    mcps: [github]
    prompt: "Focus on research and design."

groups:                    # Optional: session groups
  elaboration:
    conductors:
      clint:
        description: "Elaboration conductor"
        heartbeat: true
        claude_md: |
          # Conductor instructions
          Manage elaboration sessions.
    sessions:
      intent:
        shell: research
        mcps: [exa]        # Merged with shell MCPs
        prompt: "Begin elaboration."
  operations:
    sessions:
      monitor:
        command: ralph-tui run --parallel 3
        path: .
```

### Session Fields

| Field | Description |
|-------|-------------|
| `agent` | Agent to run (e.g., `claude`, `gemini`) |
| `command` | Arbitrary command (e.g., `ralph-tui`) |
| `shell` | Reference to a shell template |
| `worktree` | Worktree strategy: `subdirectory`, `sibling`, or path |
| `mcps` | MCP servers to attach |
| `skills` | Skills to attach |
| `prompt` | Initial message (inline string or file path) |
| `path` | Project directory (default: `.`) |
| `parent` | Parent session for sub-session linking |

Exactly one of `agent`, `command`, or `shell` must be set.

### Shell Merging

When a session references a shell, fields are merged:
- **MCPs**: Union (shell + session)
- **Skills**: Union (shell + session)
- **All other fields**: Session overrides shell

### Markdown Fields

`claude_md`, `policy_md`, and `prompt` accept either inline content or a file path:
- Contains `/` or ends in `.md` → treated as file path
- Otherwise → treated as inline content
- YAML `|` block scalars are inline content

## Behavior

- **Idempotent**: Running `shuffle deal` twice makes no duplicate changes
- **Additive**: Only creates missing entities, never deletes
- **Non-destructive**: Manual sessions/groups are never touched
- **Ordered**: Profile → MCPs → Groups → Conductors → Sessions → MCP/Skill attachments

## Related Repositories

- [codecorral](https://github.com/codecorral/codecorral) — Agent orchestration framework (uses Shuffle for agent-deck setup)
- [agent-deck](https://github.com/asheshgoplani/agent-deck) — Terminal session manager for AI coding agents

## License

TBD
