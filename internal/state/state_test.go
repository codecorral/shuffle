package state

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
)

// These test fixtures are captured from real agent-deck v0.24.x output.

const realProfileListJSON = `{
  "default_profile": "default",
  "profiles": [
    {"is_default": false, "name": "codecorral"},
    {"is_default": true, "name": "default"},
    {"is_default": false, "name": "willdan"}
  ],
  "success": true,
  "total": 3
}`

const realGroupListJSON = `{
  "groups": [
    {
      "name": "conductor",
      "path": "conductor",
      "session_count": 2,
      "status": {"running": 0, "waiting": 0, "idle": 2},
      "children": []
    },
    {
      "name": "My Sessions",
      "path": "my-sessions",
      "session_count": 5,
      "status": {"running": 0, "waiting": 2, "idle": 3},
      "children": [
        {
          "name": "openmob",
          "path": "my-sessions/openmob",
          "session_count": 3,
          "status": {"running": 0, "waiting": 1, "idle": 2}
        }
      ]
    },
    {
      "name": "codecorral",
      "path": "codecorral",
      "session_count": 2,
      "status": {"running": 0, "waiting": 1, "idle": 1}
    }
  ],
  "total_groups": 4,
  "total_sessions": 12
}`

const realSessionListJSON = `[
  {
    "id": "94079c52-1772748442",
    "title": "conductor-mob-forge",
    "path": "/Users/chuck/.agent-deck/conductor/mob-forge",
    "group": "conductor",
    "tool": "claude",
    "command": "claude",
    "status": "idle",
    "profile": "default"
  },
  {
    "id": "08000d32-1773184364",
    "title": "cc-1 elaboration",
    "path": "/Users/chuck/Code/github_afterthought/codecorral",
    "group": "codecorral",
    "tool": "claude",
    "command": "claude",
    "status": "waiting",
    "profile": "default"
  },
  {
    "id": "0172286f-1772678917",
    "title": "mob-forge",
    "path": "/Users/chuck/Code/github_afterthought/mob-forge",
    "group": "my-sessions/openmob",
    "tool": "claude",
    "command": "claude",
    "status": "waiting",
    "profile": "default"
  }
]`

func TestParseProfileList(t *testing.T) {
	result := parseProfileList([]byte(realProfileListJSON))
	if len(result) != 3 {
		t.Fatalf("expected 3 profiles, got %d: %v", len(result), result)
	}
	expected := map[string]bool{"codecorral": true, "default": true, "willdan": true}
	for _, name := range result {
		if !expected[name] {
			t.Errorf("unexpected profile: %q", name)
		}
	}
}

func TestParseGroupPaths(t *testing.T) {
	result := parseGroupPaths([]byte(realGroupListJSON))
	expected := []string{"conductor", "my-sessions", "my-sessions/openmob", "codecorral"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d groups, got %d: %v", len(expected), len(result), result)
	}
	for _, exp := range expected {
		found := false
		for _, got := range result {
			if got == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing group path: %q", exp)
		}
	}
}

func TestParseGroupPaths_IncludesChildren(t *testing.T) {
	result := parseGroupPaths([]byte(realGroupListJSON))
	found := false
	for _, p := range result {
		if p == "my-sessions/openmob" {
			found = true
		}
	}
	if !found {
		t.Error("child group path my-sessions/openmob not found in flattened list")
	}
}

func TestParseSessionList_GroupIsPath(t *testing.T) {
	sessions := parseSessionList([]byte(realSessionListJSON))
	if len(sessions) != 3 {
		t.Fatalf("expected 3 sessions, got %d", len(sessions))
	}
	if sessions[2].Group != "my-sessions/openmob" {
		t.Errorf("expected path-style group, got %q", sessions[2].Group)
	}
}

func TestParseStringList_WrappedObject(t *testing.T) {
	// Test the generic wrapper-object fallback
	data := `{"items": [{"name": "a"}, {"name": "b"}], "total": 2}`
	result := parseStringList([]byte(data), "name")
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(result), result)
	}
}

func TestParseStringList_FlatArray(t *testing.T) {
	data := `[{"name": "a"}, {"name": "b"}]`
	result := parseStringList([]byte(data), "name")
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(result), result)
	}
}

func TestConfigTOMLAppend(t *testing.T) {
	// Simulate the append strategy for MCP creation
	original := `default_tool = "claude"

[claude]
  allow_dangerous_mode = true

# My important comment
[mcps]
  [mcps.github]
    args = ["-y", "@anthropic/mcp-github"]
    command = "npx"
    [mcps.github.env]
      GITHUB_TOKEN = "${GITHUB_TOKEN}"
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}

	// Append a new MCP
	fragment := `
  [mcps.exa]
    args = ["-y", "@anthropic/mcp-exa"]
    command = "npx"
    description = "Exa search MCP"
`

	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(fragment)
	f.Close()

	// Read back and verify
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	// Original content preserved
	if !strings.Contains(content, "# My important comment") {
		t.Error("append strategy lost comments")
	}
	if !strings.Contains(content, `allow_dangerous_mode = true`) {
		t.Error("append strategy lost original values")
	}

	// New MCP added
	if !strings.Contains(content, "[mcps.exa]") {
		t.Error("new MCP not appended")
	}

	// Parse the result to verify it's valid TOML
	_, err = tomlDecode(content)
	if err != nil {
		t.Fatalf("appended config is not valid TOML: %v", err)
	}
}

// tomlDecode is a test helper to validate TOML
func tomlDecode(s string) (map[string]interface{}, error) {
	// Use the same library as production code
	var m map[string]interface{}
	_, err := toml.Decode(s, &m)
	return m, err
}
