#!/bin/bash

sign() {
    if [ -z "$APPLE_DEVELOPER_ID" ]; then
        echo "skipping macOS code-signing; APPLE_DEVELOPER_ID not set" >&2
        return 0
    fi
    
    if [[ $1 == *.zip ]]; then
        xcrun notarytool submit "$1" --apple-id "${APPLE_ID?}" --team-id "${APPLE_DEVELOPER_ID?}" --password "${APPLE_ID_PASSWORD?}"
    else
        codesign --sign "Developer ID Installer: Jennifer Leigh Hamon (NM6K5DCNF3)" --timestamp --options runtime dist/dist/macos_darwin_amd64_v1/pinecone
        codesign --timestamp --options=runtime -s "${APPLE_DEVELOPER_ID?}" -v "$1"
    fi
}

sign $1