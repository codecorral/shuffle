package state

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Current represents the current agent-deck state.
type Current struct {
	Profiles   []string
	Groups     []string
	Sessions   []SessionInfo
	Conductors []ConductorInfo
	MCPs       []string // MCP names from config.toml
}

type SessionInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Group string `json:"group"`
}

type ConductorInfo struct {
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

// Discover queries agent-deck to build the current state.
func Discover(profile string) (*Current, error) {
	cur := &Current{}

	// Profiles
	profiles, err := runAgentDeckJSON("profile", "list", "--json")
	if err != nil {
		// Profile list may fail if no profiles; treat as empty
		cur.Profiles = []string{}
	} else {
		cur.Profiles = parseProfileList(profiles)
	}

	// Groups — use paths (not names) since sessions reference groups by path
	args := buildProfileArgs(profile, "group", "list", "--json")
	groups, err := runAgentDeckJSONArgs(args)
	if err != nil {
		cur.Groups = []string{}
	} else {
		cur.Groups = parseGroupPaths(groups)
	}

	// Sessions
	args = buildProfileArgs(profile, "list", "--json")
	sessions, err := runAgentDeckJSONArgs(args)
	if err != nil {
		cur.Sessions = []SessionInfo{}
	} else {
		cur.Sessions = parseSessionList(sessions)
	}

	// Conductors
	args = buildProfileArgs(profile, "conductor", "list", "--json")
	conductors, err := runAgentDeckJSONArgs(args)
	if err != nil {
		cur.Conductors = []ConductorInfo{}
	} else {
		cur.Conductors = parseConductorList(conductors)
	}

	// MCPs from config.toml
	mcpNames, err := readMCPsFromConfig()
	if err != nil {
		cur.MCPs = []string{}
	} else {
		cur.MCPs = mcpNames
	}

	return cur, nil
}

// AttachedMCPs returns the MCPs attached to a session.
func AttachedMCPs(profile, sessionID string) ([]string, error) {
	args := buildProfileArgs(profile, "mcp", "attached", sessionID, "--json")
	out, err := runAgentDeckJSONArgs(args)
	if err != nil {
		return nil, err
	}
	return parseStringList(out, "name"), nil
}

func buildProfileArgs(profile string, args ...string) []string {
	if profile != "" {
		return append([]string{"-p", profile}, args...)
	}
	return args
}

func runAgentDeckJSON(args ...string) ([]byte, error) {
	return runAgentDeckJSONArgs(args)
}

func runAgentDeckJSONArgs(args []string) ([]byte, error) {
	cmd := exec.Command("agent-deck", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("agent-deck %s: %w", strings.Join(args, " "), err)
	}
	return out, nil
}

func parseStringList(data []byte, key string) []string {
	// Try as array of objects first (e.g., [{"name": "foo"}, ...])
	var items []map[string]interface{}
	if err := json.Unmarshal(data, &items); err == nil {
		var result []string
		for _, item := range items {
			if v, ok := item[key]; ok {
				result = append(result, fmt.Sprintf("%v", v))
			}
		}
		return result
	}
	// Try as array of strings
	var strs []string
	if err := json.Unmarshal(data, &strs); err == nil {
		return strs
	}
	// Try as wrapper object with a known list key (agent-deck wraps lists in objects)
	var wrapper map[string]json.RawMessage
	if err := json.Unmarshal(data, &wrapper); err == nil {
		// Look for a key whose value is an array containing objects with our target key
		for _, raw := range wrapper {
			var nested []map[string]interface{}
			if err := json.Unmarshal(raw, &nested); err == nil && len(nested) > 0 {
				if _, ok := nested[0][key]; ok {
					var result []string
					for _, item := range nested {
						if v, ok := item[key]; ok {
							result = append(result, fmt.Sprintf("%v", v))
						}
					}
					return result
				}
			}
		}
	}
	return nil
}

// parseProfileList extracts profile names from agent-deck's profile list JSON.
// Format: {"profiles": [{"name": "foo", "is_default": true}, ...], "default_profile": "...", ...}
func parseProfileList(data []byte) []string {
	return parseStringList(data, "name")
}

// parseGroupPaths extracts all group paths (flattened, including children) from agent-deck's group list JSON.
// Format: {"groups": [{"name": "...", "path": "...", "children": [...]}, ...]}
func parseGroupPaths(data []byte) []string {
	var wrapper struct {
		Groups []groupNode `json:"groups"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		// Fallback: try as flat array
		var flat []groupNode
		if err := json.Unmarshal(data, &flat); err != nil {
			return nil
		}
		wrapper.Groups = flat
	}
	var paths []string
	for _, g := range wrapper.Groups {
		paths = flattenGroupPaths(g, paths)
	}
	return paths
}

type groupNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Children []groupNode `json:"children"`
}

func flattenGroupPaths(g groupNode, paths []string) []string {
	paths = append(paths, g.Path)
	for _, child := range g.Children {
		paths = flattenGroupPaths(child, paths)
	}
	return paths
}

func parseSessionList(data []byte) []SessionInfo {
	var sessions []SessionInfo
	// Try direct unmarshal
	if err := json.Unmarshal(data, &sessions); err == nil {
		return sessions
	}
	// Try as generic map array
	var items []map[string]interface{}
	if err := json.Unmarshal(data, &items); err == nil {
		for _, item := range items {
			s := SessionInfo{}
			if v, ok := item["id"]; ok {
				s.ID = fmt.Sprintf("%v", v)
			}
			if v, ok := item["title"]; ok {
				s.Title = fmt.Sprintf("%v", v)
			}
			if v, ok := item["group"]; ok {
				s.Group = fmt.Sprintf("%v", v)
			}
			sessions = append(sessions, s)
		}
	}
	return sessions
}

func parseConductorList(data []byte) []ConductorInfo {
	var conductors []ConductorInfo
	if err := json.Unmarshal(data, &conductors); err == nil {
		return conductors
	}
	var items []map[string]interface{}
	if err := json.Unmarshal(data, &items); err == nil {
		for _, item := range items {
			c := ConductorInfo{}
			if v, ok := item["name"]; ok {
				c.Name = fmt.Sprintf("%v", v)
			}
			if v, ok := item["profile"]; ok {
				c.Profile = fmt.Sprintf("%v", v)
			}
			conductors = append(conductors, c)
		}
	}
	return conductors
}

func readMCPsFromConfig() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(home, ".agent-deck", "config.toml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	mcps, ok := config["mcps"]
	if !ok {
		return nil, nil
	}

	mcpMap, ok := mcps.(map[string]interface{})
	if !ok {
		return nil, nil
	}

	var names []string
	for name := range mcpMap {
		names = append(names, name)
	}
	return names, nil
}

// AttachedSkills returns the skills attached to a session.
func AttachedSkills(profile, sessionID string) ([]string, error) {
	args := buildProfileArgs(profile, "skill", "attached", sessionID, "--json")
	out, err := runAgentDeckJSONArgs(args)
	if err != nil {
		return nil, err
	}
	return parseStringList(out, "name"), nil
}

// DefaultProfile queries agent-deck for the current default profile name.
func DefaultProfile() (string, error) {
	out, err := runAgentDeckJSON("profile", "default", "--json")
	if err != nil {
		return "default", nil // fallback
	}
	var result map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		return "default", nil
	}
	if name, ok := result["default_profile"]; ok {
		return fmt.Sprintf("%v", name), nil
	}
	return "default", nil
}

// ConfigPath returns the path to agent-deck's config.toml.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agent-deck", "config.toml"), nil
}
