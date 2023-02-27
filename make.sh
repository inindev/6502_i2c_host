#!/bin/sh

set -e

#BOARD='qtpy-rp2040'
BOARD='feather-rp2040'
#BOARD='xiao-rp2040'
#BOARD='pico'
FWBIN='firmware.uf2'
VNAME='RPI-RP2'

# exit codes
#  1 - build failure (no .uf2 file)
#  2 - unknown os
#  3 - volume not found
#  4 - block device for volume not found (linux only)
#  5 - unknown error

flash_macosx() {
    uf2file="$1"

    vol="/Volumes/$VNAME"
    if [ ! -d "$vol" ]; then
#        echo "${red}error: $vol not found${rst}"
#        exit 3
        echo "mounting uf2 volume..."
        stty -f /dev/cu.usbmodem* 1200
        while [ ! -d "$vol" ]; do sleep 0.1; done
        sleep 0.5
    fi

    build_uf2 "$uf2file"

    cp -v "$uf2file" "$vol"
}

flash_linux() {
    uf2file="$1"

    vol="/dev/disk/by-label/$VNAME"
    if [ ! -d "$vol" ]; then
        echo "${red}error: $vol not found${rst}"
        exit 3
    fi

    if ! bd=$(readlink -f "$vol"); then
        echo "${red}error: block device for $vol not found${rst}"
        exit 4
    fi

    build_uf2 "$uf2file"

    td=$(mktemp -d)
    sudo mount "$bd" "$td"
    sudo cp -v "$uf2file" "$td"
    sudo sync
    sudo umount "$td"
    rm -rf "$td"
}

flash() {
    uf2file="$1"
    if [ -d '/Volumes' ]; then
        echo "${yel}mac osx uf2-flash${rst}"
        flash_macosx "$uf2file"
    elif [ -d '/dev/disk' ]; then
        echo "${yel}linux uf2-flash${rst}"
        flash_linux "$uf2file"
    else
        error "${red}error: system does not appear to be mac osx nor linux${rst}"
        exit 2
    fi
}

build_uf2() {
    uf2file="$1"

    rm -rf "$uf2file"
    if ! tinygo build -target "$BOARD" -o "$uf2file" .; then
        echo "${yel}uf2-flash: ${red}failed to build uf2 file${rst}"
        exit 1
    fi
}

main() {
    flash "$FWBIN"
}

rst='\033[m'
bld='\033[1m'
red='\033[31m'
grn='\033[32m'
yel='\033[33m'
blu='\033[34m'
mag='\033[35m'
cya='\033[36m'

main $@
