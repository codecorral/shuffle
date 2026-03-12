## ADDED Requirements

### Requirement: Deck file structure
A deck file SHALL be a YAML file with the `.deck.yaml` extension containing these top-level keys: `name` (required string), `profile` (optional object), `mcps` (optional map), `shells` (optional map), `groups` (optional map).

#### Scenario: Valid minimal deck
- **WHEN** a deck file contains only `name: my-deck`
- **THEN** the file SHALL be considered valid

#### Scenario: Unknown top-level key
- **WHEN** a deck file contains an unrecognized top-level key
- **THEN** validation SHALL fail with an error identifying the unknown key

### Requirement: Profile declaration
The `profile` section SHALL declare an agent-deck profile with `name` (required string) and optional agent-specific config blocks (e.g., `claude.config_dir`).

#### Scenario: Profile with claude config
- **WHEN** a deck declares `profile.name: my-profile` and `profile.claude.config_dir: ~/.claude-work`
- **THEN** shuffle SHALL ensure a profile named "my-profile" exists with the specified claude config directory

### Requirement: MCP declarations
The `mcps` section SHALL be a map of MCP name to MCP definition. Each definition SHALL include `command` (required string) and optionally `args` (string array), `env` (string map), and `description` (string).

#### Scenario: MCP with all fields
- **WHEN** a deck declares an MCP with command, args, env, and description
- **THEN** all fields SHALL be preserved when written to agent-deck's config.toml under `[mcps.<name>]`

#### Scenario: MCP with only command
- **WHEN** a deck declares an MCP with only a command field
- **THEN** the MCP SHALL be valid and written with just the command

### Requirement: Shell templates
The `shells` section SHALL be a map of shell name to session template. A shell template MAY contain any field valid on a session declaration: `agent`, `command`, `worktree`, `mcps`, `skills`, `prompt`.

#### Scenario: Session references shell
- **WHEN** a session declares `shell: my-shell`
- **THEN** the session SHALL inherit all fields from the referenced shell template

#### Scenario: Session overrides shell field
- **WHEN** a session declares `shell: my-shell` and also declares `agent: gemini`
- **THEN** the session's `agent` field SHALL override the shell's `agent` field

#### Scenario: Session merges MCP lists with shell
- **WHEN** a shell declares `mcps: [github]` and a session using that shell declares `mcps: [exa]`
- **THEN** the resolved session SHALL have `mcps: [github, exa]` (union, not override)

#### Scenario: Shell references nonexistent
- **WHEN** a session references a shell name that does not exist in the `shells` section
- **THEN** validation SHALL fail with an error identifying the missing shell

### Requirement: Group declarations
The `groups` section SHALL be a map of group name to group definition. A group definition MAY contain `sessions` (map) and `conductors` (map). An empty group (`{}` or `sessions: {}`) SHALL be valid.

#### Scenario: Empty group
- **WHEN** a deck declares a group with no sessions or conductors
- **THEN** the group SHALL be created in agent-deck with no sessions

#### Scenario: Group with sessions and conductors
- **WHEN** a deck declares a group with both `sessions` and `conductors` sections
- **THEN** both sessions and conductors SHALL be created within that group

### Requirement: Session declarations
A session declaration SHALL support: `agent` (string), `command` (string), `shell` (string), `worktree` (string), `mcps` (string array), `skills` (string array), `prompt` (string — inline or file path), `path` (string), `parent` (string). Exactly one of `agent`, `command`, or `shell` MUST be present.

#### Scenario: Agent session
- **WHEN** a session declares `agent: claude`
- **THEN** shuffle SHALL create the session with `agent-deck launch -c claude`

#### Scenario: Command session
- **WHEN** a session declares `command: ralph-tui run --parallel 3`
- **THEN** shuffle SHALL create the session with `agent-deck launch -c "ralph-tui run --parallel 3"`

#### Scenario: Session with parent
- **WHEN** a session declares `parent: my-conductor`
- **THEN** shuffle SHALL create the session with `--parent "my-conductor"`

#### Scenario: Session with worktree
- **WHEN** a session declares `worktree: subdirectory`
- **THEN** shuffle SHALL create the session with agent-deck's worktree support using the specified location strategy

#### Scenario: Neither agent nor command nor shell
- **WHEN** a session declares none of `agent`, `command`, or `shell`
- **THEN** validation SHALL fail

### Requirement: Conductor declarations
A conductor declaration SHALL support: `description` (string), `heartbeat` (boolean, default true), `claude_md` (string — inline or file path), `policy_md` (string — inline or file path).

#### Scenario: Conductor with inline claude_md
- **WHEN** a conductor declares `claude_md` as a YAML multiline string
- **THEN** shuffle SHALL write the content to the conductor's CLAUDE.md file in `~/.agent-deck/conductor/<name>/`

#### Scenario: Conductor with file-ref claude_md
- **WHEN** a conductor declares `claude_md: ./conductors/my-conductor.md`
- **THEN** shuffle SHALL pass the resolved absolute path to `conductor setup -claude-md`

#### Scenario: Conductor placed in group
- **WHEN** a conductor is declared under `groups.my-group.conductors`
- **THEN** shuffle SHALL run `conductor setup` then `group move` to place the conductor's session in "my-group"

### Requirement: Markdown field resolution
Fields `claude_md`, `policy_md`, and `prompt` SHALL accept either a file path or inline content. A value SHALL be treated as a file path if it contains `/` or ends with `.md`. Otherwise it SHALL be treated as inline content.

#### Scenario: Inline prompt
- **WHEN** a session declares `prompt: "Start the monitoring loop"`
- **THEN** the value SHALL be treated as inline content and sent via `launch -m`

#### Scenario: File-ref prompt
- **WHEN** a session declares `prompt: ./prompts/monitor.md`
- **THEN** shuffle SHALL read the file contents and send them via `launch -m`

#### Scenario: Multiline inline
- **WHEN** a field uses YAML `|` block scalar syntax
- **THEN** the value SHALL be treated as inline content (the `|` produces a string, not a path)
