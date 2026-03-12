package deck

// Deck represents the top-level .deck.yaml structure.
type Deck struct {
	Name    string              `yaml:"name"`
	Profile *Profile            `yaml:"profile,omitempty"`
	MCPs    map[string]MCP      `yaml:"mcps,omitempty"`
	Shells  map[string]Shell    `yaml:"shells,omitempty"`
	Groups  map[string]Group    `yaml:"groups,omitempty"`
}

// Profile declares an agent-deck profile.
type Profile struct {
	Name string `yaml:"name"`
}

// MCP declares an MCP server definition for config.toml.
type MCP struct {
	Command     string            `yaml:"command"`
	Args        []string          `yaml:"args,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	Description string            `yaml:"description,omitempty"`
}

// Shell is a reusable session template.
type Shell struct {
	Agent    string   `yaml:"agent,omitempty"`
	Command  string   `yaml:"command,omitempty"`
	Worktree string   `yaml:"worktree,omitempty"`
	MCPs     []string `yaml:"mcps,omitempty"`
	Skills   []string `yaml:"skills,omitempty"`
	Prompt   string   `yaml:"prompt,omitempty"`
}

// Group declares a session group with optional sessions and conductors.
type Group struct {
	Sessions   map[string]Session   `yaml:"sessions,omitempty"`
	Conductors map[string]Conductor `yaml:"conductors,omitempty"`
}

// Session declares an agent-deck session.
type Session struct {
	Agent    string   `yaml:"agent,omitempty"`
	Command  string   `yaml:"command,omitempty"`
	Shell    string   `yaml:"shell,omitempty"`
	Worktree string   `yaml:"worktree,omitempty"`
	MCPs     []string `yaml:"mcps,omitempty"`
	Skills   []string `yaml:"skills,omitempty"`
	Prompt   string   `yaml:"prompt,omitempty"`
	Path     string   `yaml:"path,omitempty"`
	Parent   string   `yaml:"parent,omitempty"`
}

// Conductor declares a conductor meta-agent.
type Conductor struct {
	Description string `yaml:"description,omitempty"`
	Heartbeat   *bool  `yaml:"heartbeat,omitempty"`
	ClaudeMD    string `yaml:"claude_md,omitempty"`
	PolicyMD    string `yaml:"policy_md,omitempty"`
}
