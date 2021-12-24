package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/yinheli/udppunch"
	"github.com/yinheli/udppunch/client/netx"
	"github.com/yinheli/udppunch/client/wg"
)

type Peer [32]byte

var (
	l          = log.New(os.Stdout, "", log.LstdFlags)
	iface      = flag.String("iface", "wg0", "wireguard interface")
	server     = flag.String("server", "", "udp punch server")
	continuous = flag.Bool("continuous", false, "continuously resolve peers")
	version    = flag.Bool("version", false, "show version")
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

	if *server == "" {
		l.Fatal("server is empty")
	}

	if *iface == "" {
		l.Fatal("iface is empty")
	}

	raddr, err := net.ResolveUDPAddr("udp", *server)

	if err != nil {
		l.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		l.Fatal(err)
	}

	// handshake
	go handshake(*raddr)

	// wait for handshake
	time.Sleep(time.Second * 2)

	// resolve
	for {
		clients, err := wg.GetEndpoints(*iface)
		if err != nil {
			l.Print("get endpoints error:", err)
			time.Sleep(time.Second * 10)
			continue
		}

		data := make([]byte, 0, 1+len(clients)*32)
		data = append(data, udppunch.ResolveType)

		for client := range clients {
			data = append(data, client[:]...)
		}
		conn.Write(data)

		buf := make([]byte, 4096)
		n, err := conn.Read(buf)
		if err != nil {
			l.Fatal(err)
		}

		if n < 38 {
			time.Sleep(time.Second * 5)
			continue
		}

		peers := udppunch.ParsePeers(buf[:n])
		for _, peer := range peers {
			key, addr := peer.Parse()
			if clients[key] == addr {
				continue
			}
			l.Print("reslove ", key, " ", addr)
			err = wg.SetPeerEndpoint(*iface, key, addr)
			if err != nil {
				l.Printf("set peer (%v) endpoint error: %v", key, err)
			}
		}

		if err == nil && !*continuous {
			break
		} else {
			// sleep for a while then continue resolve
			time.Sleep(time.Second * 5)
		}
	}
}

func handshake(raddr net.UDPAddr) {
	defer func() {
		if x := recover(); x != nil {
			l.Print("handshake error:", x)
			time.Sleep(time.Second * 10)
			handshake(raddr)
		}
	}()

	for {
		port, err := wg.GetIfaceListenPort(*iface)
		if err != nil {
			l.Print("get interface listen-port:", err)
			time.Sleep(time.Second * 10)
			continue
		}

		pubKey, err := wg.GetIfacePubKey(*iface)
		if err != nil {
			l.Fatal("get interface public key:", err)
			time.Sleep(time.Second * 10)
			continue
		}

		doHandshake(raddr.IP, port, uint16(raddr.Port), pubKey)

		time.Sleep(time.Second * 25)
	}
}

func doHandshake(ip net.IP, srcPort uint16, dstPort uint16, pubKey udppunch.Key) {
	conn, err := netx.Dial(ip, srcPort, dstPort)
	if err != nil {
		l.Print("handshake dial error:", err)
		time.Sleep(time.Second * 10)
		return
	}

	defer conn.Close()

	data := make([]byte, 0, 32+1)
	data = append(data, udppunch.HandshakeType)
	data = append(data, pubKey[:]...)

	conn.Write(data)
}
