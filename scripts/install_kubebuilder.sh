#!/usr/bin/env bash

# Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.

set -euo pipefail

install=YES
version=""
while [[ $# -gt 0 ]]; do
  case "$1" in
      -d|--download-only) install=NO; shift;;
      *) version="$1"; shift ;;
  esac
done

[ -z "${version}" ] && { echo "Usage: $0 [-d|--download-only] <version>"; exit 1; }

target_dir_name=kubebuilder-${version}
link_path=bin/kubebuilder

mkdir -p bin

if [ $install = YES ]; then
  [ -e ${link_path} ] && rm -r ${link_path}
  ln -s "${target_dir_name}" ${link_path}
fi

if [ ! -e bin/"${target_dir_name}" ]; then
    os=$(go env GOOS)
    arch=$(go env GOARCH)

    # download kubebuilder and extract it to tmp
    curl -L "https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${version}/kubebuilder_${version}_${os}_${arch}.tar.gz" | tar -xz -C /tmp/

    # extract the archive
    mv "/tmp/kubebuilder_${version}_${os}_${arch}" bin/"${target_dir_name}"
fi
