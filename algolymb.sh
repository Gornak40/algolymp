#!/bin/bash

: ${POLYGON_URL:="https://polygon.codeforces.com"}

pcurl() {
	[[ $# -eq 0 ]] && echo "Usage: $0 METHOD [PARAM=VALUE]... [PARAM=@FILE]..." >&2 && return 2

	: ${POLYGON_API_KEY:?"env required"}
	: ${POLYGON_SECRET:?"env required"}

	local method="$1"
	shift

	local curl_params=("apiKey=${POLYGON_API_KEY}" "time=$(date +%s)" "$@")
	local sig_params=()
	for param in "${curl_params[@]}"; do
		local key= val=
		IFS='=' read -r key val <<< "${param}"
		[[ "${val}" == @* ]] && val=$(<"${val#@}")
		sig_params+=("${key}=${val}")
	done

	local sig_rand=$(head /dev/urandom | LC_ALL=C tr -dc '[:alnum:]' | head -c 6)
	local sig_string=$(printf "%s\n" "${sig_params[@]}" | sort | paste -sd '&')
	local sig_data="${sig_rand}/${method}?${sig_string}#${POLYGON_SECRET}"
	local sig_hash=$(echo -n "${sig_data}" | sha512sum | awk '{print $1}')

	curl_params+=("apiSig=${sig_rand}${sig_hash}")
	local response=$(curl -sSL -w "\n%{http_code}" "${curl_params[@]/#/-F}" "${POLYGON_URL}/api/${method}")

	local http_code=$(tail -n 1 <<< "${response}")
	local body=$(head -n -1 <<< "${response}")
	case "${http_code}" in
		200)
			jq -r ".result // .status" <<< "${body}" ;;
		400)
			jq -r ".comment" <<< "${body}"; return 1 ;;
		*)
			printf "%s" "${body}"; return 2 ;;
	esac
}

algp() {
	pcurl "problem.$1" problemId="$2" "${@:3}"
}
