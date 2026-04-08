#!/bin/bash
cd "$(dirname "$0")"

echo "========== BOT ANALYSIS REPORT =========="
echo ""

for bot in bot_*; do
    if [ ! -f "$bot" ]; then continue; fi
    
    echo "========================================="
    echo "📁 FILE: $bot"
    echo "========================================="
    
    # File type
    echo "📋 FILE TYPE:"
    file "$bot" | sed 's/^/   /'
    echo ""
    
    # ELF Header
    echo "🔧 ELF HEADER:"
    readelf -h "$bot" 2>/dev/null | grep -E "(Class|Machine|OS/ABI|Type|Entry point)" | sed 's/^/   /'
    echo ""
    
    # Size
    echo "📦 SIZE: $(du -h "$bot" | cut -f1)"
    echo ""
    
    # Static or dynamic
    if readelf -d "$bot" 2>/dev/null | grep -q "NEEDED"; then
        echo "⚠️  WARNING: Dynamically linked (needs libraries)"
        readelf -d "$bot" 2>/dev/null | grep "NEEDED" | sed 's/^/   /'
    else
        echo "✅ Static: Statically linked (no dependencies)"
    fi
    echo ""
    
    # Stripped status
    if readelf -S "$bot" 2>/dev/null | grep -q "symtab"; then
        echo "⚠️  NOT STRIPPED (has symbol table)"
    else
        echo "✅ Stripped: No symbols (harder to analyze)"
    fi
    echo ""
done

echo "========== END OF REPORT =========="
