# Capability: hm-testing

## Purpose
Test suite for the Home Manager module, verifying that Nix option declarations produce correct TOML output.

## Requirements

### Requirement: Basic module evaluation test
The test suite SHALL include a test that evaluates the module with minimal configuration (`enable = true`) and verifies it produces valid TOML output.

#### Scenario: Minimal config evaluates
- **WHEN** the module is evaluated with only `programs.agent-deck.enable = true`
- **THEN** the evaluation succeeds and produces a non-empty TOML file

### Requirement: Full config round-trip test
The test suite SHALL include a test that sets every typed option and verifies the generated TOML contains all expected keys with correct snake_case names.

#### Scenario: All sections present
- **WHEN** every typed option is set to a non-default value
- **THEN** the generated TOML contains `[shell]`, `[claude]`, `[codex]`, `[docker]`, `[logs]`, `[updates]`, `[global_search]`, `[mcp_pool]` sections with correct key names

### Requirement: MCP definition test
The test suite SHALL verify that freeform MCP definitions produce correct `[mcps.*]` TOML sections.

#### Scenario: STDIO MCP round-trip
- **WHEN** an MCP is defined with `command`, `args`, and `env`
- **THEN** the generated TOML contains a `[mcps.<name>]` section with those exact values

#### Scenario: HTTP MCP round-trip
- **WHEN** an MCP is defined with `url`, `transport`, and `headers`
- **THEN** the generated TOML contains a `[mcps.<name>]` section with those exact values

### Requirement: Profile override test
The test suite SHALL verify that profile-level Claude overrides produce correct `[profiles.*]` TOML sections.

#### Scenario: Profile Claude config dir
- **WHEN** a profile `work` with `claude.configDir` is defined
- **THEN** the generated TOML contains `[profiles.work.claude]` with `config_dir`

### Requirement: Extra config merge test
The test suite SHALL verify that `extraConfig` values merge correctly with typed options.

#### Scenario: Extra config does not clobber typed values
- **WHEN** typed `claude.configDir` is set AND `extraConfig.claude.new_key` is set
- **THEN** the generated TOML `[claude]` section contains both `config_dir` and `new_key`

### Requirement: Empty section omission test
The test suite SHALL verify that sections with all-default values are omitted from output.

#### Scenario: No empty sections
- **WHEN** only `enable` and `defaultTool` are set
- **THEN** the generated TOML does not contain `[shell]`, `[docker]`, or other unconfigured sections
