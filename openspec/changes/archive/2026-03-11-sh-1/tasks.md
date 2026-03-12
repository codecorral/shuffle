## 1. Project Setup

- [x] 1.1 Choose implementation language and initialize project (Go recommended — matches agent-deck, single binary, good YAML/TOML libs)
- [x] 1.2 Add YAML parsing dependency and define Go structs for deck schema (name, profile, mcps, shells, groups, sessions, conductors)
- [x] 1.3 Add TOML parsing dependency for reading/writing agent-deck config.toml

## 2. Deck Parsing & Validation

- [x] 2.1 Implement deck file parser — load `.deck.yaml`, unmarshal into structs
- [x] 2.2 Implement validation: required fields, shell reference resolution, MCP reference validation, mutual exclusivity of agent/command/shell on sessions
- [x] 2.3 Implement shell template expansion — merge shell fields into session, with session fields overriding (except MCPs which merge as union)
- [x] 2.4 Implement markdown field resolution — detect file-ref vs inline for claude_md, policy_md, prompt fields
- [x] 2.5 Implement `shuffle validate <deck.yaml>` CLI command

## 3. State Discovery

- [x] 3.1 Implement agent-deck state reader — query `agent-deck list --json`, `group list --json`, `conductor list --json`, `mcp list --json` to build current state
- [x] 3.2 Implement config.toml reader — parse existing MCP definitions from `~/.agent-deck/config.toml`
- [x] 3.3 Implement diff engine — compare deck desired state against current state, produce a list of actions (create profile, create group, create session, create conductor, attach MCP, move conductor to group)

## 4. Diff Command

- [x] 4.1 Implement `shuffle diff <deck.yaml>` CLI command — runs state discovery + diff engine, prints planned actions without executing

## 5. Deal Execution

- [x] 5.1 Implement profile creation — `agent-deck profile create <name>` if not exists
- [x] 5.2 Implement MCP creation — write `[mcps.<name>]` entries to config.toml if not exists
- [x] 5.3 Implement group creation — `agent-deck group create <name>` for each missing group
- [x] 5.4 Implement conductor creation — `agent-deck conductor setup <name>` with description, heartbeat, claude-md, policy-md flags; then `agent-deck group move <conductor-session> <group>`
- [x] 5.5 Implement session creation — `agent-deck launch <path> -c <agent|command> -t <name> -g <group>` with optional `--parent`, `--mcp`, `-w <worktree>`, `-m <prompt>` flags
- [x] 5.6 Implement MCP attachment reconciliation — check `agent-deck mcp attached <session> --json`, then `agent-deck mcp attach <session> <mcp>` for each missing binding
- [x] 5.7 Implement `shuffle deal <deck.yaml>` CLI command — runs validation, state discovery, diff, then executes actions in order

## 6. Example & Documentation

- [x] 6.1 Create a sample `codecorral.deck.yaml` demonstrating all features (profile, MCPs, shells, conductors, sessions with inline/file-ref prompts, ralph-tui command session)
- [x] 6.2 Update README.md with deck format reference and CLI usage
