# Review PR Reference

On-demand sections loaded when needed. Not part of the core review flow.

## Batch Mode

When `$ARGUMENTS` contains **multiple MR numbers** (e.g., `399 400 402 403 404`), switch to batch triage mode:

1. **Run all MRs in parallel** using sub-agents (one per MR). Each agent follows the full review process (setup, triage, dedup, phases) and writes findings to `.reviews/pr-<ID>/findings/` per [Findings on Disk](SKILL.md#findings-on-disk). Skips timeline pause.
2. **Read each MR's `summary.json`** and build the triage table:

```markdown
## Safe to Approve

| MR | Title | Notes |
|----|-------|-------|
| !399 | Migrate GET /ordersnew/{uuid}/history | No issues |
| !403 | Migrate DELETE /ordersnew/{uuid} | No issues, cache eviction already fixed |

## Needs Your Attention

| MR | Title | Issues |
|----|-------|--------|
| !400 | Migrate GET /ordersnew/{uuid}/items | 1 high, 1 medium |
| !402 | Migrate PATCH /ordersnew/{uuid} | 1 high |
```

3. **"Safe to Approve"**: no findings, or only low/info items that don't block merge.
4. **"Needs Your Attention"**: at least one finding. Show severity counts, not details.
5. **After the triage table, ask the user which MR to start with.** Do NOT auto-continue. Wait.
6. **When the user picks an MR**, read the first (most severe) finding file from that MR's findings directory and present it. That is your entire task. When the user signals, read the next file. After that MR's findings are done, show its summary and ask which MR to continue with.

---

## Draft MR Detection

Check whether the MR is a draft/WIP:
- GitHub: `gh pr view <ID> --json isDraft -q .isDraft`
- GitLab: check the title for `Draft:` or `WIP:` prefix, or `glab mr view <ID> --output json | jq -r .draft`

**If draft**: Switch to lightweight mode:
- No phases. No critique. No bug hunting.
- Instead, read the diff and ask **3-5 curious questions** about intent, scope, and open decisions. Examples:
  - "What's the plan for the TODO on line 42?"
  - "Will the cache eviction be added in a later MR?"
  - "Is this endpoint meant to match v1 behavior exactly, or is it a fresh design?"
- No verdict (no Approve/Comment/Request Changes).
- End with: "Let me know when this is ready for a full review."

---

## Newbie Lens

**Trigger**: User says "run the newbie lens", "newbie check", or similar.

**Perspective**: Read the diff as if you are someone new to this codebase, seeing these files for the first time. You do not know the project history, the migration context, or the team's conventions. You only know the language.

**Focus**: Things that would make a newcomer stop and say "wait, what?"
- Methods that do more than their name suggests (e.g., `Update` that also inserts history and writes cache)
- Non-obvious control flow (sentinel errors, implicit ordering dependencies)
- Magic values or patterns that only make sense with tribal knowledge
- Missing or misleading comments on complex logic
- Struct fields or function signatures that are hard to reason about without reading three other files

**Not in scope**: Bugs, style, architecture. Those are covered by the standard phases. This lens is purely about "can the next person read this?"

**Output**: Present as a numbered list of questions/observations, not as issues. Keep the tone curious, not critical.

No confidence scores, no verification steps. This is lightweight and subjective by design.

---

## Comment Learning

Every time the user says they've posted a comment (e.g., "commented", "added a comment", "posted"):

1. Fetch the latest comment(s) from the MR.
2. Compare the user's actual comment against any draft you provided and the general patterns in [comment-drafting.md](comment-drafting.md).
3. Look for learnable differences:
   - **Tone**: Did the user soften, sharpen, or restructure your draft? Did they add humor or remove it?
   - **Structure**: Did they use a different format (bullets vs prose, shorter vs longer)?
   - **Content**: Did they add context you missed, remove details you included, or reframe the issue?
   - **Phrasing**: Are there recurring patterns in how the user writes comments (e.g., always leads with the consequence, uses specific idioms)?
4. If you identify a pattern that isn't already captured in `comment-drafting.md`, update it. Only update for patterns you've seen at least twice, or for a single strong signal (e.g., the user completely rewrote your draft in a consistent style).
5. Briefly note what you updated (one sentence) so the user knows the skill is learning.

---

## Comment Verification (Phase 5)

**When**: After the user has added review comments to the PR and asks to check them.

**Steps**:
1. Re-fetch comments: `gh pr view <ID> --comments` or `glab mr view <ID> --comments`
2. Map each new comment to the issues found in Phases 1-4
3. For each issue, report:
   - **Covered**: Comment addresses the issue, or **Missing**: No comment for this issue
   - Quote the relevant comment text
4. Flag any comments that don't correspond to a review finding (new observations by the reviewer)

**Output**:
```markdown
## Comment Verification

| # | Issue | Status | Comment |
|---|-------|--------|---------|
| 1 | [short description] | Covered / Missing | [quote or "-"] |

**Coverage**: X of Y issues commented on.
```

---

## Self-Improvement

This skill should evolve based on feedback. When the user corrects your review behavior, update the relevant skill file immediately.

**Triggers**:
- Corrects a false positive ("is that part of the diff?", "is this also an issue in v1?")
- Asks you to change how you present findings ("explain this as a timeline", "don't summarize")
- Points out a missed verification step ("wasn't there a comment about this?")
- Gives feedback on comment tone or style

**What to update**:
- `SKILL.md` for process changes
- `comment-drafting.md` for tone and style
- `standards.md` for new language-specific rules

**How**: Apply the fix immediately after acknowledging the feedback. If the change is risky (removing content, changing invocation behavior), confirm first.
