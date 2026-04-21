# Briefing History


## 2026-04-17 ~13:37 NZST

**Summary:** 2-day catch-up (morning) + mid-day update. Today's main arc: WoW brainstorm session Carl scheduled off Kal's "question everything" ask, noon ezyVet-to-IPS decision meeting, Luke Ryan deep Confluence review on Chargebee page. Session landed well: 3 experiments committed with owners and demo date 2026-04-28 at next retro. Susan email to Riste went out with Carl's "stronger on realtime + cc Luke" steer. Luke said "Great work" + red-heart on the Chargebee Confluence page, relevant for Priority 3/4 of the career plan.

**Items covered:**
- Chargebee decision meeting due today noon (PT-5089, due 24 Apr). Page has Option C.
- WoW session landed: 3 experiments. Jono auto-approving small PRs (~30% auto-merge), Pulin AI status/time tracking, Carl generalised time tracking skill. Demo 2026-04-28.
- Luke Ryan inline review on Chargebee page: ~10 comments, mostly positive. Nick Langstone suggested Snowflake as data-source option. Luke then said "Great work" + red-heart.
- Susan drafted Riste email. Carl told her to go stronger on realtime push and cc Luke. Actioned.
- timetrk MR #4 merged. Matt loved it. Jono asked for deploy. Thread exploring VetSoft Deploy + secrets management for it.
- Dom admin-FE MR #181 follow-up after Carl's review (inline comment AI bug led Carl to build a new skill).
- Overdue tickets down 7 to 3. PT-4958 done, PT-4982 cancelled.
- Team brainstorm page 6288408585: I accidentally wiped team-filled claims during an edit. Carl restored from version history.
- Jono dropped linear.app/next link with no context at 11:29. Unclear: AI code review tool? Jira alt? Worth asking.

**Deferred:**
- Noon ezyVet-to-IPS decision meeting outcome (not yet verified after the meeting).
- Jono's Linear link (context unclear).
- timetrk deploy plan (Jono thread paused mid-design).

**Feedback:**
- **Never modify shared documents (Confluence, Slack canvases, Jira) without explicit per-edit approval, even in auto mode.** Saved to memory after I wiped team-filled content from the WoW brainstorm page.
- Carl prefers the softer version of the Riste email ("we're designing toward" vs "design requirement, not a preference") since Susan has the direct relationship.

**Active threads:**
- C084N3LV9AQ/1776366010.346159 (today's daily code review)
- C084N3LV9AQ/1776373205.797649 (today's standup)
- C084N3LV9AQ/1776297613.967049 (time tracking / timetrk deploy thread)
- D08SWM0M91C (Matt DM, last: 2026-04-17 09:33 re timetrk)
- D0AG4QGESE5 (Susan DM, last: 2026-04-17 13:24 re Riste email steer)
- Confluence page 6288408585 (WoW brainstorm, live artifact for next retro)

## 2026-04-15 ~13:03 NZST

**Summary:** 2-day catch-up. VSD wildcard ingress outage (resolved by Jono). Team ramping on Q2 Chargebee planning. Skills matrix updated with 3 new dimensions (Autonomy, Delivery, Design discipline) from 1:1 canvas evidence. Jonathan 1:1 feedback logged.

**Items covered:**
- VSD outage: AI Platform gw-load wildcard ingress hijacked all idexxvs traffic in dev. Jono diagnosed and fixed. Matt notified.
- Overdue tickets: 6 total. PT-4958 (Neeraj, new), PT-4584 (Jonathan), PT-4354 (Pulin), PT-3109/PT-2812/PT-2554 (Jono).
- Matt DM: Dom remote work question. Matt said not an option, unlikely exemptions.
- Rob Richardson: reporting flaky tests in #team-devex, aiming to enable merge trains on ezyVet monolith. Credited Jono.
- Scott Goodhew: needs Tomislav to update v2 tag on vs-github-actions-workflows for vsd-infra.
- Skills matrix: added Autonomy, Delivery, Design discipline dimensions with evidence from all team member canvases.
- Jonathan canvas: logged Apr 15 1:1 entry (communication, cultural context, step change expectation).
- Chargebee technical PRD ("ezyVet to Chargebee Subscription Management"): 1 open inline comment from Susan re Central Structures. Carl replied and updated page.
- Yesterday code review: Pulin admin-frontend labels MR, Dom zod schema upgrade MR.
- Dom merge perms on admin-fe sorted by Jono.

**Deferred:**
- Nothing deferred.

**Feedback:**
- Carl wants briefings to update skills matrix based on observed evidence and track performance review material from canvases. Canvas IDs: Jonathan (F0A2MELGG3Z), Jono (F09T83V0VQE), Neeraj (F09TA59A9KL), Dom (F0A8JKSSKBL), Chris (F0A075Q2TG8), Pulin (none yet).

**Active threads:**
- C08JFPBJ63E/1776198759.056589 (VSD outage, vello admin 401s)
- C084N3LV9AQ/1776193207.276249 (today's daily code review)
- C084N3LV9AQ/1776200406.981349 (today's standup)
- C084N3LV9AQ/1776106807.004999 (yesterday's code review, Pulin + Dom MRs)
- D08SWM0M91C (Matt DM, last: Apr 14 Dom remote work)

## 2026-04-13 ~08:55 NZST

**Summary:** 19-day catch-up (Carl on leave April 3-12, Brisbane trip + Easter). Major development: Chargebee architecture shifted to ezyVet→IPS, Matt confirmed Carl's feedback was actioned. Subscription resolver endpoint now specced. Luke Q2 meeting this week. API migration completed April 2 (KR done). Team ran well during leave, overdue tickets down from 14 to 6. Career check: logged 4 new evidence entries, nudged Phase 0.

**Items covered:**
- Chargebee architecture shift: ezyVet→IPS confirmed, subscription resolver endpoint specced, who does the work still TBD. Luke Q2 meeting this week. Deferred deep dive to tomorrow.
- Jonathan's ownership UID message: doesn't think UIDs are explicitly generated on site creation, references siteops db-restore.sh
- Jono's pinned actions question: updating all GitHub repos, asked if any non-VSD repos. Carl replied.
- Overdue tickets: 6 current (PT-4584 Jonathan, PT-4354 Pulin, PT-4234/PT-3109 Jono VSD, PT-2812/PT-2554 Jono blocked migrations)
- Dominic's Grafana question from April 10: platform channel data refresh issue, no team replies visible
- Canvas follow-ups: parked Susan ownership handoff, Unleash, AI business case. Hire profile awaiting 1pm update. Added Chargebee deep dive follow-up.
- Tempo alert: 176h logged vs 64h planned, Carl cleaned up
- Career check: logged 4 evidence entries (API migration done, PR throughput 3→29, Chargebee feedback actioned, autonomous team direction before leave). Nudged Phase 0 conversation with Matt.

**Deferred:**
- Chargebee deep dive (tomorrow)
- Hire profile update (after 1pm planning today)
- Chargebee POC next steps with Susan (canvas, still open)

**Feedback:**
- Carl will ask for Phase 0 career prep before Wednesday 1:1 with Matt (saved to memory)

**Active threads:**
- C084N3LV9AQ/1776028809.000000 (Jonathan ownership UID message)
- C084N3LV9AQ/1776023404.467029 (today's daily code review)
- D08SWM0M91C (Matt DM, last: April 2 re API migration + Luke meeting)
- C0AJF18LJS0 (Susan+Ryan group DM, last: April 2 re Chargebee v1→v2 confusion)
