---
name: create-jira-issue
description: Use when Carl asks to create, draft, or file a Jira issue/ticket. Captures Carl's drafting conventions, MCP tool quirks, and the workflow for getting tickets, links, and Confluence references right on the first try.
---

# create-jira-issue

Create Jira issues on Carl's behalf cleanly, the first time.

## Workflow

1. **Draft first, create second.** Always present the full draft in conversation (summary, description, parent, priority, assignee, due date, intended links) and wait for explicit approval before calling the Jira API. Editing shared systems without per-edit approval violates Carl's standing rule.
2. **Load MCP tool schemas before calling.** The Atlassian tools are deferred; use `ToolSearch` with `select:<tool_name>` to load them. Core set: `createJiraIssue`, `editJiraIssue`, `getJiraIssue`, `createIssueLink`, `getIssueLinkTypes`, `getJiraIssueRemoteIssueLinks`, `searchJiraIssuesUsingJql`, `lookupJiraAccountId`, `getAccessibleAtlassianResources`.
3. **Find the cloudId once per session** via `getAccessibleAtlassianResources` (or reuse the IDEXX cloudId below).
4. **Create the issue.** Use `createJiraIssue`. After creation, add Jira-to-Jira links via `createIssueLink`. For Confluence/web links, see the "Remote links gap" section below, do not fake it via the description.
5. **Verify with `getJiraIssue`** after changes. Do not assume the write succeeded based on the response alone; re-read the issue to confirm title, parent, links, and any smart-link rendering.

## Carl's drafting conventions

- **Title = action only.** Describe what will be done. No "Follow-up to PT-XXXX", no "Part 2 of", no context framing. The title should read cleanly on a board with no surrounding context.
- **Never link other tickets in the description.** Use Jira issue links (Blocks / is blocked by / relates to / clones) instead. Propose the link type alongside the draft so Carl approves both content and relationships in one pass.
- **Descriptions are terse.** Short lead paragraph, then structured sections (**Acceptance:**, **Due:**, etc.). No prose narrating predecessor context that belongs in issue links or Confluence.
- **Due dates are real.** If Carl names a due date, set it via `additional_fields.duedate` in ISO format (`YYYY-MM-DD`). If the date is contingent (e.g. "Thursday if IPS answers today"), record that nuance in the description, and still set the `duedate` field.
- **Parent epic.** Use the `parent` parameter on `createJiraIssue` with the epic key (e.g. `PT-2585`). Confirm by reading the predecessor ticket's parent if Carl says "same epic as X".
- **Labels are semantic, not inherited.** Start with an empty label set. Do not inherit labels from the parent ticket. The `bst` label specifically means "BST (Business Support Team) should be aware of this work"; only apply it when Carl confirms BST needs visibility. Sub-tasks covering pure implementation under a `bst`-labelled parent should usually be unlabelled.

## IDEXX environment (cache these)

- **cloudId:** `e83eb8ac-b27d-416b-ae94-4335aee7044a`
- **Carl's accountId:** `712020:62eb2a84-8cfe-49e6-854d-40d85b36f93b`
- **Common project:** `PT` (VetSoft Platforms)
- **Jira base URL:** `https://idexx.atlassian.net`

If a project key or accountId is unknown, use `lookupJiraAccountId` or `searchJiraIssuesUsingJql` with `project = X` to confirm before creating.

## createJiraIssue quick reference

```
cloudId: e83eb8ac-b27d-416b-ae94-4335aee7044a
projectKey: PT
issueTypeName: Task                     # or Bug, Story, Epic
summary: "<action-only title>"
parent: "PT-XXXX"                       # epic key for the parent
assignee_account_id: "712020:..."
contentFormat: "markdown"
responseContentFormat: "markdown"
description: |
  <lead paragraph>

  **Acceptance:**
  - ...

  **Due:** ...
additional_fields:
  priority: { name: "Medium" }          # or Low, High, Highest
  duedate: "YYYY-MM-DD"
  labels: ["..."]
```

## Adding Jira-to-Jira links

Use `createIssueLink` **after** the issue is created. Direction matters:

- "A is blocked by B" -> `type: "Blocks"`, `inwardIssue: B` (blocker), `outwardIssue: A` (blocked).
- "A relates to B" -> `type: "Relates"`, direction irrelevant.
- "A clones B" -> `type: "Clones"`, `inwardIssue: B` (original), `outwardIssue: A` (clone).

Use `getIssueLinkTypes` if unsure of the link name. The type parameter is the link-type name (e.g. `Blocks`), not the directional phrase.

