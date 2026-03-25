#!/usr/bin/env python3
"""Update the app table in a PR body.

Reads APP_NAME, NEW_VERSION, and SOURCE_PR_URL from environment variables,
the existing PR body from stdin, and prints the updated body to stdout.

Usage:
    APP_NAME=kyverno NEW_VERSION=0.24.1-abc1234 SOURCE_PR_URL=https://... \
        python3 hack/update_pr_body.py < existing_body.txt
"""

import os
import sys

app_name = os.environ["APP_NAME"]
new_version = os.environ["NEW_VERSION"]
source_pr_url = os.environ.get("SOURCE_PR_URL", "")
existing_body = sys.stdin.read()

# Format the source PR link, falling back to a dash if no PR was found
if source_pr_url:
    pr_number = source_pr_url.rstrip("/").split("/")[-1]
    source_link = f"[#{pr_number}]({source_pr_url})"
else:
    source_link = "—"

new_row = f"| `{app_name}` | `{new_version}` | {source_link} |"
table_header = "| App | Version | Pull Request |"
separator = "| --- | ------- | ------------ |"

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
