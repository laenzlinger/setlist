#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

BAND=$1
SETLIST=$2

function text2pdf() {
    local id=$1
    local text=$2

    if [ "${BAND}" == "Howlers" ]
    then
 
		cat <<- EOF > ${PROJECT_DIR}/out/txt/reminder.txt
			%%MediaBox 0 0 842 595
			%%Font TmRm Times-Roman
			BT /TmRm 40 Tf 50 530 Td (${text}) Tj ET
		EOF

    else

		cat <<- EOF > ${PROJECT_DIR}/out/txt/reminder.txt
			%%MediaBox 0 0 595 842
			%%Font TmRm Times-Roman
			BT /TmRm 12 Tf 50 760 Td (${text}) Tj ET
		EOF

    fi

    mutool create -o "${PROJECT_DIR}/out/setlist/$id reminder.pdf" ${PROJECT_DIR}/out/txt/reminder.txt
}

if [ -z "$1" ]
then
    echo "No band provided"
    exit 1
fi

if [ -z "$2" ]
then
    echo "No setlist provided"
    exit 1
fi

mkdir -p ${PROJECT_DIR}/out/setlist ${PROJECT_DIR}/out/txt

i=1
while IFS=";" read name 
do
    name="${name:2}"
    printf -v id "%03d" $i
    echo "$id -> $name"
    if [ -f "${PROJECT_DIR}/${BAND}/songs/${name}.odt" ]; then
        ${SCRIPT_DIR}/pages2pdf.sh "$name" "${BAND}"
        cp "${PROJECT_DIR}/${BAND}/songs/${name}.pdf" "${PROJECT_DIR}/out/setlist/${id} ${name}.pdf"
    else
        text2pdf "$id" "$name"
    fi
    let i++
done < "${PROJECT_DIR}/${BAND}/gigs/${SETLIST}.md"

TITLE=${BAND}@${SETLIST}
OUTPUTFILE=${PROJECT_DIR}/out/${TITLE}.pdf

mutool merge -O garbage=deduplicate -o "${OUTPUTFILE}" ${PROJECT_DIR}/out/setlist/*.pdf
exiftool -Title="${TITLE}" -Author="$(date)" -Subject="${BAND}" -q -overwrite_original_in_place "${OUTPUTFILE}"
