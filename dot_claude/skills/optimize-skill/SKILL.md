---
name: optimize-skill
description: Validates and optimizes Claude Code skills against best practices. Use after editing a SKILL.md file, or when the user says "optimize skill", "check skill", or "lint skill".
allowed-tools: Read, Edit, Glob, Grep, Bash
argument-hint: "[skill-name or path]"
---

# Skill Optimizer

Validate and optimize a Claude Code skill against best practices. Apply fixes automatically unless the change is risky (e.g., changing invocation behavior, removing content). Report all optimizations with rationale and impact.

`$ARGUMENTS` is the skill name (e.g., `review-pr`) or path. If empty, ask the user which skill to optimize.

## Locate the Skill

1. If `$ARGUMENTS` is a path, use it directly.
2. If it's a skill name, look in:
   - `.claude/skills/<name>/SKILL.md` (project-level)
   - `~/.claude/skills/<name>/SKILL.md` (user-level)
3. If not found, list available skills and ask.

Also read any supporting files (referenced `.md` files in the same directory).

## Checks

Run all checks below. For each issue found, classify as:
- **Auto-fix**: Safe to apply without asking (formatting, missing fields with obvious values)
- **Confirm**: Behavior-changing — describe the fix and ask before applying

### Frontmatter

| Check | Auto-fix? | Detail |
|-------|-----------|--------|
| `name` exists and matches directory name | Auto | Add or correct if mismatched |
| `name` uses only lowercase, numbers, hyphens; max 64 chars | Auto | Fix invalid characters |
| `description` exists, 20-300 chars, specific (not vague) | Confirm | Flag vague descriptions like "helps with", "useful for" |
| `disable-model-invocation` set for skills with side effects | Confirm | Flag if skill runs git push, deploy, send messages, etc. without this |
| `disable-model-invocation: true` + `user-invocable: false` not both set | Auto | Remove contradictory field |
| `argument-hint` exists if `$ARGUMENTS` is used in content | Auto | Add hint based on usage context |
| `allowed-tools` lists only valid tool names | Flag | Warn on unrecognized tool names |
| `context: fork` has matching `agent` field and task content | Flag | Warn if guidelines-only with fork context |

### Prompt Structure

| Check | Auto-fix? | Detail |
|-------|-----------|--------|
| Uses `$ARGUMENTS` if skill accepts args (not relying on append) | Auto | Replace `$0` or missing reference with `$ARGUMENTS` |
| SKILL.md line count | Flag | Warn if >500 lines; suggest moving reference content to separate files |
| No hardcoded absolute paths | Flag | Warn on `/Users/`, `/home/`, `C:\` etc. |
| Supporting files referenced correctly | Flag | Verify `[text](file.md)` targets exist in skill directory |
| Example output in templates uses consistent format | Flag | Check that markdown templates are well-formed |

### Content Quality

| Check | Auto-fix? | Detail |
|-------|-----------|--------|
| Redundant instructions | Auto | Remove duplicated paragraphs or near-identical rules |
| Conflicting instructions | Confirm | Flag rules that contradict each other |
| Dead sections | Flag | Sections referenced but not populated |
| Overly prescriptive formatting | Flag | Warn if >30% of content is template formatting vs actual guidance |

### Context Window Efficiency

The entire SKILL.md loads into context on every invocation. Wasted tokens here are wasted on every single use.

| Check | Auto-fix? | Detail |
|-------|-----------|--------|
| Line count >500 | Flag | Suggest extracting reference material to separate files linked from SKILL.md |
| Large example blocks (>30 lines) | Confirm | Move to a `examples.md` or `reference.md` and link |
| Repeated instructions | Auto | Consolidate rules that say the same thing in different words |
| Verbose phrasing | Auto | Tighten wordy instructions (e.g., "You should make sure to always" → "Always") |
| Output templates >20% of total content | Flag | Templates are useful but shouldn't dominate; consider shortening |
| Supporting files that could be loaded on-demand | Flag | If a linked `.md` file is >200 lines and only needed for one phase, note it |

When auto-fixing verbosity, preserve meaning exactly — only remove filler words and redundant qualifiers. Never change the author's intent or remove nuance.

### Tool Scoping

| Check | Auto-fix? | Detail |
|-------|-----------|--------|
| `allowed-tools` matches actual tool usage in instructions | Flag | Warn if skill instructs use of a tool not in the allowlist |
| Overly broad Bash permissions | Flag | Warn on `Bash(*)` without justification |
| Missing tools for described workflow | Flag | If skill says "edit the file" but Edit isn't in allowed-tools |

## Output Format

After running all checks, present results as:

```markdown
## Skill Optimization: <skill-name>

### Applied (auto-fixed)
- **[check name]**: [what was changed] — [why] — [impact: e.g., "skill now auto-invocable"]

### Needs Confirmation
- **[check name]**: [what should change] — [why] — [risk if applied]

### Warnings
- **[check name]**: [observation] — [suggestion]

### Passed
[count] checks passed with no issues.
```

If everything passes, just say so — don't pad with a full report.

## Constraints

- Never remove content without confirmation
- Never change `disable-model-invocation` without confirmation (affects when skill triggers)
- Never modify supporting files without reading them first
- Preserve the author's voice and intent — optimize structure, not style
