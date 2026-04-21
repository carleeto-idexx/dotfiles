---
name: review-pr
description: Reviews pull requests in phases (design, correctness, tests, style) with verified findings. Supports GitHub and GitLab, stacked PRs, local branches, comment drafting, and comment verification. Use when the user says "review PR", "review this PR", or invokes /review-pr with a PR number, URL, or branch name (prefixed with "branch:").
disable-model-invocation: false
argument-hint: "[PR number, URL, or branch:<name>]"
allowed-tools: Read, Edit, Write, Grep, Glob, Bash, Agent
---

# Pull Request Review

## PACING RULE (read this first)

**Your task scope is exactly ONE finding per response.** Your entire job is to read one finding file from disk and present it. There is nothing else to generate. You are done.

The user will ask for the next one. At that point you read the next file and present it.

Think of it like a function: you are called with an index, you read one file, you return its contents formatted. You never load all findings into context.

**Why this matters:** The user steers the review interactively. This rule was violated repeatedly when framed as "stop after one" because generating-then-stopping fights generation momentum. This design is different: findings live on disk, you only read one at a time, so there is only one thing to generate.

**How findings get to disk:** The analysis Agent writes each finding as a separate markdown file in `.reviews/pr-<ID>/findings/`, sorted by severity. It also writes a `summary.json` with counts. You (the main thread) never see the raw analysis. You read files.

---

`$ARGUMENTS` contains the PR number, URL, or local branch reference passed to `/review-pr`.

