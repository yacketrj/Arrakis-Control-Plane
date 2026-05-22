# Apply Multi Give Items UI

1. Copy `apply_multi_give_ui.py` to the root of your local `dune-admin-fork` clone.
2. Run:

```bash
git pull origin main
python3 apply_multi_give_ui.py
npm --prefix web run build
git diff -- web/src/tabs/PlayersTab.tsx
git add web/src/tabs/PlayersTab.tsx
git commit -m "feat: add multi-item give UI"
git push origin main
```

The script replaces the existing single-item modal inside `PlayersTab.tsx` with a multi-item modal.
