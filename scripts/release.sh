#!/bin/bash

set -e

NEXT_TAG=`git describe --tags --abbrev=0 | awk -F. '{$NF+=1; OFS="."; print $0}'`

read -p "Tag [${NEXT_TAG}]" tag

if [[ -z "${tag}" ]]; then
 tag="${NEXT_TAG}"
fi

git tag -a ${tag} -m "${tag}"
git push origin ${tag}

goreleaser --snapshot --rm-dist

exit 0