#!/bin/bash

# It\'s just an experiment. Don\'t worry :)

set -euo pipefail

# TODO: remove.
set -x

readonly POLYGON_URL="https://polygon.codeforces.com/api"

pcurl() {
	local method="$1"
	shift

	local params=("apiKey=${POLYGON_API_KEY}" "time=$(date +%s)" "$@")

	local sig_rand=$(head /dev/urandom | tr -dc '[:alnum:]' | head -c 6)
	local sig_params=$(printf "%s\n" "${params[@]}" | sort | paste -sd '&')
	local sig_data="${sig_rand}/${method}?${sig_params}#${POLYGON_SECRET}"
	local sig_hash=$(echo -n "$sig_data" | sha512sum | awk '{print $1}')

	local api_sig="${sig_rand}${sig_hash}"

	curl -fsSL "${POLYGON_URL}/${method}?${sig_params}&apiSig=${api_sig}"
}

pcurl problem.info problemId=436963
