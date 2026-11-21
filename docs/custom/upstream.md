# Upstream synchronization

| Branch | Role |
| --- | --- |
| `upstream/main` | Canonical OpenChore at `liftedkilt/openchore`. |
| `origin/main` | Fast-forward-only mirror of upstream. Keep customization commits on `custom`. |
| `origin/custom` | Reviewed integration branch for the fork. |
| `feature/*` | Branch from current `origin/custom` and open a PR targeting `custom`. Use `feature/sync-upstream-<sha>` for upstream integrations. |

In a clean development clone, add the upstream remote once and fetch both
remotes. At the M1 follow-up baseline, both main branches pointed to
`4173074` and `custom` included the first M1 PR.

```bash
git remote -v
git remote add upstream https://github.com/liftedkilt/openchore.git
git fetch origin
git fetch upstream
git log --left-right --oneline origin/main...upstream/main
git merge-base --is-ancestor origin/main upstream/main
```

If the ancestor check succeeds and upstream is ahead, a maintainer can
fast-forward the mirror with `git push origin upstream/main:refs/heads/main`
after reviewing what the update will trigger. A push to `main` can start the
inherited [image publishing](../../.github/workflows/build.yml) and
[Pages](../../.github/workflows/pages.yml) workflows. Check repository
permissions and workflow settings before a mirror update. If the branches
diverged, investigate the commits and remotes before changing either ref.
Use ordinary pushes, preserving the branch history.

To integrate the mirror into `custom`, first confirm that `origin/main`
matches the chosen upstream revision, then inspect the divergence. If the
mirror is already an ancestor of `origin/custom`, the sync is complete.

```bash
git fetch origin
git merge-base origin/custom origin/main
git log --left-right --oneline origin/custom...origin/main
git diff --stat origin/custom...origin/main
```

With a common ancestor and new upstream commits, prepare the integration.

```bash
sync_branch="feature/sync-upstream-$(git rev-parse --short origin/main)"
git switch -c "$sync_branch" origin/custom
git merge --no-ff --no-commit origin/main
```

Review user models, API and auth, scoped styling, Compose, workflow triggers
and SQL migrations. If both branches introduce a migration version, test
upgrade paths from disposable databases representing both histories.
Use `git merge --abort` if the in-progress merge needs to be abandoned.
Run the [development checks](development.md), including a synthetic database
upgrade from the previous custom release when migrations change.

```bash
git diff --cached --check
git commit -m "Merge upstream OpenChore into custom"
git push -u origin "$sync_branch"
gh pr create --base custom --head "$sync_branch" --title "Sync OpenChore upstream"
```

Record both input SHAs, conflicts resolved and verification in the PR.
Use a reviewed merge that retains upstream ancestry for sync PRs.
If a sync is reverted through a reviewed PR, account for the revert before
reintroducing the same upstream commits. Evaluate database compatibility
separately before rolling back application code.
