# Ways of working

 * [Overview](#Overview)
 * [Definition of done](#Definition-of-done)
 * [Issue tracking](#Issue-tracking)
 * [Branches](#Branches)
 * [Commit message](#Commit-message)
 * [git settings](#git-settings)
 * [git workflow](#git-workflow)
 * [Code review](#Code-review)

## Overview

The code is version tracked with git and available on [gitlab](https://gitlab.3fs.si/iryo/wwm) repository hosted on 3fs servers.

Main branch (mainline) is called `master`, containing the latest stable release, which has already been or is in queue to be deployed to production.

Feature development and bug fixing is done in `topic` branches, branched of `master` branch. Upon completion and code review, `topic` branch is merged into `master` branch.

Topic branches must be merged into `master` as soon as possible but not longer than `5 days`. If the work is not completed yet, please consult with the team.

Once a deploy is done the `master` branch is tagged using [semantic versioning](https://semver.org/).

## Definition of done (DoD)

 * code written and validated (eslint)
 * tests written (unit, component, integration)
 * documentation written
 * security guidelines followed
 * successful automatic build with all checks
 * [code reviewed](#Code-review)
 * [code merged and deployed](#Merging-the-code-into-master)

## Issue tracking

Features, bugs and other tasks are tracked in [Trello](https://trello.com/b/7u6QYLFB/development).

Issue need to have as verbose description as possible. Checklists are welcome as well.

Every change needs to have an issue linked in the [commit message](#Commit-message).

## Branches

### Topic branches (features and bug fixes)

Topic branches need to be branched off `master` and named `type/short-name`, where type is `feature`, or `bug`.

For example, a branch name for a feature called `Add support for policies` would be `feature/policy-support` or similar, where a `User cannot login` bug would be `bug/user-cannot-login` or similar.

~~When pushed to remote, `topic` branch is automatically [deployed](#production-staging-environment.md) to an isolated staging environment.~~

### Master branch

`master` branch contains all the changes as defined by [Definition Of Done](#definition-of-done-dod). Topic branches need to be [merged with fast forwarding](#Merging-the-code-into-master).

~~When pushed to remote, `master` branch is automatically [deployed](production-staging-environment.md) to a production environment.~~

## Commit message

Commit message should be formatted as following:

```
part: Capitalized, short (50 chars or less) summary

More detailed explanatory text, if necessary.  Wrap it to about 72
characters or so.  In some contexts, the first line is treated as the
subject of an email and the rest of the text as the body.  The blank
line separating the summary from the body is critical (unless you omit
the body entirely); tools like rebase can get confused if you run the
two together.

Write your commit message in the imperative: "Fix bug" and not "Fixed bug"
or "Fixes bug."  This convention matches up with commit messages generated
by commands like git merge and git revert.

Further paragraphs come after blank lines.

- Bullet points are okay, too

- Typically a hyphen or asterisk is used for the bullet, followed by a
  single space, with blank lines in between, but conventions vary here

- Use a hanging indent

Closes #ID.
```

See [A Note About Git Commit Messages](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html) for details.

The first line of the change description is conventionally a one-line summary of the change, prefixed by the primary affected part of the website, and is used as the subject for code review. The rest of the description elaborates and should provide context for the change and explain what it does. Write in complete sentences with correct punctuation. If there is a helpful reference, mention it here. At the end of the message, reference the Trello issue with an URL (Example https://trello.com/c/MmVjmSTS).

Mind that there is no connection between Trello issues and Merge requests in Gitlab. Issues have to be moved or closed manually once MR is closed.

An example message:

```
version: Add versioning support

Basic Semantic Versioning style version information is all that is
needed to distinct between different builds. So far we have been
releasing the API as Semantic versions, but this information was limited
to pretty much just the developers of the API. With this change, version
is available via `/v1/version` endpoint.

Closes https://trello.com/c/YMGPbIMJ.
```

As a recommendation, use:

 * `Closes` to reference a feature,
 * `Fixes` to reference a bug.

If you are using `vim` as your commit editor, add the following to your `.vimrc` for a visual representation of a commit message:

```
autocmd Filetype gitcommit set textwidth=72
```

## git settings

### Name and email

Ensure the commits are done using your full name and your work email. Allowed working emails are currently `@3fs.si` and `@iryo.io`. Gitlab will link your account properly in the repository when the email is added to your profile.

### GPG signing

All git tags have to be signed, see [Gitlab's help page](https://gitlab.3fs.si/help/user/project/repository/gpg_signed_commits/index.md) how to set up and use it.

Commit signing is encouraged but not required for the time being.

### Good to have

Following configuration should be in your `$HOME/.gitconfig`:

```
[push]
	# prevents your from accidentally pushing to a wrong branch
	default = nothing

[alias]
	dc     = diff --cached
	co     = checkout
	ci     = commit
	cm     = commit -m
	st     = status
	br     = branch
	lg     = !git log --graph --pretty='format:%C(yellow)%h%C(reset) -%C(red)%d%C(reset) %s %Cgreen(%cr) %C(bold blue)<%an>%Creset'
	slog   = log --oneline
	rb     = rebase
	rbc    = rebase --continue
	rbi    = rebase -i
	fap    = fetch --all --prune --progress
	latest = for-each-ref --sort=-committerdate refs/heads refs/remotes --format='%(committerdate:iso8601)  %(refname:short)                        %(authorname)'
```

The most useful alias is `git fap` which fetches the changes from all remotes and at the same time prunes your local copy.


## git workflow

### 1. Create a topic branch

Create and check out a new local branch where you will do your work.

```bash
git checkout -b <topic-branch> origin/master
```

When pushing the branch for the first time, ensure correct tracking is set up:

```bash
git push -u origin <topic-branch>
```

### 2. Keep topic branch up-to-date

When changes have been pushed to `master` branch, ensure your topic branch is up-to-date.

```bash
git fap
git checkout <topic-branch>
git rebase origin/master
```

### 3. Code, test, ...

... do your best :)

### 4. Commit the changes

Changes should be grouped into logical commits. Refer to above [commit message](#commit-message) for details.

```bash
git add -p # only for existing files, when adding a file, use git add <file>
git commit
```

If you want to sign a particular commit, use `git commit -S`. If you want to sign all the commits, add the following to your git config file:

```
[commit]
	gpgsign = true
```

### 5. Push the changes

Changes should be pushed to a correct origin branch:

```
git push -u origin <topic-branch>
```

You may need to force push your topic-branch.

```
git push -fu origin <topic-branch>
```

### 6. Open a merge request

When the changes in your branch are completed, a merge request has to be opened on Gitlab.

Before opening it, please ensure your branch is up-to-date with `master` branch and your commits are properly formatted.

```
git fap
git checkout <topic-branch>
git rebase origin/master
git push -fu origin <topic-branch>
```

### 7. Code review

Code review is done on web through Gitlab's merge request and code review mechanism.

Reviewer should do the following:

- ensure all the automated checks have passed successfully,
- commits are logical (propose rebase when needed),
- [commit messages](#commit-message) are as agreed,
- code is formatted correctly (coding style, coding format),
- propose code modifications where necessary,
- submit review with by clicking the thumb up icon (:+1:),
- submit comments,
- once comments are resolved the reviewer should mark them as resolved.

Code can be merged into the mainline when:

- at least one reviewer gave the merge request a thumbs up (:+1:),
- when no one is requesting changes,
- when there are no thumbs down given (:-1:).

When a change is made after receiving a comment, is has to be done as a new commit, do not rebase existing commits. This ensure the reviewers to verify changes more easily. When criteria to merge the code has been reached, rebase the commits (squash, fix, ...) before merging.

### 8. Merging the code into master

After a successful code review, `topic` branch is merged into `master`. Merging **MUST NOT** be done through the Gitlab's web interface, but rather via your command line applying below commands.

**CAUTION** before pushing to `origin master`, you need to ensure Deploy guidelines are followed. **

```bash
git fetch
git checkout <topic-branch>
git rebase origin/master
git push -fu origin <topic-branch>
git checkout master
git reset --hard origin/master
git merge --ff-only <topic-branch>
git push origin master
git push origin :<topic-branch>
git branch -d <topic-branch>
```

Above commands will ensure:

 * your topic branch is up-to-date with `master` (in case there were changes),
 * force push it to origin (to ensure Gitlab's merge request shows correct merged status),
 * merge topic branch into master,
 * push the changes to origin,
 * remove topic branch locally and remotely.

### 9. Tagging a release

After required changes are accumulated in the `master` the release has to tagged. All tags must signed and [semver versioning](https://semver.org/) has to used.

```bash
git fetch --all --tags
git checkout master
git rebase origin/master
git tag -s v0.5.1
git push origin v0.5.1
```

Above commnds will ensure:

 * you have all remote tags present on your machine
 * your master branch is up to date with the origin
 * a gpg signed tag is created
 * tag is pushed to gitlab
