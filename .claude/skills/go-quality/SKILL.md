---
name: go-quality
description: Run full Go quality gate (lint, test, format, build) before committing or creating a PR
---

Run the full CI suite locally to catch issues before pushing:

```bash
task ci
```

This runs in order: `lint` → `test` (with -race) → `fmt` (gofumpt check) → `build`.

Report any failures and suggest fixes:
- **Lint failures**: try `task lint-fix` for auto-fixable issues; otherwise fix manually
- **Test failures**: investigate the failing test, check the error, fix the root cause
- **Format failures**: run `task fmt` to auto-fix
- **Build failures**: fix compilation errors before proceeding

Do not proceed with a commit or PR until `task ci` exits 0.
