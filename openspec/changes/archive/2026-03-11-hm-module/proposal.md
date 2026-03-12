## Why

Agent-deck's `~/.agent-deck/config.toml` must be hand-edited or managed ad-hoc. Users on NixOS or nix-darwin who manage their dotfiles through Home Manager have no way to declare agent-deck configuration alongside the rest of their system. A Home Manager module would let them version-control global settings, profile-level overrides, MCP definitions, and custom tools in Nix — with type checking, cross-platform support (Linux + Darwin), and integration with the broader Nix ecosystem (e.g., referencing packages, sops-nix for secrets).

## What Changes

- Create a new repository (`github_afterthought/hm-agent-deck`) containing a Home Manager module for agent-deck
- Provide fully-typed Nix options for stable config sections: `shell`, `claude`, `codex`, `docker`, `logs`, `updates`, `globalSearch`, `mcpPool`
- Provide flexible `attrsOf` options for dynamic sections: `mcps`, `tools`, `profiles`
- Generate `~/.agent-deck/config.toml` from the declared options
- Do NOT manage secrets inline — document that users should use `envFile` references pointing to sops-nix/agenix-managed files
- Support only `profiles.<name>.claude` for now (matching current agent-deck surface)
- Reference agent-deck binary from `llm-agents.nix` overlay/flake, not nixpkgs

## Capabilities

### New Capabilities
- `hm-module`: The Home Manager module — option types, TOML generation, activation, and cross-platform support
- `hm-testing`: Module test suite — NixOS test or standalone build checks validating generated config against known-good TOML

### Modified Capabilities

## Impact

- New repository: `github_afterthought/hm-agent-deck`
- Depends on Home Manager's module system (`mkOption`, `types`, `config.home.file`)
- Depends on `llm-agents.nix` for the agent-deck package
- No changes to agent-deck itself or to shuffle
- Users add the flake as an input and import the module in their Home Manager config
