#!/bin/bash

export LC_ALL=en_US.UTF-8

NC='\033[0m'       # Text Reset

# Regular Colors
Red='\033[0;31m'          # Red
Green='\033[0;32m'        # Green
IYellow='\033[0;93m'
Cyan='\033[0;36m'         # Cyan

if [ "$#" -lt 3 ]; then
    echo "Usage: $0 <contestId> <filter> <count> [<0/1 disable PR filter>]"
    exit 1
fi

contestId=$1
if [ $4 == '1' ]; then
  filter=$2
  echo -e "${Red}WARNING${NC}    Pending review filter is disabled, don't OK the run uncontrollably!"
else
  filter="($2)&&status==PR"
fi
count=$3

if [ -d "$contestId" ]; then
  rm -r $contestId
fi

echo -e "${Cyan}INFO${NC}       Filtering runs from contest [$contestId] with filter [$filter] and limit [$count]"
echo -e "${Cyan}INFO${NC}       Reviewing ${IYellow}[$(boban -i "$contestId" -f "$filter" -c "$count" -d . | wc -l | xargs)]${NC} filtered runs"

cleanup() {
  rm -r "$contestId"
}

for file in "$contestId"/*; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        runId="${filename%.*}"
        lang="${filename##*.}"

        current="./main.$lang"
        cp "$file" "$current"
        echo -e "${Cyan}INFO${NC}       Review $file => $current"

        runLinesCount=$(cat "$current" | awk '/=============== COMMENTS/{print NR-2; found=1; exit} END{if(!found) print NR}') # количестов строк кода в посылке до секции комментов

        while [ 1 ]; do
            read -r -p "$(echo -e "${Green}INTERACT${NC}")   Type the resolution ($(echo -e "${Green}OK(щл)${NC}/${Red}RJ(ко)${NC}/skip/rejudge")): " verdict
            verdict=$(echo "$verdict" | tr '[:lower:]' '[:upper:]') # Преобразование в верхний регистр

            comment=$(head -n "$runLinesCount" "$current" | nl -w4 -s"] " | sed -e 's/ \([0-9]\)/[\1/' | grep -C 1 "///")
            if [[ -n $comment ]]; then
              comment=$'--\n'"$comment"$'\n--'
            fi

            if [ "$verdict" = "ЩЛ" ]; then
              verdict="OK"
            fi

            if [ "$verdict" = "КО" ]; then
              verdict="RJ"
            fi

            if [ "$verdict" = "OK" ] || [ "$verdict" = "RJ" ]; then
                echo $runId | ripper -i "$contestId" -s "$verdict" -c "$comment"
                break
            elif [ "$verdict" = "REJUDGE" ]; then
                echo $runId | ripper -i "$contestId" -s "rejudge"
                break
            elif [ "$verdict" == "SKIP" ]; then
                break
            elif [ -z "$verdict" ]; then
                cleanup
                exit 0
            fi
        done
    fi
done

cleanup