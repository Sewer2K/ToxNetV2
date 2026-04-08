#!/bin/bash
# save as build_mips_libs.sh

BASE_DIR="/home/vboxuser/Desktop/toxnet"
LIBS_DIR="$BASE_DIR/cross_libs/mips/usr"
mkdir -p "$LIBS_DIR/lib/mips-linux-gnu"

cd "$BASE_DIR" || exit 1

# 1. Download Toxcore Source (Corrected Link)
if [ ! -d "toxcore-0.2.18" ]; then
    echo "[+] Downloading Toxcore 0.2.18 source..."
    # Using the direct source code archive from GitHub
    wget -q "https://github.com" -O toxcore_v0.2.18.tar.gz
    
    if file toxcore_v0.2.18.tar.gz | grep -q "gzip compressed data"; then
        tar -xf toxcore_v0.2.18.tar.gz
        mv c-toxcore-0.2.18 toxcore-0.2.18
        echo "[✓] Toxcore source extracted successfully."
    else
        echo "[!] ERROR: Download failed or file is not a valid gzip. Deleting broken file."
        rm -f toxcore_v0.2.18.tar.gz
        exit 1
    fi
fi

# 2. Build Libsodium (Skips if already present)
if [ ! -f "$LIBS_DIR/lib/mips-linux-gnu/libsodium.a" ]; then
    echo "[+] Building libsodium for MIPS..."
    cd "$BASE_DIR/libsodium-stable" || { echo "libsodium-stable folder missing!"; exit 1; }
    make clean 2>/dev/null
    ./configure --host=mips-linux-gnu \
                --prefix="$LIBS_DIR" \
                --libdir="$LIBS_DIR/lib/mips-linux-gnu" \
                --enable-static --disable-shared
    make -j$(nproc) && make install
    cd "$BASE_DIR"
fi

# 3. Build Toxcore
echo "[+] Building toxcore for MIPS..."
rm -rf toxcore_build_mips && mkdir -p toxcore_build_mips && cd toxcore_build_mips

cmake "../toxcore-0.2.18" \
  -DCMAKE_SYSTEM_NAME=Linux \
  -DCMAKE_C_COMPILER=mips-linux-gnu-gcc \
  -DCMAKE_CXX_COMPILER=mips-linux-gnu-g++ \
  -DSODIUM_LIBRARY="$LIBS_DIR/lib/mips-linux-gnu/libsodium.a" \
  -DSODIUM_INCLUDE_DIR="$LIBS_DIR/include" \
  -DCMAKE_INSTALL_PREFIX="$LIBS_DIR" \
  -DBUILD_AV_SUPPORT=OFF \
  -DENABLE_STATIC=ON \
  -DBUILD_SHARED=OFF \
  -DBOOTSTRAP_DAEMON=OFF \
  -DENABLE_TESTS=OFF \
  -DENABLE_EXAMPLES=OFF \
  -DENABLE_PROGRAMS=OFF \
  -DCMAKE_FIND_ROOT_PATH="$LIBS_DIR"

make -j$(nproc)
make install

# 4. Sync for main build script
cp "$LIBS_DIR/lib/libtoxcore.a" "$LIBS_DIR/lib/mips-linux-gnu/" 2>/dev/null

echo "================================================="
echo "[✓] MIPS Build Complete"
echo "================================================="
