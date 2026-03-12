package deck

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load reads and parses a .deck.yaml file.
func Load(path string) (*Deck, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading deck file: %w", err)
	}

	// First unmarshal into a raw map to check for unknown keys.
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing deck YAML: %w", err)
	}

	knownKeys := map[string]bool{
		"name": true, "profile": true, "mcps": true,
		"shells": true, "groups": true,
	}
	for key := range raw {
		if !knownKeys[key] {
			return nil, fmt.Errorf("unknown top-level key: %q", key)
		}
	}

	var d Deck
	if err := yaml.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("parsing deck YAML: %w", err)
	}

	return &d, nil
}

// Validate checks the deck for semantic errors.
func Validate(d *Deck) []error {
	var errs []error

	if d.Name == "" {
		errs = append(errs, fmt.Errorf("deck 'name' is required"))
	}

	if d.Profile != nil && d.Profile.Name == "" {
		errs = append(errs, fmt.Errorf("profile 'name' is required when profile is declared"))
	}

	for mcpName, mcp := range d.MCPs {
		if mcp.Command == "" {
			errs = append(errs, fmt.Errorf("mcp %q: 'command' is required", mcpName))
		}
	}

	for groupName, group := range d.Groups {
		for sessName, sess := range group.Sessions {
			count := 0
			if sess.Agent != "" {
				count++
			}
			if sess.Command != "" {
				count++
			}
			if sess.Shell != "" {
				count++
			}
			if count == 0 {
				errs = append(errs, fmt.Errorf("group %q session %q: must declare 'agent', 'command', or 'shell'", groupName, sessName))
			}
			if count > 1 {
				errs = append(errs, fmt.Errorf("group %q session %q: 'agent', 'command', and 'shell' are mutually exclusive", groupName, sessName))
			}
			if sess.Shell != "" {
				if d.Shells == nil {
					errs = append(errs, fmt.Errorf("group %q session %q: shell %q not found (no shells defined)", groupName, sessName, sess.Shell))
				} else if _, ok := d.Shells[sess.Shell]; !ok {
					errs = append(errs, fmt.Errorf("group %q session %q: shell %q not found", groupName, sessName, sess.Shell))
				}
			}
			for _, mcp := range sess.MCPs {
				if d.MCPs != nil {
					if _, ok := d.MCPs[mcp]; !ok {
						errs = append(errs, fmt.Errorf("group %q session %q: mcp %q not declared in mcps section", groupName, sessName, mcp))
					}
				}
			}
		}
	}

	return errs
}

// ResolveShells expands shell references in all sessions.
// After this, every session will have agent or command set directly.
func ResolveShells(d *Deck) error {
	for groupName, group := range d.Groups {
		for sessName, sess := range group.Sessions {
			if sess.Shell == "" {
				continue
			}
			shell, ok := d.Shells[sess.Shell]
			if !ok {
				return fmt.Errorf("group %q session %q: shell %q not found", groupName, sessName, sess.Shell)
			}

			// Merge shell into session. Session fields override, MCPs merge.
			if sess.Agent == "" {
				sess.Agent = shell.Agent
			}
			if sess.Command == "" {
				sess.Command = shell.Command
			}
			if sess.Worktree == "" {
				sess.Worktree = shell.Worktree
			}
			if sess.Prompt == "" {
				sess.Prompt = shell.Prompt
			}

			// MCPs: union merge
			mcpSet := make(map[string]bool)
			for _, m := range shell.MCPs {
				mcpSet[m] = true
			}
			for _, m := range sess.MCPs {
				mcpSet[m] = true
			}
			merged := make([]string, 0, len(mcpSet))
			// Preserve order: shell MCPs first, then session MCPs
			for _, m := range shell.MCPs {
				if mcpSet[m] {
					merged = append(merged, m)
					delete(mcpSet, m)
				}
			}
			for _, m := range sess.MCPs {
				if mcpSet[m] {
					merged = append(merged, m)
					delete(mcpSet, m)
				}
			}
			sess.MCPs = merged

			// Skills: union merge (same logic)
			skillSet := make(map[string]bool)
			for _, s := range shell.Skills {
				skillSet[s] = true
			}
			for _, s := range sess.Skills {
				skillSet[s] = true
			}
			mergedSkills := make([]string, 0, len(skillSet))
			for _, s := range shell.Skills {
				if skillSet[s] {
					mergedSkills = append(mergedSkills, s)
					delete(skillSet, s)
				}
			}
			for _, s := range sess.Skills {
				if skillSet[s] {
					mergedSkills = append(mergedSkills, s)
					delete(skillSet, s)
				}
			}
			sess.Skills = mergedSkills

			sess.Shell = "" // Clear shell reference
			group.Sessions[sessName] = sess
		}
		d.Groups[groupName] = group
	}
	return nil
}

// IsFilePath returns true if the value looks like a file path (contains / or ends in .md).
func IsFilePath(value string) bool {
	return strings.Contains(value, "/") || strings.HasSuffix(value, ".md")
}

// ResolveMarkdownField resolves a field value that may be inline content or a file path.
// basePath is the directory to resolve relative paths from.
func ResolveMarkdownField(value, basePath string) (string, error) {
	if value == "" {
		return "", nil
	}
	if !IsFilePath(value) {
		return value, nil
	}
	// It's a file path
	path := value
	if !filepath.IsAbs(path) {
		path = filepath.Join(basePath, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading markdown file %q: %w", path, err)
	}
	return string(data), nil
}
