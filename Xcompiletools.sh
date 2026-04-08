#!/bin/bash
# save as setup_cross_compile.sh

echo "[+] Installing cross-compilers..."
sudo apt update
sudo apt install -y gcc-arm-linux-gnueabi g++-arm-linux-gnueabi \
                    gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf \
                    gcc-aarch64-linux-gnu g++-aarch64-linux-gnu \
                    gcc-mips-linux-gnu g++-mips-linux-gnu \
                    gcc-mipsel-linux-gnu g++-mipsel-linux-gnu \
                    gcc-powerpc-linux-gnu g++-powerpc-linux-gnu \
                    mingw-w64

echo "[+] Installing static libraries..."
sudo apt install -y libsodium-dev:arm64 libsodium-dev:armhf libsodium-dev:armel