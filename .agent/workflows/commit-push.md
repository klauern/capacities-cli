---
description: Commit and Push with Conventional Commits
---

1. **Check Branch Status**
   - Run `git branch --show-current` to get the current branch.
   - If the branch is `main` or `master`:
     - **STOP**. Do not commit directly to main/master unless told to directly.
     - If not told to commit to main, ask the user for a new branch name (e.g., `feature/add-auth`, `fix/bug-123`).
     - Create and switch to the new branch: `git checkout -b <new-branch-name>`.

2. **Stage Changes**
   - Run `git status` and `git diff` to see what changed.
   - Stage files using `git add <file>`. Use `git add .` only if you are sure all changes should be committed.

3. **Create Commit**
   - Generate a Conventional Commit message based on the changes.
     - Format: `<type>(<scope>): <description>`
     - Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`.
     - Example: `feat(auth): implement login functionality`
   - Run `git commit -m "<message>"`.

4. **Push Changes**
   - Push the branch to the remote:
     - If it's a new branch: `git push -u origin <branch-name>`
     - If it's an existing branch: `git push`

5. **Verify**
   - Run `git status` to ensure the working directory is clean and the branch is up to date.