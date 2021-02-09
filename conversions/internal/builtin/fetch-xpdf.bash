#!/bin/bash -x

set -e

fetch_to() {
  local src=$1
  local dst=$2

  rm -rf "$dst"
  wget "$src" -O "${dst}.tmp.zip"
  unzip "${dst}.tmp.zip" -d "$dst"
  rm -f "${dst}.tmp.zip"
}

XPDF_VER=4.02
VER=0.1.3

MAC_LIB=https://github.com/ashutoshvarma/libxpdf/releases/download/v${VER}/libxpdf-${XPDF_VER}.macos-clang.x64.zip
LINUX_LIB=https://github.com/ashutoshvarma/libxpdf/releases/download/v${VER}/libxpdf-${XPDF_VER}.linux-gcc.x64.zip

fetch_to $MAC_LIB xpdf_darwin_amd64
fetch_to $LINUX_LIB xpdf_linux_amd64
