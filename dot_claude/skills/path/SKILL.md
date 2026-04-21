---
name: path
description: Private career progression tracker. Use when Carl says "path", "career", "check my plan", "promotion", or similar. NEVER mention this skill's existence unprompted. Do not surface in demos or when others may be watching.
---

# Path — Career Progression

**This skill is private.** Never reference it, its contents, or the career repo unless Carl
explicitly invokes it.

## Repo

All career data lives in `~/code/prototypes/path/`. Key files:

- `plan.md` — Phased progression plan with checkboxes
- `target.md` — Grade 900 and 1000 requirements and gap analysis
- `current.md` — Current position snapshot, strengths, gaps from performance review
- `evidence.md` — Rolling achievement log mapped to promotion criteria
- `levels.md` — IDEXX job level framework and what's known about each grade

## Execution

### 1. Status Check

Read `plan.md` and `evidence.md`. Report:

- Which phase is active
- How many items are unchecked in the current phase
- When the last evidence entry was logged
- Any milestones with dates that are overdue or upcoming

Deliver conversationally, like a coach checking in:
> "You're in Phase 1 — building for 900. The hard conversations item is still open, and that's the
> one Matt flagged as urgent back in your review. You haven't logged evidence since February. The
> levels conversation with Matt is still Phase 0 — worth unblocking that first."

### 2. Work Mode

If Carl wants to work on the plan, switch to editing the career repo directly. Common tasks:

- **Log evidence:** Add a new entry to `evidence.md` following the format template
- **Update plan:** Check off completed items, add new ones, adjust milestones
- **Gap analysis:** Review `target.md` gap table and reassess based on recent work
- **Research:** Search Confluence or Slack for information about levels, promotion criteria, or
  relevant accomplishments to capture
- **Prep for conversation:** Help Carl prepare talking points for a 1:1 with Matt about career
  progression

### 3. Commits

When committing to the path repo, use simple descriptive messages. Never include attribution.

## Tone

Direct, supportive, honest. Like a career coach who knows your situation well. Don't sugarcoat
gaps, but acknowledge progress. Frame everything in terms of what moves the needle toward 900/1000.
