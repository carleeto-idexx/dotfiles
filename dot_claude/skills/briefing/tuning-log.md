# Tuning Log

| Date | Change | Reason |
|------|--------|--------|
| 2026-02-27 | Initial version | Created skill |
| 2026-02-27 | Check emoji reactions on commitments | Carl uses "done" emoji to mark completed actions — don't surface these |
| 2026-02-27 | Include high-confidence recommendations | Carl wants actionable takes when presenting items, not just summaries |
| 2026-02-27 | Always include clickable Slack permalinks | Carl wants to jump straight to the message for additional context |
| 2026-02-27 | Follow-up reminders via self-DMs | Carl can defer items; stored as Slack messages to self, resurfaced next briefing |
| 2026-02-27 | Renamed to "briefing" (from "daily-briefing") | Can be run multiple times a day, not just mornings |
| 2026-02-27 | Added rolling memory (history.md) | Provides continuity across briefings, pruned to 2 weeks |
| 2026-02-27 | Dynamic search windows based on elapsed time | Reads local time + last briefing timestamp to adjust how far back to search, so weekends/gaps don't miss items |
| 2026-03-09 | Removed duplicate "When to Use" section | Already covered by YAML description; saves tokens per best practices |
| 2026-03-09 | Added rationale to search window constants | Prevents misinterpretation of magic numbers |
| 2026-03-09 | Added "Close Out" verification step | Feedback loop to catch missed items and improve over time |
| 2026-03-12 | Avoid acronyms — spell things out | Carl prefers full terms (e.g., "objectives and key results" not "OKRs") |
| 2026-03-16 | Added career check behind "path" codeword | Opt-in only — never surfaces automatically, safe for demos |
| 2026-03-18 | Check reactions on all messages, not just commitments | Carl reacts (thumbsup, +1) to indicate agreement — without checking, briefing misses that he's already weighed in |
| 2026-03-19 | Switched follow-ups from self-DMs to Me canvas (F0AMKMLUWVA) | Canvases are editable, structured, and checkable — self-DMs are immutable and hard to manage |
| 2026-03-23 | Verify Jira ticket status before presenting overdue alerts | Bot alerts are stale by the time the briefing runs — tickets may already be Done. Check Jira to avoid wasting Carl's time on resolved items |
| 2026-03-23 | Read full Slack threads before presenting items as unresolved | Search results show parent messages without replies — must read the thread to check if Carl or others already responded. Avoids presenting resolved items as needing action |
| 2026-03-23 | Extracted career-check.md and tuning-log.md | Saves ~55 lines of context per invocation. Career check only loaded when triggered. Tuning log only loaded when updating. |
| 2026-03-23 | Auto-update Achievements canvas from briefing items | Silently scan briefing items for promotion-relevant achievements and append to Achievements canvas (F0AAXN1EJ9Z). Maps to career plan priorities. |
| 2026-03-25 | Always use /review-pr skill for PR reviews | Carl wants the structured multi-phase review, not raw gh CLI agents. Saved to memory. |
| 2026-03-25 | ALL briefings: re-read active threads from previous briefing | Slack search indexes lag and threads span multiple days. Always re-read threads/channels from the previous briefing to catch new replies, not just for mid-day follow-ups. |
| 2026-03-25 | history.md now stores active thread references (channel_id/message_ts) | Step 0.6 needs concrete thread IDs to re-read. Text descriptions alone aren't reliable for finding threads. |
| 2026-03-25 | Added step 0.5: read canvas follow-ups as explicit execution step | Was buried in Follow-Up Reminders section, easy for the model to skip. Now in the numbered flow. |
| 2026-03-25 | Added explicit daily code review thread step (g) | The code review thread in #team-devex-internal is the most active daily thread and primary source for Carl's review queue. Was only caught incidentally by channel read. |
| 2026-03-25 | Added explicit group DM step (h) with known channel IDs | High-signal group DMs (Matt/Ryan/Carl, Susan/Ryan/Carl) contain strategic discussions that search often misses. |
| 2026-03-25 | Added close-out step: check replies to messages sent during briefing | Carl may send messages during the briefing (review comments, go/no-go messages). Replies can arrive within minutes and should be caught before closing out. |
