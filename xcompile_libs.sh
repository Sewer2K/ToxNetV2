#!/bin/bash
# save as xcompile_libs.sh

mkdir -p cross_libs

# Mirror Base Paths
BASE_MIRROR="https://reflection.grit.ucsb.edu/debian/debian/pool/main"
SODIUM_DIR="$BASE_MIRROR/libs/libsodium"
TOXCORE_DIR="$BASE_MIRROR/libt/libtoxcore"

# Architectures to download
ARCHS=("amd64" "arm64" "armel" "armhf" "i386" "mips64el" "mipsel" "ppc64el" "s390x")

download_and_extract() {
    local arch=$1
    local pkg_name=$2
    local url=$3
    
    echo "[+] Downloading $pkg_name for $arch..."
    
    # Download the file
    if wget -q --show-progress "$url" -O "/tmp/${pkg_name}_${arch}.deb"; then
        # Check if it's a valid Debian archive
        if file "/tmp/${pkg_name}_${arch}.deb" | grep -q "Debian binary package"; then
            dpkg -x "/tmp/${pkg_name}_${arch}.deb" "cross_libs/${arch}"
            echo "[✓] Extracted $pkg_name for $arch"
        else
            echo "[!] ERROR: Downloaded file for $pkg_name ($arch) is invalid (check URL)."
        fi
    else
        echo "[!] ERROR: Could not download $pkg_name for $arch from $url"
    fi
}

for arch in "${ARCHS[@]}"; do
    # 1. Download libsodium-dev
    SODIUM_URL="${SODIUM_DIR}/libsodium-dev_1.0.18-1_${arch}.deb"
    download_and_extract "$arch" "libsodium" "$SODIUM_URL"

    # 2. Download libtoxcore-dev
    TOXCORE_URL="${TOXCORE_DIR}/libtoxcore-dev_0.2.18-1_${arch}.deb"
    download_and_extract "$arch" "toxcore" "$TOXCORE_URL"
    
    echo "-------------------------------------------"
done

echo "[+] All libraries processed in ./cross_libs/"
