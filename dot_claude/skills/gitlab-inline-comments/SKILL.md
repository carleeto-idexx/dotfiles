---
name: gitlab-inline-comments
description: Post inline (file/line) review comments on GitLab merge requests via glab. Use whenever the user asks to post, draft-then-post, or reply to an inline MR comment on GitLab. Triggers on phrases like "post as inline comment", "comment on line X", "reply inline", or any MR review workflow that needs a comment anchored to a specific diff line (not a general MR discussion).
---

# GitLab inline MR comments via glab

## The trap (read this every time)

`glab api -f 'position[base_sha]=xxx'` does **not** produce a nested JSON object. `glab api`'s `-f` / `--field` / `--raw-field` all JSON-encode each flag as a top-level string key. The server receives `{"position[base_sha]": "xxx"}`, silently ignores it, and creates a plain `DiscussionNote` with `position: null` (a regular MR comment, not inline).

Symptoms that you fell into the trap:
- Response `.notes[0].type` is `"DiscussionNote"` (should be `"DiffNote"`).
- Response `.notes[0].position` is `null`.
- User says "that's not inline".

## The correct method

Use `--input <file>` with a JSON body AND set `Content-Type: application/json`. Without the header, GitLab returns `HTTP 415 {"error":"The provided content-type '' is not supported."}`.

### Step 1. Gather the required context

```bash
# Project ID and diff refs (base_sha, start_sha, head_sha)
glab mr view <IID> --output json | jq '{project: .target_project_id, iid: .iid, diff_refs}'
```

You also need:
- `new_path`: the file path as it appears in the MR's new tree (repo-root relative).
- `new_line`: the line in the **new** file you want to anchor to. Verify by reading the file in the worktree (line numbers shift from the diff header). **Do not cite a line you haven't read.**
- `old_path`: same as `new_path` for a modified file. For a newly added file, omit `old_line` but still pass `old_path` equal to `new_path`.

### Step 2. Build the JSON body

```json
{
  "body": "your markdown comment here",
  "position": {
    "base_sha": "<from diff_refs.base_sha>",
    "start_sha": "<from diff_refs.start_sha>",
    "head_sha": "<from diff_refs.head_sha>",
    "position_type": "text",
    "new_path": "path/to/file.ts",
    "old_path": "path/to/file.ts",
    "new_line": 42
  }
}
```

Variants:
- **Comment on a removed line** (the `-` side of the diff): use `"old_line": N` instead of `"new_line"`, and keep `old_path`.
- **Comment on an unchanged context line**: pass both `old_line` and `new_line`. GitLab needs both to anchor to an unchanged line.
- **Image position**: `"position_type": "image"` with `x`, `y`, `width`, `height` instead of line numbers.

Write the body to a temp file. Heredocs are safest for preserving backticks and newlines:

```bash
cat > /tmp/inline_note.json <<'EOF'
{ ... }
EOF
```

If the comment body has dynamic content with special chars, build the JSON with `jq` so quoting is correct:

```bash
jq -n --arg body "$COMMENT_BODY" --arg base "$BASE_SHA" ... \
  '{body: $body, position: {base_sha: $base, ...}}' > /tmp/inline_note.json
```

### Step 3. POST

```bash
glab api "projects/<PROJECT_ID>/merge_requests/<IID>/discussions" \
  --method POST \
  -H "Content-Type: application/json" \
  --input /tmp/inline_note.json
```

### Step 4. Verify it stuck

Always check the response. If it silently degraded to a non-inline note, retry before telling the user it posted inline.

```bash
# Pipe the POST response through:
jq '.notes[] | {type, position_present: (.position != null), new_line: .position.new_line, new_path: .position.new_path}'
```

Expected: `type: "DiffNote"`, `position_present: true`, matching `new_line` and `new_path`.

If you see `DiscussionNote` / `position_present: false`, delete it and retry:

```bash
glab api "projects/<PROJECT_ID>/merge_requests/<IID>/discussions/<DISCUSSION_ID>/notes/<NOTE_ID>" -X DELETE
```

## Replying to an existing inline discussion

Reply in the same thread (stays inline, no position needed):

```bash
glab api "projects/<PROJECT_ID>/merge_requests/<IID>/discussions/<DISCUSSION_ID>/notes" \
  --method POST \
  -H "Content-Type: application/json" \
  --input <(jq -n --arg body "$REPLY" '{body: $body}')
```

## Resolving a thread

```bash
glab api "projects/<PROJECT_ID>/merge_requests/<IID>/discussions/<DISCUSSION_ID>" \
  --method PUT -f resolved=true
```

(Plain boolean field, no nesting needed, so `-f` is fine here.)

## When inline is not possible

Inline (DiffNote) comments can only anchor to lines that appear in the MR's diff (added, removed, or adjacent context). If you try to anchor to a file or line that isn't in the diff (common when commenting on orphaned / dead code that the PR doesn't touch), GitLab returns:

```
400 Bad request - Note {:line_code=>["can't be blank", "must be a valid line code"]}
```

Fallback: post a general MR note via `POST /projects/:id/merge_requests/:iid/notes` with `{"body": "..."}`. Name the file path in the body so readers can navigate. This is the correct choice, do not try to work around by attaching a fake position to a line that happens to be in the diff.

## What NOT to do

- Do not use `-f 'position[base_sha]=...'`. It silently fails.
- Do not read raw tokens out of `~/.config/glab-cli/config.yml` or `glab auth status -t`. `glab api` already handles auth.
- Do not trust line numbers from the diff header. Read the file at `HEAD` of the MR and count.
- Do not claim "posted inline" without verifying `type == "DiffNote"` in the response.

## Quick reference: full working example

```bash
# 1. Get refs
REFS=$(glab mr view 181 --output json)
PROJECT=$(echo "$REFS" | jq -r .target_project_id)
BASE=$(echo "$REFS" | jq -r .diff_refs.base_sha)
START=$(echo "$REFS" | jq -r .diff_refs.start_sha)
HEAD=$(echo "$REFS" | jq -r .diff_refs.head_sha)

# 2. Build body
jq -n \
  --arg body "Your comment here" \
  --arg base "$BASE" --arg start "$START" --arg head "$HEAD" \
  --arg path "src/features/admin/hooks/useUser.ts" \
  --argjson line 18 \
  '{body: $body, position: {base_sha: $base, start_sha: $start, head_sha: $head, position_type: "text", new_path: $path, old_path: $path, new_line: $line}}' \
  > /tmp/inline_note.json

# 3. Post
glab api "projects/$PROJECT/merge_requests/181/discussions" \
  --method POST -H "Content-Type: application/json" \
  --input /tmp/inline_note.json \
  | jq '.notes[] | {type, new_line: .position.new_line}'

# 4. Expect: {"type": "DiffNote", "new_line": 18}
```
