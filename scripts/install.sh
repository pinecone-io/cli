#!/bin/sh
# Copyright (c) Pinecone Systems, Inc.
#
# Install script for the Pinecone CLI (pc).
#
# This script detects your operating system and architecture, downloads the
# latest release of the Pinecone CLI from GitHub, verifies the checksum, and
# installs it to /usr/local/bin (or a directory of your choice).
#
# Usage:
#   curl -fsSL https://pinecone.io/install.sh | sh
#
# Environment variables:
#   PINECONE_VERSION   Pin to a specific version (e.g. "0.4.2"). Default: latest.
#   PINECONE_INSTALL   Installation directory. Default: /usr/local/bin.
#   PINECONE_NO_VERIFY Set to 1 to skip checksum verification.

set -eu

GITHUB_REPO="pinecone-io/cli"
BINARY_NAME="pc"

# All the code is wrapped in a main function that gets called at the
# bottom of the file, so that a truncated partial download doesn't end
# up executing half a script.
main() {
    INSTALL_DIR="${PINECONE_INSTALL:-/usr/local/bin}"
    VERSION="${PINECONE_VERSION:-}"
    NO_VERIFY="${PINECONE_NO_VERIFY:-0}"

    need_cmd uname

    # -------------------------------------------------------
    # Step 1: Detect OS and architecture
    # -------------------------------------------------------
    OS=""
    ARCH=""
    FILENAME=""

    detect_platform

    # -------------------------------------------------------
    # Step 2: Find an HTTP download tool (curl or wget)
    # -------------------------------------------------------
    HTTP=""
    if command_exists curl; then
        HTTP="curl"
    elif command_exists wget; then
        HTTP="wget"
    else
        err "Either curl or wget is required to download files. Please install one and try again."
    fi

    # -------------------------------------------------------
    # Step 3: Resolve version
    # -------------------------------------------------------
    if [ -z "$VERSION" ]; then
        log "Fetching latest release version..."
        VERSION=$(get_latest_version)
        if [ -z "$VERSION" ]; then
            err "Could not determine latest version. Set PINECONE_VERSION and try again."
        fi
    fi

    log "Installing Pinecone CLI v${VERSION} (${OS}/${ARCH})"

    # -------------------------------------------------------
    # Step 4: Build download URLs
    # -------------------------------------------------------
    BASE_URL="https://github.com/${GITHUB_REPO}/releases/download/v${VERSION}"
    ARCHIVE_URL="${BASE_URL}/${FILENAME}"
    CHECKSUMS_URL="${BASE_URL}/pc_${VERSION}_checksums.txt"

    # -------------------------------------------------------
    # Step 5: Download and extract to a temp directory
    # -------------------------------------------------------
    TMPDIR_ROOT="${TMPDIR:-/tmp}"
    WORK_DIR=$(mktemp -d "${TMPDIR_ROOT}/pinecone-cli.XXXXXX")
    trap 'rm -rf "$WORK_DIR"' EXIT

    log "Downloading ${ARCHIVE_URL}..."
    download "$ARCHIVE_URL" "${WORK_DIR}/${FILENAME}"

    # -------------------------------------------------------
    # Step 6: Verify checksum (unless opted out)
    # -------------------------------------------------------
    if [ "$NO_VERIFY" != "1" ]; then
        verify_checksum
    else
        log "Skipping checksum verification (PINECONE_NO_VERIFY=1)"
    fi

    # -------------------------------------------------------
    # Step 7: Extract the binary
    # -------------------------------------------------------
    log "Extracting..."
    tar -xzf "${WORK_DIR}/${FILENAME}" -C "$WORK_DIR"

    if [ ! -f "${WORK_DIR}/${BINARY_NAME}" ]; then
        err "Archive did not contain expected binary '${BINARY_NAME}'. Contents of archive:" \
            "$(ls -la "$WORK_DIR")"
    fi

    chmod +x "${WORK_DIR}/${BINARY_NAME}"

    # -------------------------------------------------------
    # Step 8: Install the binary
    # -------------------------------------------------------
    install_binary

    # -------------------------------------------------------
    # Done
    # -------------------------------------------------------
    log ""
    log "Pinecone CLI v${VERSION} installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"
    log ""
    log "Run 'pc --help' to get started."
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

need_cmd() {
    if ! command_exists "$1"; then
        err "'$1' is required but not found on your system."
    fi
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# =========================================================
#  Platform detection
# =========================================================

detect_platform() {
    local uname_os uname_arch

    uname_os="$(uname -s)"
    uname_arch="$(uname -m)"

    case "$uname_os" in
        Darwin)
            OS="Darwin"
            ARCH="all"
            FILENAME="pc_Darwin_all.tar.gz"
            ;;
        Linux)
            OS="Linux"
            case "$uname_arch" in
                x86_64|amd64)
                    ARCH="x86_64"
                    FILENAME="pc_Linux_x86_64.tar.gz"
                    ;;
                aarch64|arm64)
                    ARCH="arm64"
                    FILENAME="pc_Linux_arm64.tar.gz"
                    ;;
                i386|i686)
                    ARCH="i386"
                    FILENAME="pc_Linux_i386.tar.gz"
                    ;;
                *)
                    err "Unsupported Linux architecture: ${uname_arch}" \
                        "Supported architectures: x86_64, arm64, i386"
                    ;;
            esac
            ;;
        *)
            err "Unsupported operating system: ${uname_os}" \
                "This installer supports macOS (Darwin) and Linux." \
                "For Windows, download the .zip from https://github.com/${GITHUB_REPO}/releases"
            ;;
    esac
}

