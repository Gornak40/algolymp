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

legend() {
  echo "Usage: $0 <filter> [<contestId>] [<count>] [<showAll>]"
  echo
  echo "             Run review for pending review filtered runs in contest."
  echo
  echo "Arguments:"
  echo
  echo "  contestId    (int)        Contest id"
  echo "  filter       (string)     Filter expression to filter runs"
  echo "  count        (int)        Maximum count of runs to review (default: 10)"
  echo "  showAll      (bool)       Disable pending review runs filter"
  echo
  echo "Description:"
  echo
  echo "  To get certain run by id instead of 'id==137' you can use just '137':"
  echo
  echo "            review.sh 137 46500"
  echo
  echo
  echo "  To get all runs for certain problem instead of 'prob==\"A\"' you can use just 'A':"
  echo
  echo "            review.sh A 46500"
  echo
  echo
  echo "  Instead of full 'contestId' argument you may specify 'templateContestId' in this script"
  echo "  and then use just number of the day, instead of 'xxx01' use '1':"
  echo
  echo "            review.sh A 1"
  echo
  echo
  echo "  After first specification of 'contestId' you may just omit it for further requests:"
  echo
  echo "            review.sh A"
  echo
  echo
  echo "Examples:"
  echo
  echo "  To get first 20 pending runs for problem 'C' with 'id % 6 = 1' in contest 46501:"
  echo
  echo "            review.sh \"prob=='A'&&id%6==1\" 46501 20"
}

if [ "$#" -lt 1 ]; then
    legend
    exit 1
fi

contestId=$2
filter=$1
if [ -z $3 ]; then
  count=10
else
  count=$3
fi
showAll=$4

if [ -z "$contestId" ]; then
  if ! [ -d ~/.config/ejudge ]; then
    mkdir ~/.config/ejudge
  fi
  contestId=$(cat ~/.config/ejudge/last_contest)
fi
if [ -z "$contestId" ]; then
  legend
  echo
  echo -e "${Green}Please specify contestId argument, then it will be saved for future requests!${NC}"
  exit 1
fi

if [ $contestId -le 31 ]; then
  contestId=$((contestId+templateContestId))
fi

echo $contestId > ~/.config/ejudge/last_contest

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