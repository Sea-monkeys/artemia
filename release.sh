#!/bin/bash
set -o allexport; source release.env; set +o allexport

: <<'COMMENT'
Todo:
- update of release.env:
  - TAG
  - ABOUT

Remark: delete tag: git tag -d v0.0.1
COMMENT

echo "Generating release: ${TAG} ${ABOUT}"

find . -name '.DS_Store' -type f -delete

git add .
git commit -m "📦 ${ABOUT}"
git push

git tag -a ${TAG} -m "${ABOUT}"
git push origin ${TAG}
