#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

if [ -z "$1" ]
  then
    echo "No setlist provided"
fi

BAND=$1
SETLIST=$2

echo "tell application \"Mail\"
    activate

    set MyEmail to make new outgoing message with properties {visible:true, subject:\"${BAND} Spick\", content:\"Generated at:\"}
    tell MyEmail
        make new to recipient at end of to recipients with properties {address:\"${EMAIL}\"}
        make new attachment with properties {file name:((\"${PROJECT_DIR}/out/${BAND}@${SETLIST}.pdf\" as POSIX file) as alias)}
        delay 1
        send MyEmail
    end tell

    delay 5
    quit
end tell
" | osascript
