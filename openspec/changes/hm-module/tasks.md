## 1. Repository Scaffolding

- [x] 1.1 Create `github_afterthought/nix-agent-deck` repo with flake.nix exposing a Home Manager module
- [x] 1.2 Set up flake inputs: nixpkgs, home-manager, flake-utils
- [x] 1.3 Create directory structure: `modules/`, `tests/`

## 2. Option Definitions

- [x] 2.1 Define `programs.agent-deck.enable` toggle and `defaultTool` option
- [x] 2.2 Define `shell` submodule options (envFiles, initScript, ignoreMissingEnvFiles)
- [x] 2.3 Define `claude` submodule options (configDir, dangerousMode, allowDangerousMode, envFile)
- [x] 2.4 Define `codex` submodule options (yoloMode)
- [x] 2.5 Define `docker` submodule options (defaultEnabled, defaultImage, cpuLimit, memoryLimit, mountSsh, autoCleanup, environment, volumeIgnores)
- [x] 2.6 Define `logs` submodule options (maxSizeMb, maxLines, removeOrphans)
- [x] 2.7 Define `updates` submodule options (autoUpdate, checkEnabled, checkIntervalHours, notifyInCli)
- [x] 2.8 Define `globalSearch` submodule options (enabled, tier enum, memoryLimitMb, recentDays, indexRateLimit)
- [x] 2.9 Define `mcpPool` submodule options (enabled, autoStart, poolAll, excludeMcps, fallbackToStdio, showPoolStatus)
- [x] 2.10 Define `mcps` as attrsOf freeform attrs
- [x] 2.11 Define `tools` as attrsOf freeform attrs
- [x] 2.12 Define `profiles` as attrsOf submodule with claude.configDir
- [x] 2.13 Define `extraConfig` freeform escape hatch

## 3. TOML Generation

- [x] 3.1 Build camelCase-to-snake_case mapping table for all section and key names
- [x] 3.2 Implement config attrset builder that merges typed options, mcps, tools, profiles, and extraConfig
- [x] 3.3 Filter out sections where all values are defaults/null (omit empty sections)
- [x] 3.4 Generate `~/.agent-deck/config.toml` via `home.file` using `pkgs.formats.toml`

## 4. Tests

- [x] 4.1 Create minimal config evaluation test (enable-only produces valid TOML)
- [x] 4.2 Create full config round-trip test (all typed options set, verify snake_case keys)
- [x] 4.3 Create MCP definition test (STDIO and HTTP variants)
- [x] 4.4 Create profile override test (profiles.work.claude.configDir)
- [x] 4.5 Create extraConfig merge test (typed + extra don't clobber)
- [x] 4.6 Create empty section omission test

## 5. Documentation

- [x] 5.1 Write README with usage examples: minimal, full config, MCP definitions, profiles, secrets pattern
- [x] 5.2 Add example flake.nix showing how to consume the module alongside llm-agents.nix
