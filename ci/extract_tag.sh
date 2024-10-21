#!/bin/bash

# tag format: v[1-9].[1-9].[1-9]-rc[1-9][1-9]?
extract_tag() {
    in=$1
    if [ -z "$in" ]; then
        echo "error: empty input"
    fi

    regex="^(v[1-9]\.[1-9]\.[1-9])-rc[1-9][0-9]?$"
    if [[ $in =~ $regex ]]; then
        echo "${BASH_REMATCH[1]}"
    else
        echo "error: invalid format"
    fi
}

extract_tag "$1"
