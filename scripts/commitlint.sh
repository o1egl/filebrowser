#!/usr/bin/env bash
set -e

for commit_hash in $(git log --pretty=format:%H origin/master..HEAD); do
   commitlint -d ./tools -f ${commit_hash}~1 -t ${commit_hash};
done;
