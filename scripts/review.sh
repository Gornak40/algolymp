#!/bin/bash

templateContestId=46500

export LC_ALL=en_US.UTF-8
export PATH="$(dirname $0)/../bin:$PATH"

NC='\033[0m'       # Text Reset
# Regular Colors
Red='\033[0;31m'          # Red
Green='\033[0;32m'        # Green
IYellow='\033[0;93m'
Cyan='\033[0;36m'         # Cyan

if [ "$#" -lt 2 ]; then
    echo "Usage: $0 <contestId> <filter> [<count>] [<showAll>]"
    echo
    echo "             Run review for pending review filtered runs in contest."
    echo
    echo "Arguments:"
    echo
    echo "  contestId    (int)        Contest id"
    echo "  filter       (string)     Filter expression to filter runs"
    echo "  count        (int)        Maximum count of runs to review (default: 10)"
    echo "  showAll      (bool)       Disable pending review runs filter"
    exit 1
fi

contestId=$1
filter=$2
if [ 0 -le $3 ]; then
  count=$3
else
  count=10
fi
showAll=$4

if [ $contestId -le 31 ]; then
  contestId=$((contestId+templateContestId))
fi

if [[ $filter =~ ^[A-Za-z]$ ]]; then
  filter="prob=='$filter'"
elif [[ $filter =~ ^[0-9]+$ ]]; then
  filter="prob=='$filter'"
  showAll='true'
fi

if [ "$4" != 'true' ]; then
  filter="($filter)&&status==PR"
fi

cleanup() {
  if [ -d "$contestId" ]; then
    rm -r $contestId
  fi
}

cleanup

echo -e "${Cyan}INFO${NC}       Filtering runs from contest [$contestId] with filter [$filter] and limit [$count]"
echo -e "${Cyan}INFO${NC}       Reviewing ${IYellow}[$(boban -i "$contestId" -f "$filter" -c "$count" -e -d . | wc -l | xargs)]${NC} filtered runs"

if [ "$4" == 'true' ]; then
  echo -e "${Red}WARNING${NC}    Pending review filter is disabled, be aware for OK wrong run!"
fi

for file in "$contestId"/*; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        runId="${filename%.*}"
        lang="${filename##*.}"

        current="./main.$lang"
        cp "$file" "$current"
        echo -e "${Cyan}INFO${NC}       Review $file => $current"

        runLinesCount=$(cat "$current" | awk '/COMMENTS FOR CURRENT RUN/{print NR-2; found=1; exit} END{if(!found) print NR}') # количестов строк кода в посылке до секции комментов

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