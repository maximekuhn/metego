#!/bin/bash

print_err_and_exit_1() {
    local errMsg="$1"
    if [ -z "$errMsg" ]; then
        echo "error: no message and provided"
        exit 1
    fi
    echo "error: $errMsg"
    exit 1
}

get_latest_release_tag() {
    latest_tag=$(curl -sL "https://github.com/maximekuhn/metego/releases/latest" | grep 'Raspberry Pi release' | tail -1 | sed -E 's/.*(v[0-9]+\.[0-9]+\.[0-9]+).*/\1/')
    if [ -z "$latest_tag" ]; then
        print_err_and_exit_1 "could not get latest release tag"
    fi
    echo "$latest_tag"
}

download_assets() {
    release_tag=$(get_latest_release_tag)
    bin_download_url="https://github.com/maximekuhn/metego/releases/download/$release_tag/web"
    static_download_url="https://github.com/maximekuhn/metego/releases/download/$release_tag/static.zip"
    wget "$bin_download_url" -O web.new
    wget "$static_download_url" -O static.new.zip
}

unzip_static_files() {
    unzip static.new.zip -d static.new
}

arrange_static_files() {
    rm -rf static || "no previous static files found"
    mv static.new/static .
    rm -rf static.new
    rm static.new.zip
}

update_binary() {
    chmod +x web.new
    mv web.new web
}

restart_service() {
    sudo systemctl restart metego
}

download_assets
unzip_static_files
arrange_static_files
update_binary
restart_service

