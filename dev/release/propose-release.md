This is the instruction for propose a new release for config connector.

Please use the VERSION environment variable and STALE_VERSION environment variable. If not set, find the correct new version that should be proposed from `git tag`. For stale
version, if not set, find the latest version to be the stale version.


1. Please checkout to the release branch release-${VERSION} locally. For example, if we are releasing 1.132.0, we should run the command `git checkout release-1.132.0`.
2. Please run ../tasks/propose-tag script and commit all the changed files, including any addition and deletions, especially in the operator/ folder, with the commit message "Release ${VERSION}". Please do not use `git commit -a -m`. Use `git add .` first and then
commit the change.
3. Please run ../tasks/sync-crds-folder.sh script and commit all the changed files, including any additions and deletions, especially in the crds/ folder, with the commit message "Update alpha CRDs for Release ${VERSION}". Please do not use `git commit -a -m`. Use `git add .` first and then
commit the change.
4. Please create a PR with the above two commit to the release branch. Please use
`git push <remote_name> <local_branch_name>:<remote_branch_name>`. For example, if the VERSION is 1.132.0, assuming the <remote_name> is "upstream", the correct command is `git push upstream/release-1.132.0 release1.132.0`.