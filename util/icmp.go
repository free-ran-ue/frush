package util

import (
	"encoding/binary"
	"fmt"
	"net"
)

func CheckSum(data []byte) uint16 {
	var sum uint32
	for i := 0; i+1 < len(data); i += 2 {
		sum += uint32(binary.BigEndian.Uint16(data[i:]))
	}
	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for (sum >> 16) > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

func BuildIcmpEchoPacket(srcIp, dstIp string) ([]byte, error) {
	src := net.ParseIP(srcIp).To4()
	dst := net.ParseIP(dstIp).To4()
	if src == nil || dst == nil {
		return nil, fmt.Errorf("invalid IP address")
	}

	// ---------- ICMP ----------
	icmp := make([]byte, 8)
	icmp[0] = 8 // Type: Echo Request
	icmp[1] = 0 // Code
	// checksum remain 0
	binary.BigEndian.PutUint16(icmp[4:], 0x1234) // Identifier
	binary.BigEndian.PutUint16(icmp[6:], 1)      // Sequence

	icmpCsum := CheckSum(icmp)
	binary.BigEndian.PutUint16(icmp[2:], icmpCsum)

	// ---------- IPv4 ----------
	ip := make([]byte, 20)
	ip[0] = 0x45 // Version=4, IHL=5
	ip[1] = 0x00 // TOS
	binary.BigEndian.PutUint16(ip[2:], uint16(20+len(icmp)))
	binary.BigEndian.PutUint16(ip[4:], 0x0000) // Identification
	binary.BigEndian.PutUint16(ip[6:], 0x0000) // Flags + Fragment offset
	ip[8] = 64                                 // TTL
	ip[9] = 1                                  // Protocol = ICMP
	// checksum remain 0
	copy(ip[12:], src)
	copy(ip[16:], dst)

	ipCsum := CheckSum(ip)
	binary.BigEndian.PutUint16(ip[10:], ipCsum)

	// ---------- merge ----------
	packet := append(ip, icmp...)
	return packet, nil
}

func IsIcmpEchoReply(pkt []byte, reqSrcIP, reqDstIP string) (bool, error) {
	if len(pkt) < 28 {
		return false, fmt.Errorf("packet too short")
	}

	reqSrc := net.ParseIP(reqSrcIP).To4()
	reqDst := net.ParseIP(reqDstIP).To4()
	if reqSrc == nil || reqDst == nil {
		return false, fmt.Errorf("invalid request IP")
	}

	// ---------- IPv4 ----------
	ip := pkt[:20]

	if ip[0]>>4 != 4 {
		return false, fmt.Errorf("invalid ip version")
	}

	if ip[9] != 1 { // ICMP
		return false, fmt.Errorf("invalid ip protocol")
	}

	if CheckSum(ip) != 0 {
		return false, fmt.Errorf("invalid ip checksum")
	}

	replySrc := net.IP(ip[12:16])
	replyDst := net.IP(ip[16:20])

	// IP must be reversed
	if !replySrc.Equal(reqDst) || !replyDst.Equal(reqSrc) {
		return false, fmt.Errorf("ip must be reversed")
	}

	// ---------- ICMP ----------
	icmp := pkt[20:]

	if icmp[0] != 0 || icmp[1] != 0 {
		return false, fmt.Errorf("icmp must be 0")
	}

	if CheckSum(icmp) != 0 {
		return false, fmt.Errorf("invalid icmp checksum")
	}

	return true, nil
}
