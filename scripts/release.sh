#!/bin/bash

set -e

NEXT_TAG=`git describe --tags --abbrev=0 | awk -F. '{$NF+=1; OFS="."; print $0}'`

read -p "Tag [$NEXT_TAG]: " TAG

if [[ -z "$TAG" ]]; then
 TAG="$NEXT_TAG"
fi

git tag -a $TAG -m "$TAG"
git push origin $TAG

goreleaser --rm-dist

git add build/package/homebrew-tap build/package/scoop-bucket

msg="publishing artifacts `date`"
if [ $# -eq 1 ]
  then msg="$1"
fi
git commit -m "$msg"

git push origin master

git submodule update --recursive --remote