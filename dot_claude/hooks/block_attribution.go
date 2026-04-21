package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// HookInput is the JSON structure Claude Code sends to PreToolUse hooks via stdin.
type HookInput struct {
	HookEventName string `json:"hook_event_name"`
	ToolName      string `json:"tool_name"`
	ToolInput     struct {
		Command        string `json:"command"`
		Description    string `json:"description"`
		Timeout        int    `json:"timeout"`
		RunInBg        bool   `json:"run_in_background"`
	} `json:"tool_input"`
}

// HookOutput is the JSON structure Claude Code expects on stdout from PreToolUse hooks.
type HookOutput struct {
	HookSpecificOutput struct {
		HookEventName          string      `json:"hookEventName"`
		PermissionDecision     string      `json:"permissionDecision"`
		PermissionDecisionReason string    `json:"permissionDecisionReason"`
		UpdatedInput           interface{} `json:"updatedInput,omitempty"`
	} `json:"hookSpecificOutput"`
}

var attributionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)Co-Authored-By.*(?:Claude|anthropic)`),
	regexp.MustCompile(`(?i)Generated with.*Claude Code`),
	regexp.MustCompile(`🤖.*Claude Code`),
}

// containsAttribution checks if any attribution pattern matches the command string.
func containsAttribution(cmd string) (bool, string) {
	for _, p := range attributionPatterns {
		if loc := p.FindString(cmd); loc != "" {
			return true, loc
		}
	}
	return false, ""
}

// isGitCommitOrPR returns true if the command is a git commit or gh pr create.
func isGitCommitOrPR(cmd string) bool {
	return strings.Contains(cmd, "git commit") || strings.Contains(cmd, "gh pr create")
}

func deny(reason string) {
	out := HookOutput{}
	out.HookSpecificOutput.HookEventName = "PreToolUse"
	out.HookSpecificOutput.PermissionDecision = "deny"
	out.HookSpecificOutput.PermissionDecisionReason = reason
	json.NewEncoder(os.Stdout).Encode(out)
	os.Exit(0)
}

func allow() {
	// Empty stdout + exit 0 = allow without modification.
	os.Exit(0)
}

func main() {
	var input HookInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		fmt.Fprintf(os.Stderr, "block_attribution: failed to parse hook input: %v\n", err)
		allow() // fail open — don't block on parse errors
		return
	}

	cmd := input.ToolInput.Command
	if cmd == "" {
		allow()
		return
	}

	if !isGitCommitOrPR(cmd) {
		allow()
		return
	}

	if found, match := containsAttribution(cmd); found {
		deny(fmt.Sprintf(
			"Command contains Claude attribution (%q). Remove Co-Authored-By, "+
				"'Generated with Claude Code', or similar lines before committing.",
			match,
		))
		return
	}

	allow()
}
