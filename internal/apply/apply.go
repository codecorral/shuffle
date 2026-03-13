package apply

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/afterthought/shuffle/internal/deck"
	adiff "github.com/afterthought/shuffle/internal/diff"
	"github.com/afterthought/shuffle/internal/state"
)

// Execute runs a list of actions against agent-deck.
// profile is the agent-deck profile to use (empty for default).
func Execute(actions []adiff.Action, profile string) error {
	var errs []error
	for _, a := range actions {
		fmt.Printf("  %s\n", a.Description)
		var err error
		switch a.Type {
		case adiff.CreateProfile:
			err = createProfile(a.Name)
		case adiff.CreateMCP:
			err = createMCP(a.Name, a.MCP)
		case adiff.CreateGroup:
			err = createGroup(profile, a.Name)
		case adiff.CreateConductor:
			err = createConductor(profile, a.Name, a.Conductor)
		case adiff.MoveConductor:
			err = moveConductor(profile, a.Name, a.Group)
		case adiff.CreateSession:
			err = createSession(profile, a.Name, a.Group, a.Session)
		case adiff.AttachMCP:
			err = attachMCP(profile, a.SessionID, a.MCPName)
		case adiff.AttachSkill:
			err = attachSkill(profile, a.SessionID, a.SkillName)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ERROR: %s: %v\n", a.Description, err)
			errs = append(errs, fmt.Errorf("%s: %w", a.Description, err))
		}
	}
	if len(errs) > 0 {
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = e.Error()
		}
		return fmt.Errorf("%d action(s) failed:\n  %s", len(errs), strings.Join(msgs, "\n  "))
	}
	return nil
}

func createProfile(name string) error {
	return runAgentDeck("profile", "create", name)
}

func createMCP(name string, mcp *deck.MCP) error {
	configPath, err := state.ConfigPath()
	if err != nil {
		return err
	}

	// Read existing config to check if MCP already exists
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config.toml: %w", err)
	}

	var config map[string]interface{}
	if err := toml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parsing config.toml: %w", err)
	}

	// Check if already exists
	if mcps, ok := config["mcps"]; ok {
		if mcpMap, ok := mcps.(map[string]interface{}); ok {
			if _, exists := mcpMap[name]; exists {
				return nil
			}
		}
	}

	// Append new MCP entry to config.toml instead of re-encoding the entire file.
	// This preserves comments, formatting, and ordering of existing content.
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Build the TOML fragment for this MCP
	var buf strings.Builder
	fmt.Fprintf(&buf, "\n  [mcps.%s]\n", name)
	if len(mcp.Args) > 0 {
		// Format args as TOML array
		quotedArgs := make([]string, len(mcp.Args))
		for i, a := range mcp.Args {
			quotedArgs[i] = fmt.Sprintf("%q", a)
		}
		fmt.Fprintf(&buf, "    args = [%s]\n", strings.Join(quotedArgs, ", "))
	}
	fmt.Fprintf(&buf, "    command = %q\n", mcp.Command)
	if mcp.Description != "" {
		fmt.Fprintf(&buf, "    description = %q\n", mcp.Description)
	}
	if len(mcp.Env) > 0 {
		fmt.Fprintf(&buf, "    [mcps.%s.env]\n", name)
		for k, v := range mcp.Env {
			fmt.Fprintf(&buf, "      %s = %q\n", k, v)
		}
	}

	_, err = f.WriteString(buf.String())
	return err
}

func createGroup(profile, name string) error {
	args := buildProfileArgs(profile, "group", "create", name)
	return runAgentDeck(args...)
}

func createConductor(profile, name string, cond *deck.Conductor) error {
	args := buildProfileArgs(profile, "conductor", "setup", name)

	if cond.Description != "" {
		args = append(args, "-description", cond.Description)
	}

	if cond.Heartbeat != nil && !*cond.Heartbeat {
		args = append(args, "-no-heartbeat")
	}

	if cond.ClaudeMD != "" {
		tmpPath, err := writeTemp(cond.ClaudeMD)
		if err != nil {
			return err
		}
		defer os.Remove(tmpPath)
		args = append(args, "-claude-md", tmpPath)
	}

	if cond.PolicyMD != "" {
		tmpPath, err := writeTemp(cond.PolicyMD)
		if err != nil {
			return err
		}
		defer os.Remove(tmpPath)
		args = append(args, "-policy-md", tmpPath)
	}

	return runAgentDeck(args...)
}

func moveConductor(profile, name, group string) error {
	// Conductor session ID is typically the conductor name
	args := buildProfileArgs(profile, "group", "move", name, group)
	return runAgentDeck(args...)
}

func createSession(profile, name, group string, sess *deck.Session) error {
	path := sess.Path
	if path == "" {
		path = "."
	}

	// Determine the command (-c flag)
	cmd := sess.Agent
	if sess.Command != "" {
		cmd = sess.Command
	}

	args := buildProfileArgs(profile, "launch", path)
	if cmd != "" {
		args = append(args, "-c", cmd)
	}
	args = append(args, "-t", name)
	args = append(args, "-g", group)

	if sess.Parent != "" {
		args = append(args, "-parent", sess.Parent)
	}

	if sess.Worktree != "" {
		// Worktree field specifies the location strategy (subdirectory, sibling, or custom path).
		// agent-deck requires -w <branch> to create a worktree, plus --location for strategy.
		// Use the session name as branch name with --new-branch to auto-create it.
		args = append(args, "-w", name, "--new-branch", "--location", sess.Worktree)
	}

	for _, mcp := range sess.MCPs {
		args = append(args, "--mcp", mcp)
	}

	if sess.Prompt != "" {
		args = append(args, "-m", sess.Prompt)
	}

	return runAgentDeck(args...)
}

func attachMCP(profile, sessionID, mcpName string) error {
	// Check if already attached
	attached, err := state.AttachedMCPs(profile, sessionID)
	if err == nil {
		for _, m := range attached {
			if m == mcpName {
				return nil // already attached
			}
		}
	}

	args := buildProfileArgs(profile, "mcp", "attach", sessionID, mcpName)
	return runAgentDeck(args...)
}

func attachSkill(profile, sessionID, skillName string) error {
	// Check if already attached
	attached, err := state.AttachedSkills(profile, sessionID)
	if err == nil {
		for _, s := range attached {
			if s == skillName {
				return nil // already attached
			}
		}
	}

	args := buildProfileArgs(profile, "skill", "attach", sessionID, skillName)
	return runAgentDeck(args...)
}

// writeTemp writes content to a temp file and returns its path.
// Caller must clean up the file.
func writeTemp(content string) (string, error) {
	tmp, err := os.CreateTemp("", "shuffle-*.md")
	if err != nil {
		return "", err
	}
	if _, err := tmp.WriteString(content); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", err
	}
	tmp.Close()
	return tmp.Name(), nil
}

func buildProfileArgs(profile string, args ...string) []string {
	if profile != "" {
		return append([]string{"-p", profile}, args...)
	}
	return args
}

func runAgentDeck(args ...string) error {
	cmd := exec.Command("agent-deck", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("agent-deck %s: %w", strings.Join(args, " "), err)
	}
	return nil
}
