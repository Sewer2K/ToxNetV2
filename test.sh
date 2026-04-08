#!/bin/bash
# save as mass_test.sh

declare -A BOTS
BOTS["bot_i386"]="qemu-i386-static"
BOTS["bot_armv7"]="qemu-arm-static"
BOTS["bot_armv5"]="qemu-arm-static"
BOTS["bot_arm64"]="qemu-aarch64-static"
BOTS["bot_mips"]="qemu-mips-static"
BOTS["bot_mipsel"]="qemu-mipsel-static"
BOTS["bot_mips64el"]="qemu-mips64el-static"
BOTS["bot_ppc64el"]="qemu-ppc64le-static"
BOTS["bot_s390x"]="qemu-s390x-static"

echo "[+] Launching all bots in the background..."
for bot in "${!BOTS[@]}"; do
    if [ -f "bots/$bot" ]; then
        ${BOTS[$bot]} bots/$bot &
        echo "[>] Started $bot"
    else
        echo "[!] $bot missing, skipping."
    fi
done

# Wait for them to daemonize
echo "[+] Waiting 3 seconds for bots to stabilize..."
sleep 3

# Verify running processes
echo -e "\n[+] VERIFYING PROCESS LIST:"
ps aux | grep qemu | grep -v grep

# Verify network sockets (Tox uses UDP)
echo -e "\n[+] VERIFYING ACTIVE UDP SOCKETS:"
ss -unap | grep qemu