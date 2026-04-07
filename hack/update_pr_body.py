#!/usr/bin/env python3
"""Update the app table in a PR body.

Operates in two modes depending on which env vars are set:

Full update (NEW_VERSION is set):
  Inserts or replaces the entire row for APP_NAME.
  Status is always initialised to ⏳ Waiting.
  Env: APP_NAME, NEW_VERSION, SOURCE_PR_URL (optional)

Status-only update (NEW_VERSION is empty, STATUS is set):
  Finds the row for APP_NAME and patches only the status cell.
  Env: APP_NAME, STATUS, SHA

In both modes the existing PR body is read from stdin and the result is written to stdout.
"""

import os
import sys

app_name   = os.environ["APP_NAME"]
new_version = os.environ.get("NEW_VERSION", "")
source_pr_url = os.environ.get("SOURCE_PR_URL", "")
status     = os.environ.get("STATUS", "")
sha        = os.environ.get("SHA", "")

existing_body = sys.stdin.read()

table_header = "| App | Version | Pull Request | Status |"
separator    = "| --- | ------- | ------------ | ------ |"

STATUS_DISPLAY = {
    "failure": "❌ Failed",
    "success": "✅ Passed",
    "waiting": "⏳ Waiting",
    "":        "⏳ Waiting",
}

if new_version:
    # ── Full update mode ────────────────────────────────────────────────────────
    if source_pr_url:
        pr_number = source_pr_url.rstrip("/").split("/")[-1]
        source_link = f"[#{pr_number}]({source_pr_url})"
    else:
        source_link = "—"

    new_row = f"| `{app_name}` | `{new_version}` | {source_link} | ⏳ Waiting |"

    if table_header in existing_body:
        lines = existing_body.split("\n")
        new_lines = []
        app_updated = False
        for line in lines:
            if line.startswith(f"| `{app_name}` |"):
                new_lines.append(new_row)
                app_updated = True
            else:
                new_lines.append(line)
        if not app_updated:
            for i, line in enumerate(new_lines):
                if separator in line:
                    new_lines.insert(i + 1, new_row)
                    break
        print("\n".join(new_lines))
    else:
        base = existing_body.strip()
        if not base:
            base = (
                "## Description\n\n"
                "Updated chart references for the security-bundle.\n\n"
                "---\n"
                "> [!NOTE]\n"
                "> This PR was created by the **Handle Chart Update** workflow, "
                "triggered by repository dispatch events."
            )
        print(f"{base}\n\n## Updated Apps\n\n{table_header}\n{separator}\n{new_row}")

else:
    # ── Status-only update mode ─────────────────────────────────────────────────
    status_display = STATUS_DISPLAY.get(status.lower(), status)

    lines = existing_body.split("\n")
    new_lines = []
    for line in lines:
        if line.startswith(f"| `{app_name}` |"):
            if sha and sha not in line:
                new_lines.append(line)
                continue
            parts = line.rstrip(" |").split(" | ")
            parts[-1] = status_display
            new_lines.append(" | ".join(parts) + " |")
        else:
            new_lines.append(line)
    print("\n".join(new_lines))
