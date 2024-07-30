#!/bin/bash
set -eu

REV="${1:-}"
LATEST_TAG=$(git tag | sort --reverse | head -n 1)
echo "Old tag: $LATEST_TAG"

PY_PREFIX='import sys;
import re;
pattern = re.compile(r"v(\d+)\.(\d+)\.(\d+)(.*)");
r = pattern.match(sys.argv[1]);
if r is None: exit(1);
major, minor, patch, extra = tuple(r.groups());
major, minor, patch = int(major), int(minor), int(patch);'

PY_POSTFIX='print(f"v{major}.{minor}.{patch}{extra}");'

NEW_TAG=""
if [ "$REV" == "" ]; then
  echo "You must supply a rev type: major, minor, patch"
  exit 1
fi

NEW_TAG=$(python3 -c "$PY_PREFIX $REV += 1; $PY_POSTFIX" $LATEST_TAG)
echo "New tag: ${NEW_TAG}"

git tag "${NEW_TAG}"

while true; do
    read -p "Do you want to push this tag? " yn
    case $yn in
        [Yy]* ) git push origin "${NEW_TAG}"; break;;
        [Nn]* ) exit;;
        * ) echo "Please answer y or n.";;
    esac
done
