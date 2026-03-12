## Context

Agent-deck stores its configuration in `~/.agent-deck/config.toml`. The file has a mix of stable, well-defined sections (shell, claude, docker, logs, updates, global_search, mcp_pool) and open-ended sections (mcps, tools, profiles) that grow as users add integrations. Users on NixOS/nix-darwin manage dotfiles via Home Manager and want agent-deck config to follow the same pattern.

The agent-deck binary itself comes from the `llm-agents.nix` flake/overlay — it's not in nixpkgs. The module needs to integrate with that package source.

## Goals / Non-Goals

**Goals:**
- Declare agent-deck global and profile-level config in Nix, generating `~/.agent-deck/config.toml`
- Fully type stable config sections for validation and discoverability
- Use flexible `attrsOf` typing for dynamic sections (mcps, tools) to avoid coupling to every upstream field
- Work identically on Linux (NixOS/Home Manager) and Darwin (nix-darwin/Home Manager)
- Provide clear guidance for secrets management via env_file indirection

**Non-Goals:**
- Managing secrets inline (API keys in Nix store)
- Managing `sessions.json`, conductor files, or anything shuffle/deal handles
- Managing `skills/sources.toml` (future work)
- Installing agent-deck — the module configures it, the user brings the package from `llm-agents.nix`
- Supporting profile-level overrides beyond `claude` (only section agent-deck currently supports)

## Decisions

### D1: Separate repository
**Choice**: New repo `github_afterthought/hm-agent-deck`, not part of shuffle or agent-deck.
**Rationale**: The module has its own release cadence (tied to agent-deck config schema changes, not shuffle features). It's a Nix flake with no Go code. Keeping it separate avoids polluting shuffle's build and lets Nix users consume it without pulling in the shuffle CLI.

### D2: Hybrid typing — structured stable, flexible dynamic
**Choice**: Fully typed submodules for `shell`, `claude`, `codex`, `docker`, `logs`, `updates`, `globalSearch`, `mcpPool`. Dynamic `attrsOf` for `mcps`, `tools`, `profiles`.
**Rationale**: Stable sections change rarely and benefit from type errors and documentation. Dynamic sections vary per-user and per-upstream-release — rigid typing would require module updates for every new MCP field. The hybrid gives safety where it matters and flexibility where it's needed.
**Alternative considered**: Fully typed everything — rejected because `mcps` fields evolve (e.g., `server` sub-section, `headers`, `transport` were added over time). Also considered pure freeform `settings` attr — rejected because it loses all validation benefit.

### D3: TOML generation via `pkgs.formats.toml`
**Choice**: Use Nix's `pkgs.formats.toml` to generate the config file from an attrset.
**Rationale**: Home Manager convention. Handles TOML serialization correctly (arrays, nested tables, inline tables). The module builds a merged attrset from typed options + freeform options, then passes it to `formats.toml.generate`.

### D4: Option namespace `programs.agent-deck`
**Choice**: `programs.agent-deck.enable`, `programs.agent-deck.settings.*`, `programs.agent-deck.mcps.*`, etc.
**Rationale**: Follows Home Manager convention (`programs.git`, `programs.tmux`). The `settings` sub-attr holds the typed global sections. `mcps` and `tools` are top-level under `programs.agent-deck` for ergonomics.

### D5: No package management in module
**Choice**: The module does NOT install agent-deck. Users add `llm-agents.nix` to their flake inputs and include agent-deck in `home.packages` themselves.
**Rationale**: Agent-deck comes from a separate flake (`llm-agents.nix`), not nixpkgs. Having the module depend on that flake creates a transitive dependency. Keeping them separate lets users pin versions independently. The module can optionally accept a `package` option for path references but won't add it to PATH.

### D6: Config file placement via `home.file`
**Choice**: Write generated TOML to `~/.agent-deck/config.toml` using `home.file."agent-deck/config.toml"`.
**Rationale**: Standard Home Manager mechanism. Symlinks the generated file from the Nix store. Agent-deck reads this path by default — no env var needed.
**Caveat**: The symlink is read-only. If agent-deck ever writes to config.toml at runtime, this breaks. Currently agent-deck only reads it, so this is fine. Document this.

### D7: Nix option naming — camelCase to snake_case mapping
**Choice**: Nix options use camelCase (Nix convention). The TOML generator maps to snake_case keys (TOML convention). For example, `globalSearch.memoryLimitMb` → `[global_search] memory_limit_mb`.
**Rationale**: Each ecosystem follows its own convention. The mapping is mechanical and lives in the TOML generation layer, not in option definitions. Users write idiomatic Nix; the output is idiomatic TOML.

### D8: Profiles as nested submodule
**Choice**: `programs.agent-deck.profiles` is `attrsOf (submodule { claude = { configDir = ...; }; })`.
**Rationale**: Currently agent-deck only supports `profiles.<name>.claude.config_dir`. Modeling it as a submodule means we can add more profile-level sections later without restructuring. For now the submodule only has `claude`.

## Risks / Trade-offs

- **[Read-only symlink]** → If agent-deck adds runtime config writes, the symlinked file will cause errors. Mitigation: monitor upstream, switch to `home.activation` copy if needed.
- **[Schema drift]** → Typed options may lag behind agent-deck releases. Mitigation: freeform `extraConfig` escape hatch that merges raw attrsets into the generated TOML.
- **[camelCase mapping bugs]** → Mechanical name translation could miss edge cases (e.g., `mcpPool` → `mcp_pool` vs `mcps` → `mcps`). Mitigation: explicit mapping table in the generator, not a generic camelToSnake function.
- **[No secret management]** → Users must handle API keys outside the module. Mitigation: document the pattern (sops-nix → env_file path) and provide examples in README.
