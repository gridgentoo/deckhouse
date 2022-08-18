#!/bin/sh
DIR=$1

cat "$DIR/manifest.json" | jq -r '.[].RepoTags[0]' | \
  while read image; do
    (set -x; mkdir -p "$image")
    cat "$DIR/manifest.json" | jq -r --arg tag "$image" '.[]| select(.RepoTags[0] == $tag).Layers[]' | \
      while read layer; do
        (set -x; tar -C "$image" --overwrite --exclude='./var/run/*' -xf "$DIR/$layer" .)
      done
done
