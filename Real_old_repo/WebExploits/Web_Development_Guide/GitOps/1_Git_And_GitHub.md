# Git & GitHub for Dummies

A simple, in-order walkthrough: local changes → push → fetch/merge → pull requests → issues.

---

## 1. The Big Picture

Three places your code can live:

```
 WORKING DIR  --add-->  STAGED  --commit-->  LOCAL REPO  --push-->  REMOTE REPO (GitHub)
  (your edits)          (ready)              (saved history)        (shared with team)
```

That's it. Everything below is just different ways of moving code between these boxes.

---

## 2. Setup (once per project)

```bash
git clone <url>          # copy a remote repo to your machine
git checkout -b <name>   # create + switch to a new branch
```

---

## 3. Step 1: Make Changes → Add → Commit

```
 [ Working Dir ]
       │  git add .
       ▼
 [  Staged   ]
       │  git commit -m "message"
       ▼
 [ Local Repo ]
```

```bash
git add .                  # stage your changes
git commit -m "message"    # save them to local history
```

At this point your changes exist **only on your computer**.

---

## 4. Step 2: Push → Send Changes to Remote Branch

```
 [ Local Repo ]  ──git push──►  [ Remote Branch (GitHub) ]
```

```bash
git push -u origin <branch-name>   # first push (sets tracking)
git push                           # every push after that
```

Your teammates (and GitHub) can now see your branch.

---

## 5. Step 3: Fetch & Merge → Bring Remote Changes Down

Use this when others have pushed changes and you need them locally.

```
 [ Remote Branch ] ──git fetch──► [ Local: origin/branch ]
                                          │  git merge
                                          ▼
                                   [ Your Local Branch ]
```

```bash
git fetch origin           # download changes, don't apply them yet
git merge origin/main      # merge them into your current branch
```

> `git pull` = `git fetch` + `git merge` in one step.

---

## 6. Step 4: Pull Requests (PR) — Merging Branches on GitHub

A PR is just a **request** to merge one branch into another. It goes both ways:

```
 Branch ──► Main         "I finished a feature, merge it into main"
 Main   ──► Branch       "Main moved forward, update my branch"
```

**Branch → Main:**
```bash
gh pr create --base main --head my-branch --title "Add feature" --body "..."
```

**Main → Branch** (keep your branch up to date):
```bash
gh pr create --base my-branch --head main --title "Sync with main"
```

---

## 7. Step 5: Review, Approve, and Merge the PR

```
 [ PR opened ] → [ Reviewed ] → [ Approved ] → [ Merged into target branch ]
```

```bash
gh pr list                     # see open PRs
gh pr view <number>            # look at one PR
gh pr diff <number>            # see the code changes
gh pr review <number> --approve
gh pr merge <number> --squash  # squash = combine all commits into one clean commit
```

**Squash merge** keeps history tidy — instead of 10 messy commits, main gets 1 clean commit.

---

## 8. If Something Goes Wrong: Revert vs Reset

```
 revert  →  adds a NEW commit that undoes an old one   (safe, use on shared branches)
 reset   →  moves history backwards, deletes commits   (only use on LOCAL, unpushed work)
```

```bash
git revert <commit>          # safe undo
git reset --hard <commit>    # dangerous rewrite (local only!)
```

---

## 9. Issues (Tracking Work / Bugs)

Issues aren't code — they're just tickets to track tasks, bugs, or ideas.

```bash
gh issue create --title "Bug: login fails" --body "Steps to reproduce..."
gh issue list                       # see all open issues
gh issue view <number>              # read one issue
gh issue edit <number> --add-assignee <name>   # assign someone
gh issue close <number>             # mark it done
```

A normal workflow: open an issue → create a branch for it → fix it → PR the branch into main → close the issue.

---

## 10. Full Cheat Sheet (in order)

```bash
# 1. Setup
git clone <url>
git checkout -b feature-branch

# 2. Edit files, then save your work locally
git add .
git commit -m "message"

# 3. Send it to GitHub
git push -u origin feature-branch

# 4. Keep up to date with remote
git fetch origin
git merge origin/main

# 5. Open a PR (branch -> main)
gh pr create --base main --head feature-branch --title "..." --body "..."

# 6. Review + merge
gh pr list
gh pr review <n> --approve
gh pr merge <n> --squash

# 7. Undo if needed
git revert <commit>

# 8. Track work with issues
gh issue create --title "..." --body "..."
gh issue close <n>
```

---

## 11. One-Line Summary

**Edit → Add → Commit → Push → PR → Review → Squash Merge.**
Use **fetch + merge** to stay in sync, **revert** to undo safely, and **issues** to track what needs doing.
