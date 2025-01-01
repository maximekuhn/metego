#!/bin/bash

# -e: exit immediatly when any command returns a non-zero exit code
# -u: exit if an unknown variable is used
# -x: prints all excuted commands
# -o pipefail: prevents error in a pipeline from being masked
set -euxo pipefail

# Where to log (relative)
LOG_FILE="metego_update_logs.log"

# Current version
CURR_VERSION="curr_metego_version"

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

# check_if_update_is_needed true" if provided release tag ($1) is
# newer than the current running one.
check_if_update_is_needed() {
    log_msg "checking if update is needed..."

    local latestReleaseTag="$1"
    local currVersion=$(echo $CURR_VERSION)
    log_msg "latest release tag: $latestReleaseTag, current version $currVersion"

    # We assume nobody is deleting any release from the Github repository.
    # We can simpy check if latestReleaseTag is the same as the one stored
    # in the CURR_VERSION file. If it's not the same, then we need to update,
    # otherwise the software is already up to date.
    if [ "$latestReleaseTag" = "$currVersion" ]; then
        echo "false"
        return
    fi
    echo "true"
}

# download_assets accepts the release tag as its first argument
# and download the associated assets.
download_assets() {
    log_msg "downloading assets..."
    local release_tag="$1"
    if [ -z "$release_tag" ]; then
        log_err_and_exit_1 "missing release tag(or empty)"
    fi
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
    wget https://github.com/maximekuhn/metego/tree/main/updater -O latest_update_script.sh
    chmod +x latest_update_script.sh
    mv latest_update_script.sh updater.sh
    log_msg "update script updated!"
}

# store_new_version writes the new version to the $CURR_VERSION file.
# It accepts the release tag as the first argument.
store_new_version() {
    log_msg "storing new version..."
    local newVersion="$1"
    if [ -z "$newVersion" ]; then
        log_err_and_exit_1 "missing new version (or empty)"
    fi
    echo "$newVersion" > $CURR_VERSION
    log_msg "new version stored!"
}

## MAIN ##
log_msg "UPDATE SCRIPT START"
latest_release_tag=$(get_latest_release_tag)

update_needed=$(check_if_update_is_needed)
if [ "$update_needed" != "true" ]; then
    log_msg "metego is already up to date"
    log_msg "UPDATE SCRIPT END"
    exit 0
fi
log_msg "a new version has been found, metego needs to be updated with $latest_release_tag"

download_assets "$latest_release_tag"
unzip_static_files
arrange_static_files
update_binary
restart_service
update_update_script
store_new_version "$latest_release_tag"
log_msg "UPDATE SCRIPT END"

