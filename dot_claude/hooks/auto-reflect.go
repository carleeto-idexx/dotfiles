package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type state struct {
	Counter       int    `json:"counter"`
	Frequency     int    `json:"frequency"`
	LastReflected string `json:"last_reflected"`
}

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		os.Exit(0)
	}

	statePath := filepath.Join(home, ".claude", "state", "auto-reflect.json")

	s := state{Counter: 0, Frequency: 10}

	data, err := os.ReadFile(statePath)
	if err == nil {
		_ = json.Unmarshal(data, &s)
	}

	if s.Frequency < 1 {
		s.Frequency = 10
	}

	s.Counter++

	if s.Counter < s.Frequency {
		writeState(statePath, s)
		os.Exit(0)
	}

	// Reset counter and trigger reflection
	s.Counter = 0
	s.LastReflected = time.Now().UTC().Format(time.RFC3339)
	writeState(statePath, s)

	fmt.Println("[Auto-Skill Reflection] Review this session for skill opportunities.")
	fmt.Println("Check ~/.claude/skills/auto-manifest.md for existing auto skills.")
	fmt.Println("Create, update, or suggest pruning as needed per the Auto Skill System instructions in ~/.claude/CLAUDE.md.")
	fmt.Println("Run reflection in the background using a background Agent.")
}

func writeState(path string, s state) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, data, 0o644)
}
