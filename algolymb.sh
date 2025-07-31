#!/bin/bash

# It\'s just an experiment. Don\'t worry :)

set -euo pipefail

# TODO: remove.
set -x

################################################################################

readonly POLYGON_URL="https://polygon.codeforces.com/api"

pcurl() {
	local method="$1"
	shift

	local curl_params=("apiKey=${POLYGON_API_KEY}" "time=$(date +%s)" "$@")
	local sig_params=()
	for param in "${curl_params[@]}"; do
		local key val
		IFS='=' read -r key val <<< "${param}"
		[[ "${val}" == @* ]] && val=$(<"${val#@}")
		sig_params+=("${key}=${val}")
	done

	local sig_rand=$(head /dev/urandom | LC_ALL=C tr -dc '[:alnum:]' | head -c 6)
	local sig_string=$(printf "%s\n" "${sig_params[@]}" | sort | paste -sd '&')
	local sig_data="${sig_rand}/${method}?${sig_string}#${POLYGON_SECRET}"
	local sig_hash=$(echo -n "${sig_data}" | sha512sum | awk '{print $1}')

	curl_params+=("apiSig=${sig_rand}${sig_hash}")
	curl -fsSL "${curl_params[@]/#/-F}" "${POLYGON_URL}/${method}"
}

################################################################################

algp() {
	pcurl "problem.$1" problemId="$2" "${@:3}"
}

################################################################################

echo -n 436963 > /tmp/test.txt
algp info @/tmp/test.txt

algp statements 436963

algp validator 436963
