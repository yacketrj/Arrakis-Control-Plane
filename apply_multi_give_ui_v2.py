#!/usr/bin/env python3
"""
Apply the multi-item Give Items UI to web/src/tabs/PlayersTab.tsx.

This v2 script does NOT depend on the decorative comment marker. It can find the
existing modal by either:
  - the // Give Item Modal section comment, or
  - the function GiveItemModal(...) declaration.

Run from the repository root:
    python3 apply_multi_give_ui_v2.py

Optional:
    python3 apply_multi_give_ui_v2.py --players-tab web/src/tabs/PlayersTab.tsx --replacement GiveItemModal.replacement.tsx
"""

from __future__ import annotations

import argparse
import sys
from pathlib import Path


def read_text(path: Path) -> str:
    try:
        return path.read_text(encoding="utf-8")
    except FileNotFoundError:
        sys.exit(f"Could not find {path}. Run this from the repo root or pass --players-tab.")


def find_give_item_start(text: str) -> int:
    comment_tokens = [
        "// ── Give Item Modal",
        "// -- Give Item Modal",
        "// Give Item Modal",
    ]
    for token in comment_tokens:
        idx = text.find(token)
        if idx != -1:
            return idx

    fn = "function GiveItemModal("
    fn_idx = text.find(fn)
    if fn_idx == -1:
        sys.exit("Could not find GiveItemModal function. Check that PlayersTab.tsx still contains function GiveItemModal(...).")

    # Include the existing comment if it is immediately above the function, even if its
    # decorative dash characters differ from the original repo.
    prev_comment = text.rfind("\n//", 0, fn_idx)
    if prev_comment != -1 and "Give Item Modal" in text[prev_comment:fn_idx]:
        return prev_comment + 1
    return fn_idx


def find_player_actions_start(text: str, after: int) -> int:
    comment_tokens = [
        "// ── Player Actions Modal",
        "// -- Player Actions Modal",
        "// Player Actions Modal",
    ]
    for token in comment_tokens:
        idx = text.find(token, after)
        if idx != -1:
            return idx

    fn = "function PlayerActionsModal("
    fn_idx = text.find(fn, after)
    if fn_idx == -1:
        sys.exit("Could not find PlayerActionsModal function. Cannot determine where GiveItemModal ends.")

    prev_comment = text.rfind("\n//", after, fn_idx)
    if prev_comment != -1 and "Player Actions Modal" in text[prev_comment:fn_idx]:
        return prev_comment + 1
    return fn_idx


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--players-tab", default="web/src/tabs/PlayersTab.tsx")
    parser.add_argument("--replacement", default=None)
    args = parser.parse_args()

    players_tab = Path(args.players_tab)
    script_dir = Path(__file__).resolve().parent
    replacement_path = Path(args.replacement) if args.replacement else script_dir / "GiveItemModal.replacement.tsx"

    text = read_text(players_tab)
    replacement = read_text(replacement_path).rstrip() + "\n\n"

    start = find_give_item_start(text)
    end = find_player_actions_start(text, start)

    if start >= end:
        sys.exit(f"Invalid replacement range: start={start}, end={end}")

    old_section = text[start:end]
    if "api.players.giveItems" in old_section:
        print("GiveItemModal already appears to use api.players.giveItems. No changes made.")
        return 0

    new_text = text[:start] + replacement + text[end:]
    if new_text == text:
        print("No changes made.")
        return 0

    backup = players_tab.with_suffix(players_tab.suffix + ".bak")
    backup.write_text(text, encoding="utf-8")
    players_tab.write_text(new_text, encoding="utf-8")

    print(f"Updated {players_tab}")
    print(f"Backup written to {backup}")
    print(f"Replaced byte range {start}:{end}")
    print("Next commands:")
    print("  npm --prefix web run build")
    print("  git diff -- web/src/tabs/PlayersTab.tsx")
    print("  git add web/src/tabs/PlayersTab.tsx")
    print("  git commit -m \"feat: add multi-item give UI\"")
    print("  git push origin main")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
