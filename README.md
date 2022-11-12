# git-cherry-pick-without-worktree

This is a git-cherry-pick-without-worktree demo.

With `git merge-tree --write-tree`, we can implement `git cherry-pick` by merging given commit.

GIT_ALTERNATE_OBJECT_DIRECTORIES=$(PWD)/objects GIT_OBJECT_DIRECTORY=$(PWD)/objects/incoming-xxxxxx git-cherry-pick-without-worktree <upstream> <branch>