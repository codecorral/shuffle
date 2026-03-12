## ADDED Requirements

### Requirement: Deal command applies deck
The `shuffle deal <deck.yaml>` command SHALL parse the deck file and apply it to agent-deck, ensuring all declared entities exist.

#### Scenario: Apply a full deck
- **WHEN** user runs `shuffle deal codecorral.deck.yaml` with a deck declaring a profile, MCPs, groups, conductors, and sessions
- **THEN** shuffle SHALL create or verify each entity in order: profile → MCPs → groups → conductors → sessions → MCP attachments

#### Scenario: Apply same deck twice
- **WHEN** user runs `shuffle deal` twice with the same deck file and no external changes
- **THEN** the second run SHALL make no changes and report that everything is up to date

### Requirement: Idempotent execution
Shuffle deal SHALL be idempotent. It SHALL check for existing entities by name before creating them. It SHALL NOT create duplicates.

#### Scenario: Profile already exists
- **WHEN** the deck declares a profile that already exists in agent-deck
- **THEN** shuffle SHALL skip profile creation

#### Scenario: Session already exists in group
- **WHEN** the deck declares a session with title "intent" in group "elaboration" and a session with that title already exists in that group
- **THEN** shuffle SHALL skip session creation

#### Scenario: MCP already in config
- **WHEN** the deck declares an MCP "exa" and `[mcps.exa]` already exists in config.toml
- **THEN** shuffle SHALL skip MCP creation

### Requirement: Non-destructive behavior
Shuffle deal SHALL NOT remove, modify, or interfere with entities not declared in the deck file.

#### Scenario: Ad-hoc session preserved
- **WHEN** the user has manually created a session "experiment-7" in group "exploration" and the deck does not declare that session
- **THEN** shuffle deal SHALL leave "experiment-7" untouched

#### Scenario: Ad-hoc group preserved
- **WHEN** the user has manually created a group "sandbox" and the deck does not declare that group
- **THEN** shuffle deal SHALL leave "sandbox" untouched

#### Scenario: Entity removed from deck
- **WHEN** a session previously created by shuffle deal is removed from the deck file and deal is run again
- **THEN** shuffle SHALL NOT delete the session from agent-deck

### Requirement: Validate command
The `shuffle validate <deck.yaml>` command SHALL parse the deck file and report any schema errors without modifying agent-deck state.

#### Scenario: Valid deck
- **WHEN** user runs `shuffle validate` on a well-formed deck
- **THEN** shuffle SHALL exit 0 and report the deck is valid

#### Scenario: Invalid deck
- **WHEN** user runs `shuffle validate` on a deck with errors (missing required fields, unknown shell references, invalid types)
- **THEN** shuffle SHALL exit non-zero and list all validation errors

### Requirement: Diff command
The `shuffle diff <deck.yaml>` command SHALL compare the deck file against current agent-deck state and report what `shuffle deal` would create, without making changes.

#### Scenario: New entities needed
- **WHEN** user runs `shuffle diff` and the deck declares entities that don't exist yet
- **THEN** shuffle SHALL list each entity that would be created (e.g., "create group: elaboration", "create session: intent in elaboration")

#### Scenario: Everything exists
- **WHEN** user runs `shuffle diff` and all declared entities already exist
- **THEN** shuffle SHALL report no changes needed

### Requirement: Prompt sent on creation only
When a session has a `prompt` field, the prompt SHALL be sent only when the session is newly created. It SHALL NOT be re-sent if the session already exists.

#### Scenario: New session with prompt
- **WHEN** shuffle creates a new session that has a `prompt` field
- **THEN** shuffle SHALL use `agent-deck launch -m <prompt>` to create and send the prompt in one step

#### Scenario: Existing session with prompt
- **WHEN** shuffle encounters an existing session that matches a deck session with a `prompt` field
- **THEN** shuffle SHALL skip the session entirely (no re-send)

### Requirement: Execution order
Shuffle deal SHALL apply entities in a fixed order: profile → MCPs (config.toml) → groups → conductors (setup + group move) → sessions (launch + MCP attach).

#### Scenario: Session depends on group
- **WHEN** a deck declares sessions in a group
- **THEN** the group SHALL be created before any sessions are added to it

#### Scenario: MCP attach depends on MCP and session
- **WHEN** a session declares MCP attachments
- **THEN** both the MCP (in config.toml) and the session SHALL exist before `mcp attach` is called

### Requirement: MCP attachment reconciliation
Shuffle deal SHALL verify that each session's declared MCPs are attached. If an MCP is defined but not attached to a session, shuffle SHALL attach it.

#### Scenario: MCP already attached
- **WHEN** session "intent" should have MCP "github" and it is already attached
- **THEN** shuffle SHALL skip the attachment

#### Scenario: MCP not yet attached
- **WHEN** session "intent" should have MCP "github" and it is not attached
- **THEN** shuffle SHALL run `agent-deck mcp attach intent github`

### Requirement: Conductor existence check
Shuffle deal SHALL check if each declared conductor exists via `agent-deck conductor list --json`. If missing, it SHALL create it via `agent-deck conductor setup` and then move it to the declared group.

#### Scenario: Conductor already exists
- **WHEN** the deck declares conductor "clint" under group "clint" and a conductor named "clint" already exists
- **THEN** shuffle SHALL skip conductor setup

#### Scenario: Conductor with inline claude_md
- **WHEN** the deck declares a conductor with inline `claude_md` content
- **THEN** shuffle SHALL write the content to a temporary file, pass it to `conductor setup -claude-md`, then clean up

### Requirement: Error handling
Shuffle deal SHALL report errors clearly and continue processing remaining entities when possible.

#### Scenario: Agent-deck not on PATH
- **WHEN** the `agent-deck` CLI is not found
- **THEN** shuffle SHALL exit with a clear error message before attempting any operations

#### Scenario: One session fails
- **WHEN** session creation fails for one session (e.g., invalid path)
- **THEN** shuffle SHALL report the error and continue processing remaining sessions
