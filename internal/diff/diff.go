package diff

import (
	"fmt"
	"sort"

	"github.com/afterthought/shuffle/internal/deck"
	"github.com/afterthought/shuffle/internal/state"
)

// ActionType describes what kind of change to make.
type ActionType string

const (
	CreateProfile   ActionType = "create_profile"
	CreateMCP       ActionType = "create_mcp"
	CreateGroup     ActionType = "create_group"
	CreateConductor ActionType = "create_conductor"
	MoveConductor   ActionType = "move_conductor"
	CreateSession   ActionType = "create_session"
	AttachMCP       ActionType = "attach_mcp"
	AttachSkill     ActionType = "attach_skill"
)

// Action is a single change to apply.
type Action struct {
	Type        ActionType
	Description string
	// Context fields for execution
	Name      string
	Group     string
	Session   *deck.Session
	Conductor *deck.Conductor
	MCP       *deck.MCP
	MCPName   string
	SkillName string
	Profile   *deck.Profile
	SessionID string // for attach operations
}

// Plan computes the diff between desired deck state and current state.
func Plan(d *deck.Deck, cur *state.Current) []Action {
	var actions []Action

	// 1. Profile
	if d.Profile != nil {
		if !contains(cur.Profiles, d.Profile.Name) {
			actions = append(actions, Action{
				Type:        CreateProfile,
				Description: fmt.Sprintf("create profile: %s", d.Profile.Name),
				Name:        d.Profile.Name,
				Profile:     d.Profile,
			})
		}
	}

	// 2. MCPs (sorted for deterministic output)
	for _, name := range sortedKeys(d.MCPs) {
		if !contains(cur.MCPs, name) {
			mcp := d.MCPs[name]
			actions = append(actions, Action{
				Type:        CreateMCP,
				Description: fmt.Sprintf("create MCP: %s", name),
				Name:        name,
				MCP:         &mcp,
			})
		}
	}

	// 3. Groups (sorted)
	for _, groupName := range sortedKeys(d.Groups) {
		if !contains(cur.Groups, groupName) {
			actions = append(actions, Action{
				Type:        CreateGroup,
				Description: fmt.Sprintf("create group: %s", groupName),
				Name:        groupName,
			})
		}
	}

	// 4. Conductors (sorted by group, then conductor name)
	existingConductors := make(map[string]bool)
	for _, c := range cur.Conductors {
		existingConductors[c.Name] = true
	}
	for _, groupName := range sortedKeys(d.Groups) {
		group := d.Groups[groupName]
		for _, condName := range sortedKeys(group.Conductors) {
			if !existingConductors[condName] {
				cond := group.Conductors[condName]
				actions = append(actions, Action{
					Type:        CreateConductor,
					Description: fmt.Sprintf("create conductor: %s", condName),
					Name:        condName,
					Conductor:   &cond,
					Group:       groupName,
				})
				actions = append(actions, Action{
					Type:        MoveConductor,
					Description: fmt.Sprintf("move conductor %s to group: %s", condName, groupName),
					Name:        condName,
					Group:       groupName,
				})
			}
		}
	}

	// 5. Sessions (sorted by group, then session name)
	existingSessions := make(map[string]bool) // "group/title" -> true
	sessionIDs := make(map[string]string)     // "group/title" -> id
	for _, s := range cur.Sessions {
		key := s.Group + "/" + s.Title
		existingSessions[key] = true
		sessionIDs[key] = s.ID
	}
	for _, groupName := range sortedKeys(d.Groups) {
		group := d.Groups[groupName]
		for _, sessName := range sortedKeys(group.Sessions) {
			sess := group.Sessions[sessName]
			key := groupName + "/" + sessName
			if !existingSessions[key] {
				actions = append(actions, Action{
					Type:        CreateSession,
					Description: fmt.Sprintf("create session: %s in group %s", sessName, groupName),
					Name:        sessName,
					Group:       groupName,
					Session:     &sess,
				})
			}
		}
	}

	// 6. MCP and skill attachments (for existing sessions only — new sessions get MCPs via launch --mcp)
	for _, groupName := range sortedKeys(d.Groups) {
		group := d.Groups[groupName]
		for _, sessName := range sortedKeys(group.Sessions) {
			sess := group.Sessions[sessName]
			key := groupName + "/" + sessName
			id, exists := sessionIDs[key]
			if !exists {
				continue
			}
			for _, mcpName := range sess.MCPs {
				actions = append(actions, Action{
					Type:        AttachMCP,
					Description: fmt.Sprintf("attach MCP %s to session %s", mcpName, sessName),
					MCPName:     mcpName,
					SessionID:   id,
					Name:        sessName,
				})
			}
			for _, skillName := range sess.Skills {
				actions = append(actions, Action{
					Type:        AttachSkill,
					Description: fmt.Sprintf("attach skill %s to session %s", skillName, sessName),
					SkillName:   skillName,
					SessionID:   id,
					Name:        sessName,
				})
			}
		}
	}

	return actions
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

// sortedKeys returns the keys of a map sorted alphabetically.
// Works with any map type via type parameter (Go 1.18+).
func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
