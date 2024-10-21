#!/bin/bash

error() {
    local msg="$1"
    if [ -z "$msg" ]; then
        echo "error: no message provided"
        exit 1
    fi
    echo "error: $msg"
    exit 1
}

validate_tag() {
    in=$1
    if [ -z "$in" ]; then
        error "empty input"
    fi

    regex="^v[0-9]\.[0-9]\.[0-9]$"
    if [[ ! $in =~ $regex ]]; then
        error "invalid input"
    fi
}

get_latest_release_tag() {
    latest_tag=$(curl -sL "https://github.com/maximekuhn/metego/releases/latest" | grep 'Raspberry Pi release' | tail -1 | sed -E 's/.*(v[0-9]+\.[0-9]+\.[0-9]+).*/\1/')
    echo "$latest_tag"
}

compare_versions() {
    local latest=$1
    local next=$2

    echo "comparing tags: latest ($latest), next ($next)"

    if [[ "$latest" == "$next" ]]; then
        error "tags are the same"
    elif [[ $(printf '%s\n' "$latest" "$next" | sort -V | head -n1) == "$latest" ]]; then
        echo "latest ($latest) is older than next ($next): OK"
    else
        error "latest ($latest) is newer than next ($next)"
    fi
}

validate_tag "$1"
latest_tag=$(get_latest_release_tag "$1")
echo "latest release tag: $latest_tag"
compare_versions "$latest_tag" "$1"
