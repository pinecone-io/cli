#!/bin/sh
# Copyright (c) Pinecone Systems, Inc.
#
# Uninstall script for the Pinecone CLI (pc).
#
# This script removes the Pinecone CLI binary and optionally removes
# configuration and data files created by the CLI.
#
# Usage:
#   curl -fsSL https://pinecone.io/uninstall.sh | sh
#
# Or, if you have the script locally: 
#   cat uninstall.sh | sh 
#
# Environment variables:
#   PINECONE_INSTALL   Installation directory where pc was installed. Default: /usr/local/bin.
#   PINECONE_KEEP_CONFIG  Set to 1 to keep configuration files. Default: remove them.

set -eu

BINARY_NAME="pc"

main() {
    INSTALL_DIR="${PINECONE_INSTALL:-/usr/local/bin}"
    KEEP_CONFIG="${PINECONE_KEEP_CONFIG:-0}"

    BINARY_PATH="${INSTALL_DIR}/${BINARY_NAME}"

    # -------------------------------------------------------
    # Step 0: Check for package-manager installations
    # -------------------------------------------------------
    check_package_manager

    # -------------------------------------------------------
    # Step 1: Remove the binary
    # -------------------------------------------------------
    remove_binary

    # -------------------------------------------------------
    # Step 2: Remove configuration and data files
    # -------------------------------------------------------
    if [ "$KEEP_CONFIG" != "1" ]; then
        remove_config
    else
        log "Keeping configuration files (PINECONE_KEEP_CONFIG=1)."
    fi

    # -------------------------------------------------------
    # Done
    # -------------------------------------------------------
    log ""
    log "Pinecone CLI has been uninstalled."
}

# =========================================================
#  Helpers
# =========================================================

log() {
    printf '%s\n' "$@"
}

err() {
    printf 'Error: %s\n' "$@" >&2
    exit 1
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# =========================================================
#  Package-manager detection
# =========================================================

check_package_manager() {
    PC_PATH="$(command -v pc 2>/dev/null || true)"

    # Homebrew: check if pc resolves into a Homebrew prefix (Cellar or Caskroom)
    if [ -n "$PC_PATH" ] && command_exists brew; then
        BREW_PREFIX="$(brew --prefix 2>/dev/null || true)"
        REAL_PC="$(readlink "$PC_PATH" 2>/dev/null || true)"
        # readlink may return a relative path on macOS; resolve it to absolute
        case "$REAL_PC" in
            /*) ;;
            ?*) REAL_PC="$(cd "$(dirname "$PC_PATH")" && cd "$(dirname "$REAL_PC")" && pwd)/$(basename "$REAL_PC")" ;;
        esac
        case "$REAL_PC" in
            "${BREW_PREFIX}"/Caskroom/*)
                CASK_NAME="$(echo "$REAL_PC" | sed "s|${BREW_PREFIX}/Caskroom/||" | cut -d/ -f1)"
                log "Pinecone CLI appears to have been installed via Homebrew (cask: ${CASK_NAME})."
                log "You can uninstall it with:"
                log ""
                log "  brew uninstall --cask ${CASK_NAME}"
                log ""
                exit 0
                ;;
            "${BREW_PREFIX}"/Cellar/*)
                FORMULA_NAME="$(echo "$REAL_PC" | sed "s|${BREW_PREFIX}/Cellar/||" | cut -d/ -f1)"
                log "Pinecone CLI appears to have been installed via Homebrew (formula: ${FORMULA_NAME})."
                log "You can uninstall it with:"
                log ""
                log "  brew uninstall ${FORMULA_NAME}"
                log ""
                exit 0
                ;;
        esac
    fi

}

# =========================================================
#  Binary removal
# =========================================================

remove_binary() {
    if [ ! -f "$BINARY_PATH" ]; then
        log "Binary not found at ${BINARY_PATH}. It may have already been removed."
        return 0
    fi

    log "Removing ${BINARY_PATH}..."

    if [ -w "$INSTALL_DIR" ]; then
        rm -f "$BINARY_PATH"
    elif command_exists sudo; then
        log "Password may be required to remove ${BINARY_PATH}."
        sudo rm -f "$BINARY_PATH"
    elif command_exists doas; then
        log "Password may be required to remove ${BINARY_PATH}."
        doas rm -f "$BINARY_PATH"
    else
        err "Cannot write to ${INSTALL_DIR} and neither sudo nor doas are available." \
            "Either run this script as root or set PINECONE_INSTALL to the directory where pc was installed:" \
            "  curl -fsSL https://pinecone.io/uninstall.sh | PINECONE_INSTALL=\$HOME/.local/bin sh"
    fi

    log "Binary removed."
}

# =========================================================
#  Configuration removal
# =========================================================

remove_config() {
    CONFIG_DIR="${XDG_CONFIG_HOME:-${HOME}/.config}/pinecone"

    if [ -d "$CONFIG_DIR" ]; then
        log "Removing configuration directory ${CONFIG_DIR}..."
        rm -rf "$CONFIG_DIR"
        log "Configuration removed."
    else
        log "No configuration directory found."
    fi
}

main
