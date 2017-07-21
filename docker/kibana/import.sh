#!/bin/sh
set -e

##
## Configuration
##

IMPORT_NAME="MgoExporter"
ES_HOST="http://127.0.0.1:9200"
KIBANA_INDEX=".kibana"
INDEX_PATTERN="mgo-exporter-*"
MAKE_DEFAULT="false"

##
## Internal Configuration
##

TIME_FIELD="time"
FIELDS_FILE="$(dirname $0)/index-pattern/mgo-exporter.json"
FORMAT_FILE="$(dirname $0)/field-format/mgo-exporter.json"
VISU_FOLDER="$(dirname $0)/visualization"
DASH_FOLDER="$(dirname $0)/dashboard"

##
## Global variables
##

IMPORT_ID=
SED_INDEX_PATTERN=
SED_VISU_ID=
SED_VISU_TITLE=

##
## Import functions
##

createIndexPattern () {

    fields="[]"
    if [ -n "${FIELDS_FILE}" ]; then
        fields="$(cat ${FIELDS_FILE})"
    fi

    field_fmt="{}"
    if [ -n "${FORMAT_FILE}" ]; then
        field_fmt="$(cat ${FORMAT_FILE})"
    fi

    req_data="$(jq -ncM '.title = $idx | .timeFieldName = $tf | .fields = $fields | .fieldFormatMap = $fmt' \
        --arg idx "${INDEX_PATTERN}" \
        --arg tf "${TIME_FIELD}" \
        --arg fmt "${field_fmt}" \
        --arg fields  "${fields}")"

    echo "> Creating index-pattern for ${INDEX_PATTERN}" 1>&2
    curl -sS -X POST "${ES_HOST}/${KIBANA_INDEX}/index-pattern/${INDEX_PATTERN}?pretty" -d "${req_data}"

    if [ "$MAKE_DEFAULT" = "true" ]; then

        echo "> Making ${INDEX_PATTERN} Kibana's default index" 1>&2
        curl -sS -X POST "${ES_HOST}/${KIBANA_INDEX}/config/_update_by_query?pretty" -d '{
          "script": {
            "inline": "ctx._source.defaultIndex = \"'${INDEX_PATTERN}'\"",
            "lang": "painless"
          }
        }'

    fi
}

importVisualization () {

    for v in ${VISU_FOLDER}/*.json; do
        vis_name="mgo-${IMPORT_ID}-$(basename "${v}" '.json')"
        vis_data="$(cat "${v}" | sed -E "${SED_INDEX_PATTERN};${SED_VISU_TITLE}")"
        echo "> Importing visualization ${vis_name}" 1>&2
        curl -sS -X POST "${ES_HOST}/${KIBANA_INDEX}/visualization/${vis_name}?pretty" -d "${vis_data}"
    done
}

importDashboard () {

    for d in ${DASH_FOLDER}/*; do
        dash_name="mgo-${IMPORT_ID}-$(basename "${d}" '.json')"
        dash_data="$(cat "${d}" | sed -E "${SED_VISU_ID};${SED_VISU_TITLE}")"
        echo "> Importing dashboard ${dash_name}" 1>&2
        curl -sS -X POST "${ES_HOST}/${KIBANA_INDEX}/dashboard/${dash_name}?pretty" -d "${dash_data}"
    done
}

##
## Main
##

init_global () {
    IMPORT_ID=$(echo "${IMPORT_NAME}" | tr ' ' '-' | tr 'A-Z' 'a-z')
    SED_INDEX_PATTERN="s/mgo-exporter-\*/${INDEX_PATTERN}/g"
    SED_VISU_ID="s/(\\\\\"id\\\\\":\\\\\"mgo-)exporter(-[a-z-]+\\\\\")/\1${IMPORT_ID}\2/g"
    SED_VISU_TITLE="s/MgoExporter -/${IMPORT_NAME} -/g"
}

usage () {
    echo "Usage: $0 [-i [-m]] [-v] [-d] [-p INDEX_PATTERN] [-e ES_HOST] [-k KIBANA_INDEX] IMPORT_NAME" 1>&2
    echo 1>&2
    echo "Positional arguments:" 1>&2
    echo "  IMPORT_NAME            name to use to prefix visu titles and generate ids" 1>&2
    echo 1>&2
    echo "Options:" 1>&2
    echo "  -i                     create index pattern before importing visu" 1>&2
    echo "  -m                     make new index pattern the default one after creating it" 1>&2
    echo "  -v                     import visualizations" 1>&2
    echo "  -d                     import dashboards" 1>&2
    echo "  -p INDEX_PATTERN       use INDEX_PATTERN while importing [default: 'mgo-exporter-*']" 1>&2
    echo "  -e ES_HOST             Elastic search host to connect to while importing [default: 'http://127.0.0.1:9200']" 1>&2
    echo "  -k KIBANA_INDEX        Index to use while create index-pattern and visu [default: '.kibana']" 1>&2
    echo "  -h                     print this help and exit" 1>&2
}

main () {

    ## Variables
    impIdxP=0
    impVisu=0
    impDash=0

    ## Parse CLI arguments
    while getopts "himvdp:e:k:" opt; do
      case ${opt} in
        "h") usage; return 0 ;;
        "i") impIdxP=1 ;;
        "m") MAKE_DEFAULT="true" ;;
        "v") impVisu=1 ;;
        "d") impDash=1 ;;
        "p") INDEX_PATTERN="${OPTARG}" ;;
        "e") ES_HOST="${OPTARG}" ;;
        "k") KIBANA_INDEX="${OPTARG}"  ;;
        \?|:) echo "Try '$0 -h' for more information." 1>&2; exit 2 ;;
      esac
    done
    shift $((OPTIND-1))

    ## Get positional argument as IMPORT_NAME
    IMPORT_NAME="$*"

    ## Verify that we are asked to do something
    if [ "${impIdxP}${impVisu}${impDash}" = "000" ]; then
        usage
        return 0
    elif [ "${IMPORT_NAME}" = "" ]; then
        echo "$0: required argument not found -- IMPORT_NAME" 1>&2
    fi

    ## Generate import id and patterns
    init_global

    ## Run imports
    [ "${impIdxP}" -eq 1 ] && createIndexPattern
    [ "${impVisu}" -eq 1 ] && importVisualization
    [ "${impDash}" -eq 1 ] && importDashboard
}

main "$@"
