# Capability: hm-module

## Purpose
Home Manager module for declarative agent-deck configuration. Provides typed Nix options that generate `~/.agent-deck/config.toml`.

## Requirements

### Requirement: Module enable toggle
The module SHALL provide a `programs.agent-deck.enable` option (bool, default `false`) that controls whether agent-deck configuration is managed.

#### Scenario: Module disabled
- **WHEN** `programs.agent-deck.enable` is `false`
- **THEN** no config files are written to `~/.agent-deck/`

#### Scenario: Module enabled
- **WHEN** `programs.agent-deck.enable` is `true`
- **THEN** `~/.agent-deck/config.toml` is generated from declared options

### Requirement: Top-level default tool option
The module SHALL provide `programs.agent-deck.defaultTool` (string, default `"claude"`) mapping to the `default_tool` TOML key.

#### Scenario: Custom default tool
- **WHEN** `programs.agent-deck.defaultTool = "codex"`
- **THEN** generated TOML contains `default_tool = "codex"`

### Requirement: Shell section options
The module SHALL provide typed options under `programs.agent-deck.shell` for `envFiles` (list of string), `initScript` (string), and `ignoreMissingEnvFiles` (bool). All section options default to `null` (unset) — only explicitly set values appear in generated TOML, allowing agent-deck to apply its own defaults.

#### Scenario: Shell env files configured
- **WHEN** `programs.agent-deck.shell.envFiles = ["~/.agent-deck.env" ".env"]`
- **THEN** generated TOML contains `[shell]` with `env_files = ["~/.agent-deck.env", ".env"]`

#### Scenario: Shell section omitted when defaults
- **WHEN** no shell options are set (all defaults/null)
- **THEN** the `[shell]` section is omitted from generated TOML

### Requirement: Claude section options
The module SHALL provide typed options under `programs.agent-deck.claude` for `configDir` (string), `dangerousMode` (bool), `allowDangerousMode` (bool), and `envFile` (string). All default to `null`.

#### Scenario: Claude config dir override
- **WHEN** `programs.agent-deck.claude.configDir = "~/.claude-work"`
- **THEN** generated TOML contains `[claude]` with `config_dir = "~/.claude-work"`

### Requirement: Codex section options
The module SHALL provide typed options under `programs.agent-deck.codex` for `yoloMode` (bool). Defaults to `null`.

#### Scenario: Codex yolo mode enabled
- **WHEN** `programs.agent-deck.codex.yoloMode = true`
- **THEN** generated TOML contains `[codex]` with `yolo_mode = true`

### Requirement: Docker section options
The module SHALL provide typed options under `programs.agent-deck.docker` for `defaultEnabled` (bool), `defaultImage` (string), `cpuLimit` (string), `memoryLimit` (string), `mountSsh` (bool), `autoCleanup` (bool), `environment` (list of string), and `volumeIgnores` (list of string). All default to `null`.

#### Scenario: Docker sandbox defaults
- **WHEN** `programs.agent-deck.docker.mountSsh = true` and `programs.agent-deck.docker.autoCleanup = true`
- **THEN** generated TOML contains `[docker]` with `mount_ssh = true` and `auto_cleanup = true`

### Requirement: Logs section options
The module SHALL provide typed options under `programs.agent-deck.logs` for `maxSizeMb` (int), `maxLines` (int), and `removeOrphans` (bool). All default to `null`.

#### Scenario: Custom log limits
- **WHEN** `programs.agent-deck.logs.maxSizeMb = 50`
- **THEN** generated TOML contains `[logs]` with `max_size_mb = 50`

### Requirement: Updates section options
The module SHALL provide typed options under `programs.agent-deck.updates` for `autoUpdate` (bool), `checkEnabled` (bool), `checkIntervalHours` (int), and `notifyInCli` (bool). All default to `null`.

#### Scenario: Disable update checks
- **WHEN** `programs.agent-deck.updates.checkEnabled = false`
- **THEN** generated TOML contains `[updates]` with `check_enabled = false`

### Requirement: Global search section options
The module SHALL provide typed options under `programs.agent-deck.globalSearch` for `enabled` (bool), `tier` (enum `["auto" "instant" "balanced"]`), `memoryLimitMb` (int), `recentDays` (int), and `indexRateLimit` (int). All default to `null`.

#### Scenario: Search tier override
- **WHEN** `programs.agent-deck.globalSearch.tier = "instant"`
- **THEN** generated TOML contains `[global_search]` with `tier = "instant"`

