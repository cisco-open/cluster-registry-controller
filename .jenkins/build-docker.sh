#!/bin/bash

set -e  # errexit
set -u  # nounset
set -x

DOCKERFILE="./Dockerfile"
CONTEXT="."
BUILD_ARGS=( )

while [[ $# -gt 0 ]]; do
    case "$1" in
        -i|--image)
            IMG="$2"
            shift
            shift
            ;;
        -b|--build-id)
            BUILD_ID="$2"
            shift
            shift
            ;;
        -f|--dockerfile)
            DOCKERFILE="$2"
            shift
            shift
            ;;
        -c|--context)
            CONTEXT="$2"
            shift
            shift
            ;;
        --build-arg)
            BUILD_ARGS+=("$1 $2")
            shift
            shift
            ;;
        *)
            echo "Unknown argument $1"
            shift
    esac
done

# get artifactory creds
sed -i '1d' "${NYOTA_CREDENTIALS_FILE}"
. "${NYOTA_CREDENTIALS_FILE}"

docker build \
    -f "${DOCKERFILE}" \
    --build-arg GOPROXY="https://proxy.golang.org,https://${artifactory_user}:${artifactory_password}@${artifactory_url},direct" \
    ${BUILD_ARGS[@]} \
    "${CONTEXT}" \
    -t "${IMG}"

docker tag "${IMG}" "$(echo $IMG | sed 's?/?-?g')"
