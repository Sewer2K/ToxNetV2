#!/bin/bash
# run_with_tor.sh - Launch Toxnet C2 with Tor anonymity

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
echo "╔═══════════════════════════════════════════════════════════╗"
echo "║           Toxnet C2 - Tor Anonymous Launcher              ║"
echo "║                   Running with Tor proxy                  ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Check if running as root (required for transparent proxy)
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}[!] Please run as root (sudo) for Tor transparent proxy setup${NC}"
    echo -e "${YELLOW}[*] Alternatively, run with: torsocks ./toxnet-c2${NC}"
    exit 1
fi

# Check if Tor is installed
if ! command -v tor &> /dev/null; then
    echo -e "${RED}[!] Tor is not installed. Installing...${NC}"
    apt-get update && apt-get install -y tor
fi

# Check if toxnet-c2 binary exists
if [ ! -f "./toxnet-c2" ]; then
    echo -e "${RED}[!] toxnet-c2 binary not found. Building...${NC}"
    go build -o toxnet-c2 main.go
    if [ $? -ne 0 ]; then
        echo -e "${RED}[!] Build failed. Please fix compilation errors.${NC}"
        exit 1
    fi
fi

# Function to backup Tor config
backup_tor_config() {
    if [ -f "/etc/tor/torrc" ]; then
        cp /etc/tor/torrc /etc/tor/torrc.backup.$(date +%s)
        echo -e "${GREEN}[+] Tor config backed up${NC}"
    fi
}

# Function to restore Tor config
restore_tor_config() {
    if [ -f "/etc/tor/torrc.backup."* ]; then
        cp $(ls -t /etc/tor/torrc.backup.* | head -1) /etc/tor/torrc
        echo -e "${GREEN}[+] Tor config restored${NC}"
    fi
}

# Configure Tor for transparent proxy
configure_tor() {
    backup_tor_config
    
    echo -e "${YELLOW}[*] Configuring Tor for transparent proxy...${NC}"
    
    # Create Tor config for transparent proxy
    cat > /etc/tor/torrc << EOF
# Tor transparent proxy configuration for Toxnet C2
Log notice file /var/log/tor/notices.log
VirtualAddrNetworkIPv4 10.192.0.0/10
AutomapHostsOnResolve 1
TransPort 9040
DNSPort 5353
DNSListenAddress 127.0.0.1

# Use bridges if needed (uncomment to use)
# UseBridges 1
# Bridge obfs4 <IP>:<PORT> <FINGERPRINT>

# Exit nodes - can restrict to specific countries (optional)
# ExitNodes {US},{CA},{GB}
# StrictNodes 1
EOF

    echo -e "${GREEN}[+] Tor configuration updated${NC}"
}

# Configure iptables for transparent proxy
configure_iptables() {
    echo -e "${YELLOW}[*] Configuring iptables for Tor transparent proxy...${NC}"
    
    # Backup current iptables
    iptables-save > /tmp/iptables.backup.$(date +%s)
    
    # Flush existing rules
    iptables -t nat -F
    iptables -F
    
    # Tor Transparent Proxy Rules
    iptables -t nat -A OUTPUT -p tcp --dport 9040 -j REDIRECT --to-ports 9040 2>/dev/null || true
    iptables -t nat -A OUTPUT -p udp --dport 53 -j REDIRECT --to-ports 5353 2>/dev/null || true
    
    # Don't torify localhost, local networks
    iptables -t nat -A OUTPUT -d 127.0.0.0/8 -j RETURN
    iptables -t nat -A OUTPUT -d 10.0.0.0/8 -j RETURN
    iptables -t nat -A OUTPUT -d 172.16.0.0/12 -j RETURN
    iptables -t nat -A OUTPUT -d 192.168.0.0/16 -j RETURN
    
    # Torify everything else
    iptables -t nat -A OUTPUT -p tcp -m owner --uid-owner debian-tor -j RETURN
    iptables -t nat -A OUTPUT -p tcp -j REDIRECT --to-ports 9040
    
    echo -e "${GREEN}[+] iptables configured${NC}"
}

# Function to get Tor circuit information
get_tor_info() {
    echo -e "\n${BLUE}════════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}[+] Tor Circuit Information:${NC}"
    
    # Get Tor IP using curl through Tor
    if command -v curl &> /dev/null; then
        TOR_IP=$(curl --socks5-hostname 127.0.0.1:9050 -s https://check.torproject.org/api/ip 2>/dev/null | grep -o '"IP":"[^"]*"' | cut -d'"' -f4)
        if [ ! -z "$TOR_IP" ]; then
            echo -e "    Tor Exit IP: ${GREEN}$TOR_IP${NC}"
        fi
    fi
    
    # Show Tor circuit
    if [ -f /var/run/tor/control_cookie ]; then
        echo -e "    Tor Circuit: Active"
        echo -e "    Transparent Proxy: ${GREEN}127.0.0.1:9040${NC}"
        echo -e "    DNS Proxy: ${GREEN}127.0.0.1:5353${NC}"
    fi
}

# Start Tor service
start_tor() {
    echo -e "${YELLOW}[*] Starting Tor service...${NC}"
    
    # Kill existing Tor processes
    pkill tor 2>/dev/null || true
    sleep 2
    
    # Start Tor
    tor --runasdaemon 0 &
    TOR_PID=$!
    
    # Wait for Tor to initialize
    echo -ne "${YELLOW}[*] Waiting for Tor to bootstrap"
    for i in {1..30}; do
        if nc -z 127.0.0.1 9050 2>/dev/null; then
            echo -e " ${GREEN}✓${NC}"
            break
        fi
        echo -ne "."
        sleep 1
    done
    
    echo -e "${GREEN}[+] Tor started (PID: $TOR_PID)${NC}"
}

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}[*] Cleaning up...${NC}"
    
    # Restore iptables
    if [ -f /tmp/iptables.backup.* ]; then
        iptables-restore < $(ls -t /tmp/iptables.backup.* | head -1)
        echo -e "${GREEN}[+] iptables restored${NC}"
    fi
    
    # Restore Tor config
    restore_tor_config
    
    # Kill Tor
    pkill tor 2>/dev/null || true
    
    # Kill toxnet-c2 if running
    pkill -f toxnet-c2 2>/dev/null || true
    
    echo -e "${GREEN}[+] Cleanup complete${NC}"
    exit 0
}

# Trap cleanup on exit
trap cleanup INT TERM EXIT

# Main execution
main() {
    echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
    
    # Check for existing toxnet-c2 processes
    if pgrep -f toxnet-c2 > /dev/null; then
        echo -e "${RED}[!] toxnet-c2 is already running. Killing existing instance...${NC}"
        pkill -f toxnet-c2
        sleep 2
    fi
    
    # Configure and start Tor
    configure_tor
    configure_iptables
    start_tor
    
    # Get Tor info
    sleep 3
    get_tor_info
    
    # Launch Toxnet C2
    echo -e "\n${BLUE}════════════════════════════════════════════════════════════${NC}"
    echo -e "${GREEN}[+] Launching Toxnet C2 through Tor...${NC}"
    echo -e "${YELLOW}[!] All C2 traffic will be routed through the Tor network${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}\n"
    
    # Set environment variables for Tor proxy
    export http_proxy="socks5://127.0.0.1:9050"
    export https_proxy="socks5://127.0.0.1:9050"
    export all_proxy="socks5://127.0.0.1:9050"
    
    # Run toxnet-c2 with Tor
    ./toxnet-c2
    
    # Keep running until interrupted
    wait
}

# Run main function
main