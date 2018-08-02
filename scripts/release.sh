#!/bin/bash

set -e

NEXT_TAG=`git describe --tags --abbrev=0 | awk -F. '{$NF+=1; OFS="."; print $0}'`

read -p "Tag [$NEXT_TAG]" TAG

if [[ -z "$tag" ]]; then
 TAG="$NEXT_TAG"
fi

git tag -a $TAG -m "$TAG"
git push origin $TAG

goreleaser --rm-dist

exit 0