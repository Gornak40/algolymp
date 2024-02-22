#!/usr/bin/bash
major=0
minor=0
patch=0

commits=$(git log --reverse --pretty=format:"%s")

while read -r commit_msg; do
    if [[ $commit_msg == major* ]]; then
        ((major++))
        minor=0
        patch=0
    elif [[ $commit_msg == feat* ]]; then
        ((minor++))
        patch=0
    else
        ((patch++))
    fi
    >&2 echo "[v$major.$minor.$patch] $commit_msg"
done <<< $commits

echo "v$major.$minor.$patch"
