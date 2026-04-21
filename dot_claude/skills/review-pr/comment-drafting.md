# Comment Drafting Guidelines

Write comments in the reviewer's voice.

**Tone**: Light-hearted but clear and firm. State the problem directly without hedging. Be a colleague, not a linter — human, occasionally funny, but never ambiguous about what's wrong.

**Style guidelines**:
- Use markdown formatting: tables for comparisons, `backticks` for code/identifiers, **bold** for emphasis, numbered lists for sequences
- State the problem plainly, then explain why it matters ("Are we ok with this?" not "Perhaps we should consider...")
- Use emoji or casual asides when they land naturally (`:grinning:`, "Could be kinda confusing")
- For subtle issues, walk through the request timeline step-by-step with numbered sequences and ✓/✗ markers to make the impact undeniable
- For parity or comparison issues, use a table showing expected vs actual behavior per case
- Keep simple issues short — one or two sentences. Use longer walkthroughs only when the bug isn't obvious
- For minor or non-blocking issues, lead with an emoji (e.g., `:smile:`) before the suggestion — keeps the tone light and signals it's not a hard blocker
- When the author already understands the context, skip the explanation and lead with the suggestion. Don't re-explain what they already know — just ask the question or propose the change.
- When drafting non-blocking suggestions, cut the analysis entirely. Don't explain *why* the current code is the way it is. Just state the suggestion and why it's useful. The author wrote the code; they know the backstory.
- Avoid framing suggestions as things the author is "defending" or needs to justify. Keep it collaborative.
- Frame compatibility breaks as questions to the author ("Are we ok with this?" not "You must fix this")
- When referencing other MRs or recurring issues, be matter-of-fact ("Same bug as !336")
- For recurring issues already explained in detail on an earlier MR in the stack, keep the comment minimal — just name the problem (e.g., "Double error logging"). The author already has context from the first comment.
- End with the consequence, not the fix — let the author decide how to solve it

**Example — subtle issue (cache poisoning via partial fetch)**:
```
The happy path with the current query:
1. Client sends PATCH /users/{id} with {"name": "New Name"}
2. Handler calls GetByID — queries users LEFT JOIN orgs → user has all fields including org_name: "Acme", role: "admin"
3. Patch applied on top of complete user → name changes, org fields untouched
4. DB updated, cache updated with complete user (org fields intact)
5. Response returns user with org_name: "Acme", role: "admin" ✓
6. Next GET /users/{id} → cache hit → returns complete user ✓

What happens with the new query:
1. Client sends PATCH /users/{id} with {"name": "New Name"}
2. Handler calls SelectByID — queries users only, no JOIN → user has org_name: "", role: ""
3. Patch applied on top of incomplete user → name changes, org fields are already empty
4. DB updated (org fields aren't in the update record, so DB is fine), cache updated with org_name: "", role: "" ← poisoned
5. Response returns user with org_name: null, role: null ✗
6. Next GET /users/{id} → cache hit → returns user with org_name: null ✗
7. Cache stays poisoned until a List query refreshes it

So the PATCH itself returns wrong data, and it also corrupts reads that follow.
```

**Example — compatibility break (short)**:
```
GetOrderHistory returns 500 on error, but the existing endpoint returns 404. Are we ok with this?

This is a compatibility break, so just checking.
```

**Example — minor issue with humor**:
```
This results in a 404 error with a body that says internal server error. Could be kinda confusing. :grinning:
```

---

## Posting Inline Comments on GitLab

To post a **DiffNote** (inline comment on a specific line in the diff), you must use `glab api --input -` with a JSON body and a `Content-Type: application/json` header. The `-f` form-field approach does **not** produce inline comments; it silently creates a regular discussion note instead.

**Recipe:**

```bash
HEAD_SHA=$(git rev-parse review-<ID>)
BASE_SHA=$(git merge-base origin/<target_branch> review-<ID>)

cat > /tmp/mr-comment.json << ENDJSON
{
  "body": "Your comment text here with \`backticks\` escaped.",
  "position": {
    "base_sha": "${BASE_SHA}",
    "start_sha": "${BASE_SHA}",
    "head_sha": "${HEAD_SHA}",
    "position_type": "text",
    "new_path": "path/to/file.go",
    "old_path": "path/to/file.go",
    "new_line": 42
  }
}
ENDJSON

cat /tmp/mr-comment.json | glab api --method POST \
  -H "Content-Type: application/json" \
  "projects/<url-encoded-project>/merge_requests/<ID>/discussions" \
  --input -
```

**Key rules:**
- Use `--input -` with piped JSON, not `-f` form fields. `-f` silently drops the position and creates a top-level comment.
- Set `-H "Content-Type: application/json"` explicitly.
- `old_path` must always be set (same as `new_path` for modified files).
- For added lines (no corresponding old line), omit `old_line` from the JSON entirely.
- For deleted lines, omit `new_line` and set `old_line`.
- Verify success: the response should have `notes[0].type == "DiffNote"` and `notes[0].position.new_line` set. If `type` is `"DiscussionNote"` or `position` is null, it failed silently.
- Write the JSON to a temp file first to avoid shell escaping issues with backticks and quotes in the comment body.
