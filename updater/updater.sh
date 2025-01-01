#!/bin/bash

# -e: exit immediatly when any command returns a non-zero exit code
# -u: exit if an unknown variable is used
# -x: prints all excuted commands
# -o pipefail: prevents error in a pipeline from being masked
set -euxo pipefail

# Where to log (relative)
LOG_FILE="metego_update_logs.log"

# log_err_and_exit_1 accepts an error message as first argument.
# It will be logged to the specified $LOG_FILE and the script execution
# will end with an exit code 1.
log_err_and_exit_1() {
    local errMsg="$1"
    if [ -z "$errMsg" ]; then
        log_msg "error: no message and provided"
        exit 1
    fi
    log_msg "error: $errMsg"
    exit 1
}

# log_msg appends a new log line to the $LOG_FILE.
# The msg is provided as first argument of this function.
# The date will be automatically added before the message.
log_msg() {
    local msg="$1"
    echo "[METEGO - UPDATE LOGS] $(date): $msg" >> $LOG_FILE 2>&1
}

# get_latest_release_tag fetches the latest release tag and return it.
# For example: v0.0.7
get_latest_release_tag() {
    log_msg "fetching latest release tag..."
    latest_tag=$(curl -sL "https://github.com/maximekuhn/metego/releases/latest" | grep 'Raspberry Pi release' | tail -1 | sed -E 's/.*(v[0-9]+\.[0-9]+\.[0-9]+).*/\1/')
    if [ -z "$latest_tag" ]; then
        log_err_and_exit_1 "could not get latest release tag"
    fi
    log_msg "latest tag fetched: $latest_tag"
    echo "$latest_tag"
}

download_assets() {
    log_msg "downloading assets..."
    release_tag=$(get_latest_release_tag)
    bin_download_url="https://github.com/maximekuhn/metego/releases/download/$release_tag/web"
    static_download_url="https://github.com/maximekuhn/metego/releases/download/$release_tag/static.zip"
    wget "$bin_download_url" -O web.new
    wget "$static_download_url" -O static.new.zip
    log_msg "assets downloaded!"
}

unzip_static_files() {
    log_msg "unzipping static files"
    unzip static.new.zip -d static.new
    log_msg "static files unzipped!"
}

arrange_static_files() {
    log_msg "arranging static files"
    rm -rf static || log_msg "no previous static files found"
    mv static.new/static .
    rm -rf static.new
    rm static.new.zip
    log_msg "static files arranged!"
}

update_binary() {
    log_msg "updating binary..."
    chmod +x web.new
    mv web.new web
    log_msg "binary updated!"
}

restart_service() {
    log_msg "restarting metego service..."
    sudo systemctl restart metego
    log_msg "metego service restarted!"
}

# update_update_script will fetch the latest update script from the
# repository's main branch and replace this one with it.
update_update_script() {
    log_msg "updating update script"
    # TODO
    log_msg "update script updated!"
}

## MAIN ##
log_msg "UPDATE SCRIPT START"
download_assets
unzip_static_files
arrange_static_files
update_binary
restart_service
update_update_script
log_msg "UPDATE SCRIPT END"

