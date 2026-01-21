package util_test

import (
	"encoding/binary"
	"net"
	"testing"

	"github.com/free-ran-ue/frush/util"
)

var testBuildIcmpEchoPacketCases = []struct {
	name  string
	srcIp string
	dstIp string
}{
	{
		name:  "test_build_icmp_echo_packet",
		srcIp: "10.60.0.1",
		dstIp: "1.1.1.1",
	},
}

func TestBuildIcmpEchoPacket(t *testing.T) {
	for _, testCase := range testBuildIcmpEchoPacketCases {
		t.Run(testCase.name, func(t *testing.T) {
			pkt, err := util.BuildIcmpEchoPacket(testCase.srcIp, testCase.dstIp)
			if err != nil {
				t.Fatalf("BuildICMPEchoPacket failed: %v", err)
			}

			// Total length check
			if len(pkt) != 28 {
				t.Fatalf("unexpected packet length: got %d, want 28", len(pkt))
			}

			// ---------- IPv4 header ----------
			ip := pkt[:20]

			// Version & IHL
			if ip[0] != 0x45 {
				t.Fatalf("invalid IP version/IHL: 0x%x", ip[0])
			}

			// Protocol = ICMP (1)
			if ip[9] != 1 {
				t.Fatalf("unexpected IP protocol: got %d, want 1 (ICMP)", ip[9])
			}

			// Source IP
			gotSrc := net.IP(ip[12:16]).String()
			if gotSrc != testCase.srcIp {
				t.Fatalf("unexpected src IP: got %s, want %s", gotSrc, testCase.srcIp)
			}

			// Destination IP
			gotDst := net.IP(ip[16:20]).String()
			if gotDst != testCase.dstIp {
				t.Fatalf("unexpected dst IP: got %s, want %s", gotDst, testCase.dstIp)
			}

			// IP checksum validation
			if util.CheckSum(ip) != 0 {
				t.Fatalf("invalid IP checksum")
			}

			// ---------- ICMP ----------
			icmp := pkt[20:]

			// Type = Echo Request (8)
			if icmp[0] != 8 {
				t.Fatalf("unexpected ICMP type: got %d, want 8", icmp[0])
			}

			// Code = 0
			if icmp[1] != 0 {
				t.Fatalf("unexpected ICMP code: got %d, want 0", icmp[1])
			}

			// Identifier
			id := binary.BigEndian.Uint16(icmp[4:])
			if id != 0x1234 {
				t.Fatalf("unexpected ICMP identifier: got 0x%x", id)
			}

			// Checksum validation
			if util.CheckSum(icmp) != 0 {
				t.Fatalf("invalid ICMP checksum")
			}
		})
	}
}

var testIsIcmpEchoReplyCases = []struct {
	name  string
	srcIp string
	dstIp string
}{
	{
		name:  "test_is_icmp_echo_reply",
		srcIp: "10.60.0.1",
		dstIp: "1.1.1.1",
	},
}

func TestIsIcmpEchoReply(t *testing.T) {
	buildEchoReply := func(srcIP, dstIP string) []byte {
		// ICMP
		icmp := make([]byte, 8)
		icmp[0] = 0 // Echo Reply
		binary.BigEndian.PutUint16(icmp[4:], 0x1234)
		binary.BigEndian.PutUint16(icmp[6:], 1)
		binary.BigEndian.PutUint16(icmp[2:], util.CheckSum(icmp))

		// IPv4
		ip := make([]byte, 20)
		ip[0] = 0x45
		ip[8] = 64
		ip[9] = 1
		copy(ip[12:], net.ParseIP(srcIP).To4())
		copy(ip[16:], net.ParseIP(dstIP).To4())
		binary.BigEndian.PutUint16(ip[2:], uint16(28))
		binary.BigEndian.PutUint16(ip[10:], util.CheckSum(ip))

		return append(ip, icmp...)
	}

	for _, testCase := range testIsIcmpEchoReplyCases {
		t.Run(testCase.name, func(t *testing.T) {
			reply := buildEchoReply(testCase.dstIp, testCase.srcIp)
			ok, err := util.IsIcmpEchoReply(reply, testCase.srcIp, testCase.dstIp)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatalf("expected valid echo reply")
			}
		})
	}
}
