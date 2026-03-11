## Context

Shuffle is a new tool that provides declarative configuration for agent-deck. Agent-deck currently requires imperative CLI commands to set up profiles, groups, sessions, conductors, and MCPs. This works for ad-hoc use but breaks down when you need reproducible, version-controlled layouts — especially for projects like CodeCorral that define complex multi-group session topologies.

Agent-deck's CLI surface (as of v0.24.1):
- `profile create/delete/default`
- `group create/delete/move`
- `add`/`launch` for sessions (with `-t`, `-g`, `-c`, `--mcp`, `--parent`, `-m` flags)
- `conductor setup` (with `-description`, `-heartbeat`, `-claude-md`, `-policy-md`)
- `mcp attach/detach` for per-session MCP bindings
- Config stored in `~/.agent-deck/config.toml` (MCPs, profiles) and `~/.agent-deck/sessions.json` (sessions, groups)

## Goals / Non-Goals

**Goals:**
- Define a `.deck.yaml` format that covers profiles, MCPs, shells, groups, sessions, and conductors
- Build a CLI (`shuffle deal`, `shuffle validate`, `shuffle diff`) that applies deck files idempotently
- Support both file-referenced and inline markdown for conductor/session prompts
- Be additive and non-destructive — never remove ad-hoc sessions or groups the user created manually

**Non-Goals:**
- Template variables or dynamic substitution in deck files
- Session lifecycle management (start/stop/restart) — shuffle only ensures existence
- Removing or syncing deletions — if something is removed from the deck, shuffle does not delete it from agent-deck
- Managing agent-deck's internal session state (conversation history, worktree branches)

## Decisions

### D1: YAML as deck format
**Choice**: YAML with `.deck.yaml` extension.
**Rationale**: Human-readable, supports multiline strings natively (critical for inline prompts/CLAUDE.md), widely understood. TOML's multiline handling is awkward for markdown content. JSON lacks comments and multiline.

### D2: Shell templates as internal abstraction
**Choice**: `shells` section defines reusable session templates resolved at deal-time. They are not an agent-deck entity — shuffle expands them inline before issuing CLI commands.
**Rationale**: Avoids repetition across sessions that share agent type, worktree config, and MCP sets. Keeps the deck DRY without requiring agent-deck to understand templates. Session-level fields override shell defaults.

### D3: Conductors nested under groups
**Choice**: Conductors are declared under `groups.<name>.conductors`, not as a top-level section.
**Rationale**: Users organize conductors within groups in practice. Agent-deck's `conductor setup` creates them at profile level, so shuffle runs `conductor setup` then `group move` to place them. The deck reflects the user's mental model, shuffle handles the mechanics.

### D4: File-ref vs inline markdown resolution
**Choice**: `claude_md`, `policy_md`, and `prompt` fields accept either a file path or inline YAML string. Resolution rule: if the value looks like a relative path (contains `/` or ends in `.md`), read the file; otherwise treat as inline content.
**Alternative considered**: Separate `claude_md_file` and `claude_md_inline` keys — rejected as unnecessarily verbose. The heuristic is simple and unambiguous in practice.

### D5: Idempotency via name matching
**Choice**: Shuffle matches existing entities by name (session title within group, group name, conductor name, MCP name). No metadata tagging or state file.
**Rationale**: Simple and transparent. The deck file is the source of truth. `agent-deck list --json` and `conductor list --json` provide the current state for diffing. No hidden state to get out of sync.
**Trade-off**: If a user manually renames a shuffle-created session, shuffle will create a new one. Acceptable — manual renames are an explicit user action.

### D6: Deal execution order
**Choice**: Fixed order — profile → MCPs → groups → conductors → sessions → MCP attachments.
**Rationale**: Each step depends on prior steps (sessions need groups, MCP attach needs sessions and MCPs). Deterministic ordering simplifies debugging and makes `shuffle diff` output predictable.

### D7: Command field for non-agent sessions
**Choice**: Sessions support a `command` field as an alternative to `agent`. When `command` is set, it's passed directly to `agent-deck launch -c <command>`.
**Rationale**: Supports ralph-tui, custom scripts, or any CLI tool as a session. Agent-deck's `-c` flag already accepts arbitrary commands (e.g., `codex --dangerously-bypass-approvals-and-sandbox`).

## Risks / Trade-offs

- **[Agent-deck CLI changes]** → Shuffle couples tightly to agent-deck's CLI flags and JSON output. Pin to known-good agent-deck versions and test against CLI output.
- **[No deletion]** → Stale sessions accumulate if removed from deck. Users must manually clean up. This is intentional — safety over convenience.
- **[Conductor group move race]** → Between `conductor setup` and `group move`, the conductor briefly exists in the default "conductor" group. Unlikely to cause issues in practice since deal is run as a single sequential process.
- **[Name collision]** → If two decks declare the same group or session name, they'll collide. Document that decks should use unique names scoped to their project.
