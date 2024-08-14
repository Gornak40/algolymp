#!/bin/bash

if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <contestId> <filter> <count>"
    exit 1
fi

contestId=$1
filter="($2)&&status==PR"
count=$3

if [ -d "$contestId" ]; then
  rm -r $contestId
fi

echo "Filtering runs from contest [$contestId] with filter [$filter] and limit [$count]"
echo "Reviewing [$(boban -i "$contestId" -f "$filter" -c "$count" -d . | wc -l)] filtered runs"

cleanup() {
  if [ -f "$current" ]; then
    rm "$current"
  fi
  rm -r "$contestId"
}

for file in "$contestId"/*; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        runId="${filename%.*}"
        lang="${filename##*.}"

        current="./main.$lang"
        cp "$file" "$current"
        echo "Review $file => $current"

        while [ 1 ]; do
            read -p "Type the resolution (OK, RJ): " verdict
            verdict=$(echo "$verdict" | tr '[:lower:]' '[:upper:]') # Преобразование в верхний регистр
            if [ "$verdict" = "OK" ]; then
                echo $runId | ripper -i "$contestId" -s "OK"
                break
            elif [ "$verdict" = "RJ" ]; then
                read -p "Any comment?: " comment
                echo $runId | ripper -i "$contestId" -s "RJ" -c "$comment"
                break
            elif [ -z "$verdict" ]; then
                cleanup
                exit 0
            fi
        done

        rm $current
    fi
done

cleanup