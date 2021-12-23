package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	lru "github.com/hashicorp/golang-lru"
	"github.com/yinheli/udppunch"
)

var (
	l       = log.New(os.Stdout, "", log.LstdFlags)
	port    = flag.Int("port", 19993, "udp punch port")
	version = flag.Bool("version", false, "show version")
)

func main() {
	if flag.Parse(); !flag.Parsed() {
		flag.Usage()
		os.Exit(1)
	}

	if *version {
		fmt.Println(udppunch.Version)
		os.Exit(0)
	}

	peers, _ := lru.New(1024)

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		l.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)

	if err != nil {
		l.Fatal(err)
	}

	for {
		buf := make([]byte, 1024*8)
		n, raddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		if n < 1 {
			continue
		}

		// l.Printf("\nfrom:%v\n%s", raddr, hex.Dump(buf[:n]))

		switch buf[0] {
		case udppunch.HandshakeType:
			var key udppunch.Key
			copy(key[:], buf[1:])
			peers.Add(key, udppunch.NewPeerFromAddr(key, raddr))
		case udppunch.ResolveType:
			data := make([]byte, 0, (n-1)/32*38)
			for i := 1; i < n; i += 32 {
				var key udppunch.Key
				copy(key[:], buf[i:i+32])
				if v, ok := peers.Get(key); ok {
					peer := v.(udppunch.Peer)
					data = append(data, peer[:]...)
				}
			}
			conn.WriteToUDP(data, raddr)
		}
	}
}
