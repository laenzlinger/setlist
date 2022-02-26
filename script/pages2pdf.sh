#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

fileName=$1
band=$2
sourceFile=${PROJECT_DIR}/${band}/songs/${fileName}.odt
targetDir=${PROJECT_DIR}/${band}/songs
targetFile=${targetDir}/${fileName}.pdf

if [ -e "${targetFile}" ] && [ "${targetFile}" -nt "${sourceFile}" ]
then
    # echo ${targetFile} is up to date
    exit 0
fi

libreoffice --headless --convert-to pdf --outdir ${targetDir} "${sourceFile}"
if [ $? != 0 ]
then
    exit -2
fi
exiftool -Subject="${band}" -q -overwrite_original_in_place "${targetFile}"
