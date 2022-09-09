#!/bin/sh
#
# Copyright (c) 2022 Contributors to the Eclipse Foundation
#
# See the NOTICE file(s) distributed with this work for additional
# information regarding copyright ownership.
#
# This program and the accompanying materials are made available under the
# terms of the Eclipse Public License 2.0 which is available at
# https://www.eclipse.org/legal/epl-2.0, or the Apache License, Version 2.0
# which is available at https://www.apache.org/licenses/LICENSE-2.0.
#
# SPDX-License-Identifier: EPL-2.0 OR Apache-2.0

CHANNEL="${CHANNEL:-stable}"
DOWNLOAD_URL="${DOWNLOAD_URL:-https://download.docker.com}"

command_exists() {
    command -v "$@" > /dev/null 2>&1
}

get_distribution() {
    lsb_dist=""
    # The officially support distros should have /etc/os-release
    if [ -r /etc/os-release ]; then
        lsb_dist="$(. /etc/os-release && echo "$ID")"
    fi
    echo "$lsb_dist"
}

get_distribution_version() {
    dist_version=""
    case "$lsb_dist" in

    ubuntu)
        if command_exists lsb_release; then
            dist_version="$(lsb_release --codename | cut -f2)"
        fi
        if [ -z "$dist_version" ] && [ -r /etc/lsb-release ]; then
            dist_version="$(. /etc/lsb-release && echo "$DISTRIB_CODENAME")"
        fi
    ;;

    debian|raspbian)
        dist_version="$(sed 's/\/.*//' /etc/debian_version | sed 's/\..*//')"
        case "$dist_version" in
            11)
                dist_version="bullseye"
            ;;
            10)
                dist_version="buster"
            ;;
            9)
                dist_version="stretch"
            ;;
            8)
                dist_version="jessie"
            ;;
            *)
                cat >&2 <<-'EOF'
            Error: This Linux distribution version is not supported by the script.
            Please, try installing the containerd package from you official distribution's repository.
EOF
            exit 1
            ;;
        esac
    ;;

    *)
                cat >&2 <<-'EOF'
            Error: Unsupported Linux distribution.
            The installation will be terminated
EOF
            exit 1
    ;;

    esac

    echo "$dist_version"
}

add_debian_backport_repo() {
    debian_version="$1"
    backports="deb http://ftp.debian.org/debian $debian_version-backports main"
    if ! grep -Fxq "$backports" /etc/apt/sources.list; then
        (set -x; $sh_c "echo \"$backports\" >> /etc/apt/sources.list")
    fi
}

install_containerd_io() {
    pre_reqs="apt-transport-https ca-certificates curl"
    if [ "$lsb_dist" = "debian" ]; then
        # libseccomp2 does not exist for debian jessie main repos for aarch64
        if [ "$(uname -m)" = "aarch64" ] && [ "$dist_version" = "jessie" ]; then
            add_debian_backport_repo "$dist_version"
        fi
    fi

    if ! command -v gpg > /dev/null; then
        pre_reqs="$pre_reqs gnupg"
    fi

    $sh_c 'apt-get update -qq'
    $sh_c "DEBIAN_FRONTEND=noninteractive apt-get install -y -qq $pre_reqs"
    if  ! grep -Fqr "$DOWNLOAD_URL" /etc/apt/sources.list.d; then
        apt_repo="deb [arch=$(dpkg --print-architecture)] $DOWNLOAD_URL/linux/$lsb_dist $dist_version $CHANNEL"
        echo "$apt_repo"
        $sh_c "curl -fsSL \"$DOWNLOAD_URL/linux/$lsb_dist/gpg\" | apt-key add -qq -"
        $sh_c "echo \"$apt_repo\" > /etc/apt/sources.list.d/containerd.list"
        $sh_c 'apt-get update -qq'
    fi
    $sh_c "apt-get install -y -qq --no-install-recommends containerd.io"
}

#-----------------------------------------------------------
install() {
    set -e

    user="$(id -un 2>/dev/null || true)"

    sh_c='sh -c'
    if [ "$user" != 'root' ]; then
        if command_exists sudo; then
                sh_c='sudo -E sh -c'
        elif command_exists su; then
                sh_c='su -c'
        else
                cat >&2 <<-'EOF'
            Error: this installer needs the ability to run commands as root.
            We are unable to find either "sudo" or "su" available to make this happen.
EOF
                exit 1
        fi
    fi

    set +e

    # Get distro
    lsb_dist=$( get_distribution )
    lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"

    # Get and check distro version
    dist_version=$( get_distribution_version ) || exit 1
    dist_version="$(echo "$dist_version" | tr '[:upper:]' '[:lower:]')"
    echo "Target installation $lsb_dist's version is $dist_version"

    # Install containerd.io package
    install_containerd_io || exit 1

    # Check installation
    if command_exists ctr && [ -e /var/run/containerd/containerd.sock ]; then
        (
            $sh_c 'ctr version'
        ) || exit 1
    fi
}

install
