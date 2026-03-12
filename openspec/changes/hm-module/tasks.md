## 1. Repository Scaffolding

- [ ] 1.1 Create `github_afterthought/hm-agent-deck` repo with flake.nix exposing a Home Manager module
- [ ] 1.2 Set up flake inputs: nixpkgs, home-manager, flake-utils
- [ ] 1.3 Create directory structure: `modules/`, `tests/`

## 2. Option Definitions

- [ ] 2.1 Define `programs.agent-deck.enable` toggle and `defaultTool` option
- [ ] 2.2 Define `shell` submodule options (envFiles, initScript, ignoreMissingEnvFiles)
- [ ] 2.3 Define `claude` submodule options (configDir, dangerousMode, allowDangerousMode, envFile)
- [ ] 2.4 Define `codex` submodule options (yoloMode)
- [ ] 2.5 Define `docker` submodule options (defaultEnabled, defaultImage, cpuLimit, memoryLimit, mountSsh, autoCleanup, environment, volumeIgnores)
- [ ] 2.6 Define `logs` submodule options (maxSizeMb, maxLines, removeOrphans)
- [ ] 2.7 Define `updates` submodule options (autoUpdate, checkEnabled, checkIntervalHours, notifyInCli)
- [ ] 2.8 Define `globalSearch` submodule options (enabled, tier enum, memoryLimitMb, recentDays, indexRateLimit)
- [ ] 2.9 Define `mcpPool` submodule options (enabled, autoStart, poolAll, excludeMcps, fallbackToStdio, showPoolStatus)
- [ ] 2.10 Define `mcps` as attrsOf freeform attrs
- [ ] 2.11 Define `tools` as attrsOf freeform attrs
- [ ] 2.12 Define `profiles` as attrsOf submodule with claude.configDir
- [ ] 2.13 Define `extraConfig` freeform escape hatch

## 3. TOML Generation

- [ ] 3.1 Build camelCase-to-snake_case mapping table for all section and key names
- [ ] 3.2 Implement config attrset builder that merges typed options, mcps, tools, profiles, and extraConfig
- [ ] 3.3 Filter out sections where all values are defaults/null (omit empty sections)
- [ ] 3.4 Generate `~/.agent-deck/config.toml` via `home.file` using `pkgs.formats.toml`

## 4. Tests

- [ ] 4.1 Create minimal config evaluation test (enable-only produces valid TOML)
- [ ] 4.2 Create full config round-trip test (all typed options set, verify snake_case keys)
- [ ] 4.3 Create MCP definition test (STDIO and HTTP variants)
- [ ] 4.4 Create profile override test (profiles.work.claude.configDir)
- [ ] 4.5 Create extraConfig merge test (typed + extra don't clobber)
- [ ] 4.6 Create empty section omission test

## 5. Documentation

- [ ] 5.1 Write README with usage examples: minimal, full config, MCP definitions, profiles, secrets pattern
- [ ] 5.2 Add example flake.nix showing how to consume the module alongside llm-agents.nix