# =========================================================
#  HTTP helpers
# =========================================================

download() {
    local url="$1"
    local dest="$2"

    case "$HTTP" in
        curl)
            curl -fSL --progress-bar -o "$dest" "$url" || err "Failed to download: ${url}"
            ;;
        wget)
            wget -q --show-progress -O "$dest" "$url" || err "Failed to download: ${url}"
            ;;
    esac
}

fetch() {
    # Fetch URL contents to stdout (silent)
    local url="$1"
    case "$HTTP" in
        curl) curl -fsSL "$url" ;;
        wget) wget -qO- "$url" ;;
    esac
}

get_latest_version() {
    # Use the GitHub API to get the latest release tag.
    # Falls back to following the /releases/latest redirect if the API is rate-limited.
    local tag

    tag=$(fetch "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" 2>/dev/null \
        | grep '"tag_name"' \
        | sed -E 's/.*"tag_name":[[:space:]]*"v?([^"]+)".*/\1/' ) || true

    if [ -z "$tag" ] && [ "$HTTP" = "curl" ]; then
        # Fallback: follow the redirect and parse the URL
        tag=$(curl -fsSL -o /dev/null -w '%{url_effective}' \
            "https://github.com/${GITHUB_REPO}/releases/latest" 2>/dev/null \
            | sed 's|.*/v||') || true
    fi

    printf '%s' "$tag"
}

# =========================================================
#  Checksum verification
# =========================================================

verify_checksum() {
    local sha_cmd=""

    if command_exists sha256sum; then
        sha_cmd="sha256sum"
    elif command_exists shasum; then
        sha_cmd="shasum -a 256"
    else
        log "Warning: Neither sha256sum nor shasum found. Skipping checksum verification."
        return 0
    fi

    log "Verifying checksum..."
    download "$CHECKSUMS_URL" "${WORK_DIR}/checksums.txt"

    local expected actual
    expected=$(grep "${FILENAME}" "${WORK_DIR}/checksums.txt" | awk '{print $1}')

    if [ -z "$expected" ]; then
        err "Could not find checksum for ${FILENAME} in checksums file."
    fi

    actual=$($sha_cmd "${WORK_DIR}/${FILENAME}" | awk '{print $1}')

    if [ "$expected" != "$actual" ]; then
        err "Checksum verification failed!" \
            "Expected: ${expected}" \
            "Actual:   ${actual}" \
            "The downloaded file may be corrupted. Please try again."
    fi

    log "Checksum verified."
}

# =========================================================
#  Installation
# =========================================================

install_binary() {
    # Create install directory if it doesn't exist, then move binary into it.
    if [ -d "$INSTALL_DIR" ] && [ -w "$INSTALL_DIR" ]; then
        mv "${WORK_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    elif command_exists sudo; then
        log "Password may be required to install to ${INSTALL_DIR}."
        sudo sh -c "mkdir -p \"$INSTALL_DIR\" && mv \"${WORK_DIR}/${BINARY_NAME}\" \"${INSTALL_DIR}/${BINARY_NAME}\""
    elif command_exists doas; then
        log "Password may be required to install to ${INSTALL_DIR}."
        doas sh -c "mkdir -p \"$INSTALL_DIR\" && mv \"${WORK_DIR}/${BINARY_NAME}\" \"${INSTALL_DIR}/${BINARY_NAME}\""
    elif [ ! -d "$INSTALL_DIR" ]; then
        # No sudo/doas but directory doesn't exist — try to create it (works if parent is writable)
        mkdir -p "$INSTALL_DIR" 2>/dev/null || \
            err "Cannot create ${INSTALL_DIR} and neither sudo nor doas are available." \
                "Either run this script as root or set PINECONE_INSTALL to a writable directory:" \
                "  curl -fsSL https://pinecone.io/install.sh | PINECONE_INSTALL=\$HOME/.local/bin sh"
        mv "${WORK_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        err "Cannot write to ${INSTALL_DIR} and neither sudo nor doas are available." \
            "Either run this script as root or set PINECONE_INSTALL to a writable directory:" \
            "  curl -fsSL https://pinecone.io/install.sh | PINECONE_INSTALL=\$HOME/.local/bin sh"
    fi

    # Ensure the install directory is in PATH
    case ":${PATH}:" in
        *":${INSTALL_DIR}:"*) ;;
        *)
            log ""
            log "NOTE: ${INSTALL_DIR} is not in your \$PATH."
            log "Add it by running one of the following:"
            log ""
            log "  # For bash"
            log "  echo 'export PATH=\"${INSTALL_DIR}:\$PATH\"' >> ~/.bashrc && source ~/.bashrc"
            log ""
            log "  # For zsh"
            log "  echo 'export PATH=\"${INSTALL_DIR}:\$PATH\"' >> ~/.zshrc && source ~/.zshrc"
            log ""
            ;;
    esac
}

main
