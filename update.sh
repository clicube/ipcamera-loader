#!/bin/bash -x

git fetch -p
git clean -f
git reset --hard
git checkout raspi
git reset --hard origin/raspi

./install.sh