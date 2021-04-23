#!/bin/bash

set -e  # errexit
set -u  # nounset

while [[ $# -gt 0 ]]; do
    case "$1" in
        -i|--image)
            IMG="$2"
            shift
            shift
            ;;
        *)
            echo "Unknown argument"
    esac
done

docker build . -t ${IMG}