## Remote links (Confluence pages, web links)

The MCP only exposes `getJiraIssueRemoteIssueLinks` (read). To create remote links, hit the Jira REST API directly with `curl`, because `acli` doesn't cover remote links either.

**Never** paste Confluence/web URLs into the description and call it "linked." An inline smart-link renders in the description body but does NOT appear in the sidebar "Links" section, and filters/boards that query linked Confluence pages will not see it.

### Step 1: Check for credentials

Carl stores tokens in `~/.tokens` (sourced by `~/.zshrc`). Required env vars:

- `JIRA_EMAIL` (expected: `carl-menezes@idexx.com`)
- `JIRA_API_TOKEN`

Check availability with a Bash one-liner before attempting the call:

```bash
[ -n "$JIRA_EMAIL" ] && [ -n "$JIRA_API_TOKEN" ] && echo "ready" || echo "missing"
```

If the env vars are not set in the current shell but `~/.tokens` has content, they may not have been sourced into the harness's shell yet. Try `source ~/.tokens && <check>` in a single Bash call. If still missing, **ask Carl**:

> "I need a Jira API token to attach the Confluence page as a proper remote link. Can you generate one at https://id.atlassian.com/manage-profile/security/api-tokens and add it to `~/.tokens` as `JIRA_API_TOKEN` (with `JIRA_EMAIL=carl-menezes@idexx.com`)? Tell me when it's in place."

Only ask at the moment the token is needed, not proactively at the start of a session.

### Step 2: Create the remote link

```bash
curl -sS -u "$JIRA_EMAIL:$JIRA_API_TOKEN" \
  -X POST \
  -H "Content-Type: application/json" \
  "https://idexx.atlassian.net/rest/api/3/issue/<ISSUE-KEY>/remotelink" \
  -d '<JSON_BODY>'
```

**For a Confluence page** (matches the sidebar "Wiki Page" style Carl uses):

```json
{
  "globalId": "appId=cc5a64b2-937d-35ff-8a55-de098513dc5f&pageId=<PAGE_ID>",
  "application": { "type": "com.atlassian.confluence", "name": "System Confluence" },
  "relationship": "Wiki Page",
  "object": {
    "url": "https://idexx.atlassian.net/wiki/pages/viewpage.action?pageId=<PAGE_ID>",
    "title": "<Page Title>"
  }
}
```

The `appId` above (`cc5a64b2-937d-35ff-8a55-de098513dc5f`) is IDEXX's Confluence application ID, verified from PT-5392's remote link. Reuse it. Use the `/wiki/pages/viewpage.action?pageId=<ID>` URL form (not the `/wiki/spaces/...` form); the API-created link renders with the canonical viewpage URL.

**For a plain web link** (no app registration, no globalId):

```json
{
  "object": {
    "url": "<URL>",
    "title": "<Link title>"
  }
}
```

### Step 3: Verify

After the curl call, run `getJiraIssueRemoteIssueLinks` via MCP and confirm the new link appears with the expected title, URL, and (for Confluence) `application.type: com.atlassian.confluence` + `relationship: "Wiki Page"`. A 201 Created response from curl is not sufficient proof, verify the read path.

### Step 4: If curl fails

- 401: token is missing, expired, or malformed. Re-ask Carl.
- 403: account lacks remote-link create permission on the project. Fall back to asking Carl to add via UI.
- Any other failure: note the body and stop. Do not retry blindly.

### Fallback: manual UI attachment

If the API route is unavailable (no token, permission denied, etc.), ask Carl to attach via Jira UI: "+ Add" > "Confluence page" or "Web link." Verify via `getJiraIssueRemoteIssueLinks` once he confirms.

## Post-creation checklist

After creating and linking, report to Carl:

- Issue key and URL
- Parent epic (confirmed)
- Priority, assignee, due date (confirmed)
- Jira-to-Jira links added (with types)
- Remote links NOT yet added (flag each Confluence page / web link Carl needs to attach in the UI)

## Anti-patterns to avoid

- Putting ticket references (PT-XXXX) in the summary or description instead of using issue links.
- Pasting Confluence URLs into the description and calling it "linked."
- Creating the issue before showing Carl the draft.
- Forgetting `responseContentFormat: "markdown"` and getting ADF JSON back, which is harder to read for verification.
- Hardcoding dates that are contingent without noting the contingency in the description.
- Skipping `getJiraIssue` verification after an `editJiraIssue` call.

## Updating this skill

If a new convention, tool quirk, or workaround surfaces, edit this file directly. This skill is the durable memory for Jira issue creation, keep it current.