**Local branch mode?** If `$ARGUMENTS` starts with `branch:`, see [Local Branch Mode](#local-branch-mode) below.

**Multiple MR numbers?** See [Batch Mode](reference.md#batch-mode). For a single MR, continue below.

## Execution Order

1. Setup (fetch PR, create worktree, get diff, get comments)
2. Comment Deduplication (catalog existing reviewer findings)
3. Triage (classify Tier 1/2/3 by diff size)
4. Draft MR Detection (if draft, switch to questions-only mode per [reference.md](reference.md#draft-mr-detection) and stop)
5. MR Identity Banner
6. Stacked PR Detection
7. MR Timeline (then **pause** for user acknowledgment)
8. **Analysis** (run as Agent subagent):
   - Launch an Agent to perform Phases 1-4 in the worktree
   - The agent writes each finding as a separate file to `.reviews/pr-<ID>/findings/` (see [Findings on Disk](#findings-on-disk))
   - The agent writes `summary.json` with severity counts and recommendation
   - **The agent's return message to you must be exactly:** "Findings written to .reviews/pr-<ID>/findings/" (no summaries, no finding details, no recommendations in the return value). This prevents findings from leaking into the main thread's context.
9. **Triage report**: Read `summary.json` and tell the user how many findings at each severity and the recommendation. If zero findings, show the recommendation directly. **End your response here.** Wait for the user to say "go" or pick a severity.
10. **Output**: Read and present the first (most severe) finding file. That is your entire task for this response. Update `progress.json`.
11. When the user signals, read and present the next finding file. Update `progress.json`. That is your entire task.
12. After the last finding has been presented and acknowledged, show the **Summary**.

---

## Issue Format {#issue-format}

All findings use this format:
- **Location**: `file:line` (path relative to repo root)
- **Confidence**: 0-100 (below 50 = auto-dismiss, do not present)
- **Impact**: Why it matters (lead with this)
- **Problem**: What's wrong, detailed explanation
- **Verified**: How you confirmed (trace callers, check invariants, read tests)

**Before raising any issue:**
- Prove the condition is reachable. If no code path triggers it, don't raise it.
- Steel-man first: check for comments, surrounding context, framework behavior, or legacy constraints.
- Check the codebase for answers before speculating.
- For correctness issues, illustrate with a numbered step-by-step sequence.
- **CRITICAL: Always Read the actual source file in the worktree to confirm line numbers.** Never cite a line number you haven't verified by reading the file.
- Use paths relative to the repo root. Never use absolute or worktree-relative paths.

---

## Local Branch Mode

When `$ARGUMENTS` starts with `branch:`, review a local branch instead of a remote MR.

**Syntax:** `branch:<branch-name>` with optional `base:<base-branch>` (defaults to `main`).

Examples:
- `/review-pr branch:feat/ezyvet-validate-sap`
- `/review-pr branch:feat/ezyvet-validate-sap base:develop`

### Local Branch Setup

1. **Parse arguments**: Extract `<BRANCH>` from `branch:<BRANCH>`. Extract `<BASE>` from `base:<BASE>` if present, otherwise default to `main`.

2. **Validate branch exists**: `git rev-parse --verify <BRANCH>` (must succeed).

3. **Clean up stale worktrees**:
   ```bash
   for dir in .reviews/branch-*; do [ "$dir" != ".reviews/branch-<BRANCH>" ] && git worktree remove "$dir" --force 2>/dev/null; done
   ```

4. **Create worktree from local branch**:
   - New: `git worktree add .reviews/branch-<BRANCH> <BRANCH>`
   - Existing: `git -C .reviews/branch-<BRANCH> reset --hard <BRANCH>`

5. **Get diff**: `git -C .reviews/branch-<BRANCH> diff <BASE>...HEAD`

### Local Branch Execution Order

Skip steps that require a remote MR (comment deduplication, draft detection, stacked PR detection, MR timeline, comment drafting, comment verification). Run:

1. Setup (above)
2. Triage (classify Tier 1/2/3 by diff size)
3. Identity Banner (use `> **Reviewing branch:** <BRANCH> (against <BASE>)` / `> Tier [N].`)
4. **Analysis**: Launch Agent subagent to run Phases 1-4 and write findings to `.reviews/branch-<BRANCH>/findings/`. Agent returns only "Findings written to [path]".
5. **Triage report**: Read `summary.json`, tell user finding counts by severity. End response. Wait for user.
6. **Output**: Read and present one finding file at a time (see Pacing Rule). Update `progress.json` after each.
7. Summary

After review: `git worktree remove .reviews/branch-<BRANCH> --force`

Work from the worktree. Focus ONLY on added/modified lines relative to `<BASE>`.

Then continue to the shared [Issue Format](#issue-format), [Triage](#triage), and [Phase 1-4](#phase-1-design--architecture) sections below. All phase logic, issue format, and output rules apply identically.

---

## Setup

0. **Get PR identifier**: If `$ARGUMENTS` is empty, ask the user.

1. **Clean up stale worktrees**:
   ```bash
   for dir in .reviews/pr-*; do [ "$dir" != ".reviews/pr-<ID>" ] && git worktree remove "$dir" --force; done
   for branch in $(git branch --list 'review-*'); do [ "$branch" != "review-<ID>" ] && git branch -D "$branch"; done
   ```

2. **Verify tools**: Confirm `gh` (GitHub) or `glab` (GitLab) CLI is available.

3. **Detect platform**: `git remote get-url origin` contains `github` or `gitlab`

4. **Fetch all branches** (always, even when resuming): `git fetch origin`
   This ensures the diff base is current. A stale `origin/main` produces false positives.

5. **Fetch PR head**:
   - GitHub: `git fetch origin pull/<ID>/head:review-<ID> --force`
   - GitLab: `git fetch origin merge-requests/<ID>/head:review-<ID> --force`

6. **Worktree** (always reset to fetched head):
   - New: `git worktree add .reviews/pr-<ID> review-<ID>`
   - Existing: `git -C .reviews/pr-<ID> reset --hard review-<ID>`

7. **Get target branch**:
   - GitHub: `gh pr view <ID> --json baseRefName -q .baseRefName`
   - GitLab: `glab mr view <ID> --output json | jq -r .target_branch`

8. **Get diff**: `git -C .reviews/pr-<ID> diff origin/<target_branch>...HEAD`

9. **Get comments** (both top-level and inline/resolved):
   - GitHub: `gh pr view <ID> --comments` + `gh api repos/<owner>/<repo>/pulls/<ID>/comments --paginate`
   - GitLab: `glab mr view <ID> --comments` + `glab api projects/<id>/merge_requests/<ID>/discussions`

10. **Get MR description and linked context**:
    - Read the MR description (fetched in step 9 or via `glab mr view <ID> --output json | jq -r .description`)
    - Extract any linked JIRA/issue keys (e.g., PT-4355, INGEST-123) from the title, description, or branch name
    - If the `atlassian:search-company-knowledge` skill is available (or Atlassian MCP tools like `getJiraIssue`), fetch the JIRA issue details including:
      - Title and status
      - Description
      - **Acceptance criteria** (often in an "Acceptance:" or "Scope:" block)
      - Linked Confluence pages for wider context
    - Pass all of the above to the analysis Agent. The acceptance criteria drive Phase 0 (scope alignment). The business context helps distinguish intentionally coupled changes from accidentally bundled ones in Phase 1 and beyond.

Work from the worktree. Focus ONLY on added/modified lines.

---

## Comment Deduplication

Catalog every finding another reviewer or bot already raised:
1. **Correct**: Don't present it. Note in "Existing Comments" as "Agree, already raised by [reviewer]."
2. **Wrong or incomplete**: Present it, but lead with the prior discussion and explain what it missed.
3. **Addressed by subsequent commit**: Note as resolved in timeline. Don't re-raise.

Zero redundant noise. The user should never think "the bot already said this."

---

## Triage

Classify the MR to calibrate review depth:
- **Tier 1** (light): <200 lines, <=3 files. Run Phase 0 + Phases 2-3.
- **Tier 2** (standard): 200-800 lines, 4-15 files. All phases.
- **Tier 3** (deep): >800 lines or >15 files. All phases + sub-agents for every user-visible finding.

**Phase 0 always runs** regardless of tier whenever a linked JIRA ticket is available. Scope misalignment is more important than style or coverage gaps.

Print tier in the identity banner. User can override with "go deeper" or "keep it light."

---

## MR Identity Banner

```markdown
> **Reviewing: <prefix><ID>** / <MR title>
> Tier [N]. Stack: <prefix>X, **<prefix><ID>**, <prefix>Z (if stacked)
```

Use `#` for GitHub, `!` for GitLab. Prefix all major headings with the MR number.

---

## Stacked PR Detection

**Always** check both directions:
- **Upstream**: If base is a feature branch, follow the chain back to the default branch.
- **Downstream**: Check for open PRs targeting this PR's source branch.

Present the full stack chain. Don't flag missing functionality a later PR adds. Do flag issues that compound as the stack grows.

---

## MR Timeline

Present a chronological narrative: commit sequence, reviewer comments (who, what, how addressed, resolved or open), current state. Keep it concise.

**Pause here** for user acknowledgment before proceeding to phases.

---

## Findings on Disk {#findings-on-disk}

The analysis Agent writes findings to `.reviews/pr-<ID>/findings/`. This directory persists across conversations, preserving review context over the PR's lifetime.

### File naming

Files are numbered by severity order: `01-<severity>-<phase>.md`, `02-<severity>-<phase>.md`, etc.

Severity levels (in order): `critical`, `high`, `medium`, `low`, `info`.

Example:
```
.reviews/pr-432/findings/
  summary.json
  01-high-behavior.md
  02-medium-tests.md
  03-low-style.md
```

### Finding file format

```markdown
---
severity: high
phase: behavior
confidence: 82
location: src/internal/server/entity.go:349
---

- **Impact**: [why it matters]
- **Problem**: [detailed description]
- **Verified**: [how confirmed]
```

### summary.json format

```json
{
  "mr": "!432",
  "title": "MR title here",
  "head_sha": "1d1f3071",
  "tier": 2,
  "total": 3,
  "by_severity": {"high": 1, "medium": 1, "low": 1},
  "recommendation": "Comment",
  "recommendation_reason": "one sentence",
  "existing_comments": [
    {"author": "reviewer", "status": "agree", "summary": "..."}
  ]
}
```

### Cursor tracking

After presenting each finding, update `progress.json` in the findings directory:

```json
{"last_presented": "01-high-behavior.md"}
```

To determine the next finding: list files matching `??-*.md`, sort lexically, pick the first one after `last_presented`. If `progress.json` doesn't exist, start from the first file.

### Resuming a review

If `.reviews/pr-<ID>/findings/` already exists when the skill is invoked:
1. Check if the PR HEAD has changed since the findings were written (compare `head_sha` in `summary.json`)
2. If unchanged: skip analysis, read `progress.json` to resume from where you left off
3. If changed: re-run analysis, overwrite findings, delete `progress.json`

### Safety rules

- **One Read per response**: You must call Read on at most ONE file in the findings directory per response. If you want to read another finding file, stop. The user will ask.
- **No inline fallback**: If the findings directory is empty or `summary.json` is missing after the Agent completes, tell the user the analysis failed. Do NOT attempt inline analysis as a fallback.

---

## Phase Definitions (for the analysis Agent)

The following phases define what the analysis Agent should look for. The Agent runs all applicable phases and returns structured findings. You (the main thread) then present them one at a time.

### Phase 0: Scope Alignment

**Focus**: Does the PR move in the direction of what the linked JIRA ticket (or equivalent spec) asks for, without contradicting it?

Run whenever the setup step identified a linked ticket with acceptance criteria. Skip only if no ticket was found.

**Core principle: one issue can be delivered across many PRs.** A PR that implements half of the acceptance criteria is fine if the other half is intentionally deferred to a sibling PR. Do NOT flag "missing" criteria as a defect. The question is whether the PR moves the issue forward, not whether it closes the issue.

Steps:
1. Enumerate each acceptance criterion from the linked ticket.
2. For each criterion, check whether the PR's changes move toward it, leave it untouched, or move away from it.
3. Classify each criterion as:
   - **Delivered**: code/tests implement this criterion in this PR. Note in coverage.
   - **Progressed**: PR moves toward this criterion without fully delivering it (e.g., lays foundation, implements part). Note in coverage, no finding.
   - **Untouched**: PR does not address this criterion. Assume it's deferred to another PR. Note in coverage, no finding.
   - **Contradicted**: PR actively moves *away* from this criterion or undoes prior progress. Raise a `scope` finding at `high` severity.
4. Flag changes in the PR that are clearly OUTSIDE the issue's direction (e.g., unrelated refactors, features from a different ticket). Raise a `scope` finding at `low` or `medium` severity depending on blast radius. Small drive-by cleanups in touched files are usually fine; entire new features are not.

Ask: *Does this PR move the ticket forward, or does any part of it move the ticket backward or sideways?*

**Ignore**: style, bugs (those belong to Phases 2-4). Also ignore the fact that not every AC is delivered in a single PR; multi-PR delivery is normal.

---

### Phase 1: Design & Architecture

**Focus**: Coupling, separation of concerns, pattern violations, over-engineering.
**Ignore**: Naming, formatting, bugs (unless architectural).
Ask: *Is this the simplest solution that works? Could a junior understand it?*

---

### Phase 2: Behavior & Correctness

**Focus**: Bugs, race conditions, error handling, logic gaps.
**Breaking changes** (schema, API, persistence): start with **BREAKING CHANGE DETECTED**

**Read the MR description first.** Before flagging changes as "bundled" or "unrelated," check whether the description explains why they belong together. Also check linked JIRA issues and Confluence pages for business context. Changes that look independent in the diff may be intentionally coupled when you understand the feature.

**Test gap analysis**: For every behavioral issue, investigate why tests don't catch it. Include the analysis in the finding.

**Deep verification**: For user-visible, disputed, or cross-boundary findings, launch an Agent (subagent_type=Explore) to independently verify before presenting. Include the sub-agent's evidence in the **Verified** section.

**Ignore**: Style, naming.

---

### Phase 3: Test Coverage

**Focus**: Test file changes. Runs before style because it requires semantic focus from Phase 2.

Check: tests added? names describe behavior? assertions verify intent? edge cases? test quality (pass for the right reasons)? coverage of Phase 2 findings?

---

### Phase 4: Style Guide

**Focus**: Naming, idioms, formatting.
**Reference**: See [standards.md](standards.md).

---

## Newbie Lens (optional)

User says "run the newbie lens" or "newbie check." See [reference.md](reference.md#newbie-lens).

---

## Output

See **Pacing Rule** at top.

1. Read `progress.json` to determine which finding is next (or start from 01 if no cursor)
2. Read that ONE finding file. Do not Read any other finding files in this response.
3. Present it using the format below. That is your complete task.
4. Update `progress.json` with the filename you just presented.

```markdown
### <prefix><ID>, Issue [N] of [total], [Phase emoji] [Phase name]

- **Location**: `file:line`
- **Confidence**: [score]/100
- **Impact**: [why it matters]
- **Problem**: [description]
- **Verified**: [Yes/No, how you confirmed]
```

**Per-issue comment history**: If another reviewer already raised this correctly, skip it (see Comment Deduplication). If presenting a finding with prior discussion, show a timeline section before the problem description.

**Summary** (after all issues):
```markdown
## <prefix><ID>, Summary

| Phase | Issues |
|-------|--------|
| Scope | N |
| Design | N |
| Behavior | N |
| Tests | N |
| Style | N |

## Acceptance Criteria Coverage
[If a JIRA ticket is linked: list each criterion as Delivered / Progressed / Untouched / Contradicted. Multi-PR delivery is expected; Untouched is not a finding by itself.]

## Existing Comments
[Agree / Disagree / Resolved + rationale, including deduped findings]

## Recommendation
**[Approve / Comment]**, [one sentence]
```

**Never use Request Changes.** If findings are critical, say so clearly but let the author decide.

---

## Comment Drafting

See [comment-drafting.md](comment-drafting.md). When the user posts a comment, compare it to your draft and learn patterns per [reference.md](reference.md#comment-learning).

## Comment Verification

When the user asks to check their comments, see [reference.md](reference.md#comment-verification-phase-5).

## Cleanup

After verification or on request: `git worktree remove .reviews/pr-<ID> --force && git branch -D review-<ID>`

## Self-Improvement

See [reference.md](reference.md#self-improvement). Update skill files immediately when the user corrects behavior.

## Constraints

1. **READ-ONLY**: no edits to repo code unless user says "fix this"
2. **PR changes only**: don't comment on pre-existing code
3. **Never modify PR state**: no `gh pr review`, `gh pr comment`, `glab mr approve`
4. **Skill files are fair game**: updating this skill's own files based on feedback is encouraged
