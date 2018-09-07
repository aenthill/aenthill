#!/bin/bash

echo -e "\033[0;32mDeploying documentation updates to GitHub...\033[0m"

cd docs
hugo

cd public
git add .

msg="rebuilding site `date`"
if [ $# -eq 1 ]
  then msg="$1"
fi

git commit -m "$msg"
git push origin master

cd ../..

git submodule update --recursive --remote
git add docs/public
git commit -m "$msg"
git push origin master