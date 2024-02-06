#!/usr/bin/bash
if [ "$#" -ne 2 ]; then
	echo "run with arguments: cid_1 cid_2"
	exit 1
fi
for i in $(seq $1 $2); do
	casper -i $i
done
