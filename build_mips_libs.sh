#!/bin/bash
# save as build_mips_libs.sh

BASE_DIR="/home/vboxuser/Desktop/toxnet"
LIBS_DIR="$BASE_DIR/cross_libs/mips/usr"
# Ensure this matches the folder name YOU created
SOURCE_DIR="$BASE_DIR/c-toxcore-0.2.22" 

mkdir -p "$LIBS_DIR/lib/mips-linux-gnu"
cd "$BASE_DIR" || exit 1

# 1. Build Libsodium (if missing)
if [ ! -f "$LIBS_DIR/lib/mips-linux-gnu/libsodium.a" ]; then
    echo "[+] Building libsodium for MIPS..."
    cd "$BASE_DIR/libsodium-stable" || exit 1
    ./configure --host=mips-linux-gnu \
                --prefix="$LIBS_DIR" \
                --libdir="$LIBS_DIR/lib/mips-linux-gnu" \
                --enable-static --disable-shared
    make -j$(nproc) && make install
    cd "$BASE_DIR"
fi

# 2. Build Toxcore
echo "[+] Building toxcore for MIPS..."
rm -rf toxcore_build_mips && mkdir -p toxcore_build_mips && cd toxcore_build_mips

# Crucial: Tell PkgConfig to look at our MIPS folder, not Ubuntu's system folders
export PKG_CONFIG_LIBDIR="$LIBS_DIR/lib/mips-linux-gnu/pkgconfig:$LIBS_DIR/lib/pkgconfig"
export PKG_CONFIG_SYSROOT_DIR="$BASE_DIR/cross_libs/mips"

cmake "$SOURCE_DIR" \
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
  -DMINCORE=OFF

# Use 'make' instead of 'make toxcore_static' to ensure all dependencies are hit
make -j$(nproc)
make install

# 3. Final placement for your main build script
cp "$LIBS_DIR/lib/libtoxcore.a" "$LIBS_DIR/lib/mips-linux-gnu/" 2>/dev/null

echo "================================================="
echo "[✓] MIPS Build Attempt Finished"
echo "================================================="
