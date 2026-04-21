package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type StopInput struct {
	StopHookActive       bool   `json:"stop_hook_active"`
	LastAssistantMessage string `json:"last_assistant_message"`
}

type StopOutput struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

func main() {
	var input StopInput
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		os.Exit(0)
	}

	// Prevent infinite loops: if we're already re-running after a stop hook, allow.
	if input.StopHookActive {
		os.Exit(0)
	}

	if strings.Contains(input.LastAssistantMessage, "\u2014") || strings.Contains(input.LastAssistantMessage, "--") {
		out := StopOutput{
			Decision: "block",
			Reason:   "Your response contains emdashes (\u2014) or double dashes (--). Rewrite replacing them with hyphens, commas, or other punctuation.",
		}
		json.NewEncoder(os.Stdout).Encode(out)
		fmt.Fprintf(os.Stderr, "Blocked: response contains emdashes or double dashes.\n")
		os.Exit(2)
	}

	os.Exit(0)
}
