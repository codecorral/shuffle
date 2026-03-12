package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/afterthought/shuffle/internal/apply"
	"github.com/afterthought/shuffle/internal/deck"
	"github.com/afterthought/shuffle/internal/diff"
	"github.com/afterthought/shuffle/internal/state"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "validate":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: shuffle validate <deck.yaml>")
			os.Exit(1)
		}
		os.Exit(runValidate(os.Args[2]))
	case "diff":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: shuffle diff <deck.yaml>")
			os.Exit(1)
		}
		os.Exit(runDiff(os.Args[2]))
	case "deal":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: shuffle deal <deck.yaml>")
			os.Exit(1)
		}
		os.Exit(runDeal(os.Args[2]))
	case "help", "--help", "-h":
		printUsage()
	case "version", "--version":
		fmt.Println("shuffle v0.1.0")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`shuffle - Declarative configuration for agent-deck

Usage: shuffle <command> <deck.yaml>

Commands:
  validate <deck.yaml>   Validate a deck file without applying
  diff <deck.yaml>       Show what deal would change (dry run)
  deal <deck.yaml>       Apply a deck file to agent-deck
  version                Show version
  help                   Show this help

Examples:
  shuffle validate codecorral.deck.yaml
  shuffle diff codecorral.deck.yaml
  shuffle deal codecorral.deck.yaml`)
}

func runValidate(path string) int {
	d, err := deck.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	errs := deck.Validate(d)
	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr, "Validation errors:")
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  - %v\n", e)
		}
		return 1
	}

	fmt.Println("Deck is valid.")
	return 0
}

func runDiff(path string) int {
	d, err := loadAndResolve(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	profile, profileSource := resolveProfile(d)
	fmt.Printf("Profile: %s (%s)\n\n", profile, profileSource)

	cur, err := state.Discover(profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering state: %v\n", err)
		return 1
	}

	actions := diff.Plan(d, cur)
	if len(actions) == 0 {
		fmt.Println("No changes needed. Everything is up to date.")
		return 0
	}

	fmt.Printf("Planned changes (%d):\n", len(actions))
	for _, a := range actions {
		fmt.Printf("  + %s\n", a.Description)
	}
	return 0
}

func runDeal(path string) int {
	// Check agent-deck is available
	if _, err := exec.LookPath("agent-deck"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: agent-deck not found on PATH. Install it first: https://github.com/asheshgoplani/agent-deck")
		return 1
	}

	d, err := loadAndResolve(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	profile, profileSource := resolveProfile(d)
	fmt.Printf("Profile: %s (%s)\n\n", profile, profileSource)

	cur, err := state.Discover(profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering state: %v\n", err)
		return 1
	}

	actions := diff.Plan(d, cur)
	if len(actions) == 0 {
		fmt.Println("No changes needed. Everything is up to date.")
		return 0
	}

	basePath, _ := filepath.Abs(filepath.Dir(path))
	fmt.Printf("Applying %d changes...\n", len(actions))
	if err := apply.Execute(actions, profile, basePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Println("Deal complete.")
	return 0
}

// resolveProfile returns the profile name and a description of where it came from.
func resolveProfile(d *deck.Deck) (string, string) {
	if d.Profile != nil && d.Profile.Name != "" {
		return d.Profile.Name, "from deck"
	}
	name, err := state.DefaultProfile()
	if err != nil {
		return "default", "agent-deck default"
	}
	return name, "agent-deck default"
}

func loadAndResolve(path string) (*deck.Deck, error) {
	d, err := deck.Load(path)
	if err != nil {
		return nil, err
	}

	errs := deck.Validate(d)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  - %v\n", e)
		}
		return nil, fmt.Errorf("validation failed with %d errors", len(errs))
	}

	if err := deck.ResolveShells(d); err != nil {
		return nil, fmt.Errorf("resolving shells: %w", err)
	}

	return d, nil
}
