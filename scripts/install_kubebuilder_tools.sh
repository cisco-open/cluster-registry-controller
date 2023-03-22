#!/usr/bin/env bash

# Copyright (c) 2022 Cisco and/or its affiliates. All rights reserved.

set -euo pipefail

os=$(go env GOOS)
arch=$(go env GOARCH)

# These are the current official versions for kubebuilder-tools (v1.19.2 or v1.22.x is recommended)
  # Reference: https://book.kubebuilder.io/reference/envtest.html#kubernetes-120-and-121-binary-issues
verified_kubebuilder_tools_version=(1.19.2)


FLOCK=`which flock 2> /dev/null || true`
if [ -z "${FLOCK}" ] || [ ! -x "${FLOCK}" ]; then
  # On MacOS or any system that does not have the flock binary
  # let's not use any locking
  FLOCK="true"
fi

(
  $FLOCK -x 200

  for i in "$@"; do
    case $i in
    --kubebuilder-tools-version=?*)
      kubebuilder_tools_version=${i#*=}
      ;;
    *)
    esac
  done

  if [ -z "${kubebuilder_tools_version:-}" ]; then
    echo "--kubebuilder-tools-version flag missing, possible options: ${verified_kubebuilder_tools_version[*]}"
    exit 1
  fi

  target_dir=bin/kubebuilder-tools/${kubebuilder_tools_version}

  mkdir -p "${target_dir}"

  if [[ ! -e "${target_dir}/complete" ]]; then
    echo "Installing kubebuilder-tools version v${kubebuilder_tools_version}..."
    echo "Downloading https://go.kubebuilder.io/test-tools/${kubebuilder_tools_version}/${os}/${arch}"
    #  Validate kubebuilder-tools version
    if echo ${verified_kubebuilder_tools_version[@]} | grep -q -w ${kubebuilder_tools_version}; then
        curl -sSLo ${target_dir}/envtest-bins.tar.gz "https://go.kubebuilder.io/test-tools/${kubebuilder_tools_version}/${os}/${arch}"
        tar --strip-components=1 -zvxf ${target_dir}/envtest-bins.tar.gz -C "${target_dir}"
        rm ${target_dir}/envtest-bins.tar.gz
        touch "${target_dir}/complete"

        echo " The current kubebuilder-tools version is v${kubebuilder_tools_version}."
    else
        echo "The specified kubebuilder-tools version is not valid. Overriding failed."
    fi
  else
    echo "kubebuilder-tools already installed under ${target_dir}"
  fi

) 200> /tmp/install-kubebuilder-tools.lock