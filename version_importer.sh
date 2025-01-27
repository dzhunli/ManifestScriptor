#!/bin/bash
set -e

welcome_message() {
    echo "============================================"
    echo " Welcome to the Importer!"
    echo "============================================"
    echo "Usage:"
    echo "Run this script directly in your terminal:"
    echo "./script_name.sh"
    echo ""
    echo "Let's get started!"
    echo "============================================"
}
get_full_path() {
    local base_path="/var/www/clients/var/"
    echo -e "\033[1;33m Enter the path relative to $base_path: \033[0m"
    read -r relative_path
    if [[ -z "$relative_path" ]]; then
        echo -e "\033[1;31m Path cannot be empty. Exiting... \033[0m"
        exit 1
    fi
    local full_path="${base_path}${relative_path}"
    if [[ ! -d "$full_path" ]]; then
        echo -e "\033[1;31m Path does not exist: $full_path \033[0m"
        exit 1
    fi
    echo "$full_path"
}

check_os() {
    os_name=$(uname)
    if [[ "$os_name" != "Linux" ]]; then
        echo -e "\033[1;31m Error: Unsupported operating system: $os_name \033[0m" >&2
        exit 1
    fi
    echo -e -n "Operating system is"
    echo -e -n "\033[1;32m Linux \033[0m"
    echo -e "\033[32m\u2713\033[0m"
}
check_pv() {
    if ! command -v pv &> /dev/null; then
        legacy_view=true
    else
        echo -e -n "'pv' is already installed "
        echo -e "\033[32m\u2713\033[0m"

    fi
}
check_and_run_in_tmux() {
    if [[ -z "$TMUX" ]]; then
        echo -e "\033[1;31m Not running inside tmux session. Starting script in a new tmux session with logging... \033[0m"
        script_name=$(basename "$0")
        session_name="mysession_$(date +%s)"
        log_file="tmux_session_${session_name}.log"
        tmux new-session -d -s "$session_name" "tmux pipe-pane -o -t $session_name 'cat >> $log_file'; ./$script_name $@"
        echo -e "\033[1;32m Script has been started in tmux session: $session_name. Attaching... \033[0m"
        echo -e "\033[1;32m Log file: $log_file \033[0m"        
        tmux attach-session -t "$session_name"
        exit 0
    else
        echo -e "\033[1;32m Script is already running inside tmux session. \033[0m"
    fi
}


find_directories() {
    search_path=$1
    if [[ ! -d "$search_path" ]]; then
        echo -e "\033[1;31m The specified path '$search_path' does not exist. \033[0m" >&2
        exit 1
    fi
    versions_dir=$(find "$search_path" -type d -name "versions" -print -quit  )
    assets_dir=$(find "$search_path" -type d -name "assets" -print -quit )
    if [[ -n "$versions_dir" ]]; then
        echo -e -n "Found 'versions' directory at:"
        echo -e -n "\033[1;32m $versions_dir  \033[0m"
        measure_size "$versions_dir"
    else
        echo -e "\033[1;31m versions directory not found. \033[0m" >&2
        exit 1
    fi
    if [[ -n "$assets_dir" ]]; then
        echo -e -n "Found 'assets' directory at:"
        echo -e -n "\033[1;32m $assets_dir  \033[0m"
        measure_size "$assets_dir"
    else
        echo -e "\033[1;31m 'assets' directory not found. \033[0m" >&2
        exit 1
    fi
}
measure_size() {
    local dir=$1
    if [[ -n "$dir" ]]; then
        size=$(du -sh "$dir" | cut -f1)
        #echo -e -n  "Size of '$dir':"
        echo -e -n  "$size "
        echo -e "\033[32m\u2713\033[0m"
    else
        echo -e "\033[1;31m '$dir' is not a valid directory. \033[0m" >&2
        exit 1
    fi
}
create_archive_with_progress() {
    local dir=$1
    local archive_name=$2
    if [[ -d "$dir" ]]; then
        echo -e -n "Creating archive for directory"
        echo -e "\033[1;32m '$dir' \033[0m"
        if [[ legasy_view = true ]];then
            tar -czf "$archive_name" "$dir" --checkpoint=1000 --checkpoint-action=dot
        else
            tar -cf - "$dir" | pv -p -e -s $(du -sb "$dir" | awk '{print $1}') | gzip > "$archive_name"
        fi
        echo -e "\n"
        echo -e -n "Archive"
        echo -e -n "\033[1;32m '$archive_name' \033[0m"
        echo -e -n " created successfully. "
        measure_size $2
    else
        echo -e "\033[1;31m '$dir' is not a valid directory. \033[0m" >&2
        exit 1
    fi
}
main() {
    check_os
    check_pv
    check_and_run_in_tmux
    find_directories "./var/www/clients/var/"
    create_archive_with_progress "$assets_dir" "assets_archive.tar.gz"
    create_archive_with_progress "$versions_dir" "versions_archive.tar.gz"
}

main
