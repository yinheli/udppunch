package netx

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	headerLen = 8
)

type UDPConn struct {
	ipConn   *net.IPConn
	srcPort  uint16
	destPort uint16
}

func Dial(ip net.IP, srcPort uint16, destPort uint16) (*UDPConn, error) {
	ipConn, err := net.DialIP("ip4:udp", nil, &net.IPAddr{
		IP: ip,
	})
	if err != nil {
		return nil, err
	}
	return &UDPConn{
		ipConn:   ipConn,
		srcPort:  srcPort,
		destPort: destPort,
	}, nil
}

func (c *UDPConn) Write(b []byte) (int, error) {
	if len(b) > 0xffff-headerLen {
		return 0, fmt.Errorf("datagram payload too large (max 0xffff - 8)")
	}
	enc := c.header(uint16(len(b)))
	enc = append(enc, b...)
	return c.ipConn.Write(enc)
}

func (c *UDPConn) Close() error {
	return c.ipConn.Close()
}

// header returns the encoded UDP datagram header. This header has format:
//  0      7 8     15 16    23 24    31
// +--------+--------+--------+--------+
// |     Source      |   Destination   |
// |      Port       |      Port       |
// +--------+--------+--------+--------+
// |                 |                 |
// |     Length      |    Checksum     |
// +--------+--------+--------+--------+
// Note includes a checksum of 0 to indicate no checksum.
func (c *UDPConn) header(payloadLen uint16) []byte {
	b := make([]byte, headerLen)
	binary.BigEndian.PutUint16(b[0:], c.srcPort)
	binary.BigEndian.PutUint16(b[2:], c.destPort)
	binary.BigEndian.PutUint16(b[4:], headerLen+payloadLen)
	// No checksum.
	binary.BigEndian.PutUint16(b[6:], 0)
	return b
}
