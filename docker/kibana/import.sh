#!/bin/sh
set -e

##
## Configuration
##

IMPORT_NAME="MgoTest"
ES_HOST="http://elastic:9200"
KIBANA_INDEX=".kibana"
INDEX_PATTERN="mgo-exporter-*"

##
## Internal Configuration
##

TIME_FIELD="time"
FIELDS_FILE="index-pattern/mgo-exporter.json"
FORMAT_FILE="field-format/mgo-exporter.json"
MAKE_DEFAULT="true"
VISU_FOLDER="visualization"
DASH_FOLDER="dashboard"

##
## Global variables
##

IMPORT_ID=$(echo "${IMPORT_NAME}" | tr ' ' '-' | tr 'A-Z' 'a-z')
SED_INDEX_PATTERN="s/mgo-exporter-\*/${INDEX_PATTERN}/g"
SED_VISU_ID="s/(\\\\\"id\\\\\":\\\\\"mgo-)exporter(-[a-z-]+\\\\\")/\1${IMPORT_ID}\2/g"
SED_VISU_TITLE="s/MgoExporter -/${IMPORT_NAME} -/g"

importIndexPattern () {

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

    echo "> Creating index-pattern for ${INDEX_PATTERN}"
    curl -sS -X POST "${ES_HOST}/${KIBANA_INDEX}/index-pattern/${INDEX_PATTERN}?pretty" -d "${req_data}"

    if [ "$MAKE_DEFAULT" = "true" ]; then

        echo "> Making ${INDEX_PATTERN} Kibana's default index"
        curl -sS -X POST "${ES_HOST}/${KIBANA_INDEX}/config/_update_by_query?pretty" -d '{
          "script": {
            "inline": "ctx._source.defaultIndex = \"'${INDEX_PATTERN}'\"",
            "lang": "painless"
          }
        }'

    fi
}

importVisualization() {

    for v in ${VISU_FOLDER}/*.json; do
        vis_name="mgo-${IMPORT_ID}-$(basename "${v}" '.json')"
        vis_data="$(cat "${v}" | sed -E "${SED_INDEX_PATTERN};${SED_VISU_TITLE}")"
        echo "> Importing visualization ${vis_name}"
        curl -sS -X POST "${ES_HOST}/${KIBANA_INDEX}/visualization/${vis_name}?pretty" -d "${vis_data}"
    done
}

importDashboard () {

    for d in ${DASH_FOLDER}/*; do
        dash_name="mgo-${IMPORT_ID}-$(basename "${d}" '.json')"
        dash_data="$(cat "${d}" | sed -E "${SED_VISU_ID};${SED_VISU_TITLE}")"
        echo "> Importing dashboard ${dash_name}"
        curl -sS -X POST "${ES_HOST}/${KIBANA_INDEX}/dashboard/${dash_name}?pretty" -d "${dash_data}"
    done
}

importIndexPattern
importVisualization
importDashboard
