Contributing to goamz
=====================

We encourage everyone who is familiar with the [Amazon Web Services
API](http://aws.amazon.com/documentation/) and is willing to support
and improve the project to become a contributor. Current list of
official maintainers can be found in the [AUTHORS.md](AUTHORS.md)
file.

This file contains instructions and guidelines for contributors.

Code of conduct
---------------

We are committed to providing a friendly, safe and welcoming environment
for all users and community members. Please, review the [Ubuntu Code of Conduct](https://launchpad.net/codeofconduct/1.1),
which covers the most important aspects we expect from contributors:

 * Be considerate - treat others nicely, everybody can provide valuable contributions.
 * Be respectful of others, even if you don't agree with them. Keep an open mind.
 * Be collaborative - this is essential in an open source community.
 * When we disagree, we consult others - it's important to solve disagreements constructively.
 * When unsure, ask for help.
 * Step down considerately - if you need to leave the project, minimize disruption.

Contributing a patch
--------------------

Found a bug or want to suggest an improvement?
Great! Here are the steps anyone can follow to propose a bug fix or patch.

 * You need a [GitHub account](https://github.com/signup/free) if you don't have one.
 * [Fork](https://help.github.com/articles/fork-a-repo/) the go-amz/amz repository.
 * If you found a bug, please check the existing [issues](https://github.com/go-amz/amz/issues) to see if it's a known problem. Otherwise, [open a new issue](https://github.com/go-amz/amz/issues/new) for it.
 * Clone your forked repository locally.
 * Switch to the `v1` branch.
 * Create a feature branch for your contribution.
 * Be sure to test your code changes.
 * Push your feature branch to your fork.
 * Open a pull request with a description of your change.
 * Ask a maintainer for a code review.
 * Reply to comments, fix issues, push your changes. Depending on the size of the patch, this process can be repeated a few times.
 * Once you get an approval and the CI tests pass, ask a maintainer to merge your patch.

Becoming an official maintainer
-------------------------------

Thanks for considering becoming a maintainer of goamz! It's not
required to be a maintainer to contribute, but if you find yourself
frequently proposing patches and can dedicate some of your time to
help, please consider following the following procedure.

 * You need a [GitHub account](https://github.com/signup/free) if you don't have one.
 * Review and sign the Canonical [Contributor License Agreement](http://www.ubuntu.com/legal/contributors/). You might find the [CLA FAQ](http://www.ubuntu.com/legal/contributors/licence-agreement-faq) page useful.
 * Request to become a maintainer by contacting one or more people in [AUTHORS.md](AUTHORS.md).

General guidelines
------------------

The following list is not exhaustive or in any particular order. It
providers things to keep in mind when contributing to goamz. Be
reasonable and considerate and please ask for help, if something is
not clear.

 * Commit early, commit often.
 * Before pushing your changes for the first time, use `git rebase -i v1` to minimize merge conflicts. Do not use `git pull v1`, use `git fetch` instead to avoid merging.
 * Rebase and squash small, yet unpushed changes. Let's keep the commit log cleaner. 
 * Do not rebase commits you already pushed, even when in your own fork. Others might depend on them.
 * Write new tests and update existing ones when changing the code. All changes should have tests, when possible.
 * Use `go fmt` to format your code before pushing.
 * Document exported types, functions, etc. See the excellent [Effective Go](http://golang.org/doc/effective_go.html) style guide, which we use.
 * When reporting issues, provide the necessary information to reproduce the issue.
