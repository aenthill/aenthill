#!/bin/bash

set -e

read -p "Tag to delete " TAG

git push --delete origin $TAG

git tag -d $TAG

exit 0