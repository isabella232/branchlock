Introduction
============

branchlock will add or remove a branch restriction to a Stash repository.  We use this to ensure that no commits are made to the branch during a release build, which could cause the release build to fail, when the release branch is merged back to develop.

This program will remove the branch restrictions, not just the permission for the user added to the restrictions by this program.  That is, if you already have a branch restriction on branch1 where user1 has permission to write to the branch and then you use branchlock to lock the branch with user 2 having write permission and then unlock the branch, the entire branch permission is removed, for both user1 and user2.

Build
=====

```
go get github.com/xoom/branchlock
godeb go build
```

Run
===

Lock a branch
-------------

To lock a branch the user permitted to write to the branch is required.  For us, this is the user that will run the release build.

```
branchlock -stash-rest-base-url https://git.corp.xoom.com -project-key '~mjensen' -project-slug permissions -username branchlock -password branchlock -branch "refs/heads/develop" -permitted-user branchlock -lock true
```

Unlock a branch
---------------

The full specification of the branch must be used.  For example, you can't just use develop, you need to use refs/heads/develop, or whatever is appropriate for your situation.

```
branchlock -stash-rest-base-url https://git.corp.xoom.com -project-key '~mjensen' -project-slug permissions -username branchlock -password branchlock -branch "refs/heads/develop" 
```