### Requirement: MCP pool section options
The module SHALL provide typed options under `programs.agent-deck.mcpPool` for `enabled` (bool), `autoStart` (bool), `poolAll` (bool), `excludeMcps` (list of string), `fallbackToStdio` (bool), and `showPoolStatus` (bool). All default to `null`.

#### Scenario: MCP pool enabled
- **WHEN** `programs.agent-deck.mcpPool.enabled = true` and `programs.agent-deck.mcpPool.poolAll = true`
- **THEN** generated TOML contains `[mcp_pool]` with `enabled = true` and `pool_all = true`

### Requirement: MCP definitions (flexible)
The module SHALL provide `programs.agent-deck.mcps` as `attrsOf (attrsOf unspecified)` allowing arbitrary MCP definitions that map directly to `[mcps.*]` TOML sections.

#### Scenario: STDIO MCP definition
- **WHEN** user declares `programs.agent-deck.mcps.github = { command = "npx"; args = ["-y" "@modelcontextprotocol/server-github"]; env = { GITHUB_TOKEN = "ghp_xxx"; }; }`
- **THEN** generated TOML contains `[mcps.github]` with `command`, `args`, and `env` keys

#### Scenario: HTTP MCP definition
- **WHEN** user declares `programs.agent-deck.mcps.remote = { url = "https://api.example.com/mcp"; transport = "http"; }`
- **THEN** generated TOML contains `[mcps.remote]` with `url` and `transport` keys

### Requirement: Custom tools definitions (flexible)
The module SHALL provide `programs.agent-deck.tools` as `attrsOf (attrsOf unspecified)` allowing arbitrary tool definitions that map directly to `[tools.*]` TOML sections.

#### Scenario: Custom tool definition
- **WHEN** user declares `programs.agent-deck.tools.my-ai = { command = "my-ai-assistant"; icon = "🧠"; }`
- **THEN** generated TOML contains `[tools.my-ai]` with `command` and `icon` keys

### Requirement: Profile-level Claude overrides
The module SHALL provide `programs.agent-deck.profiles` as `attrsOf (submodule)` where each profile has a `claude.configDir` option (string).

#### Scenario: Profile Claude config dir
- **WHEN** user declares `programs.agent-deck.profiles.work.claude.configDir = "~/.claude-work"`
- **THEN** generated TOML contains `[profiles.work.claude]` with `config_dir = "~/.claude-work"`

#### Scenario: Multiple profiles
- **WHEN** user declares profiles `work` and `clientx` with different `claude.configDir` values
- **THEN** generated TOML contains both `[profiles.work.claude]` and `[profiles.clientx.claude]` sections

### Requirement: Extra config escape hatch
The module SHALL provide `programs.agent-deck.extraConfig` (attrs, default `{}`) that is merged into the generated TOML attrset, allowing users to set options not covered by typed options.

#### Scenario: Unknown future option
- **WHEN** user sets `programs.agent-deck.extraConfig = { some_new_section = { key = "value"; }; }`
- **THEN** generated TOML contains `[some_new_section]` with `key = "value"`

#### Scenario: Extra config merges with typed config
- **WHEN** user sets both `programs.agent-deck.claude.configDir = "~/.claude"` and `programs.agent-deck.extraConfig.claude.new_field = true`
- **THEN** generated TOML `[claude]` section contains both `config_dir` and `new_field`

### Requirement: Config file generation
The module SHALL generate `~/.agent-deck/config.toml` using `home.file` and the `pkgs.formats.toml` generator.

#### Scenario: Generated file location
- **WHEN** module is enabled
- **THEN** `~/.agent-deck/config.toml` exists as a symlink to a Nix store path containing valid TOML

### Requirement: Omit default/empty sections
The module SHALL omit TOML sections where all values are defaults or unset, producing minimal config files.

#### Scenario: Minimal config
- **WHEN** only `enable = true` and `defaultTool = "claude"` are set
- **THEN** generated TOML contains only `default_tool = "claude"` with no empty sections

### Requirement: Cross-platform support
The module SHALL work on both Linux (NixOS + Home Manager) and Darwin (nix-darwin + Home Manager) without platform-specific conditionals in user configuration.

#### Scenario: Linux usage
- **WHEN** module is used on a NixOS system with Home Manager
- **THEN** config is generated at `~/.agent-deck/config.toml` with no errors

#### Scenario: Darwin usage
- **WHEN** module is used on a macOS system with nix-darwin + Home Manager
- **THEN** config is generated at `~/.agent-deck/config.toml` with no errors
