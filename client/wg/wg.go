package wg

import (
	"encoding/base64"
	"regexp"
	"strconv"
	"strings"

	"github.com/yinheli/udppunch"
)

var (
	reSpace = regexp.MustCompile(`\s+`)
)

func GetIfacePubKey(iface string) (udppunch.Key, error) {
	r, err := run("wg", "show", iface, "public-key")
	if err != nil {
		return udppunch.DefaultKey, err
	}
	return udppunch.NewKeyFromStr(r), nil
}

func GetIfaceListenPort(iface string) (uint16, error) {
	r, err := run("wg", "show", iface, "listen-port")
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(r, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}

func GetEndpoints(iface string) (map[udppunch.Key]string, error) {
	r, err := run("wg", "show", iface, "endpoints")
	if err != nil {
		return nil, err
	}
	peers := make(map[udppunch.Key]string, 128)
	for _, it := range strings.Split(r, "\n") {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		arr := reSpace.Split(it, -1)
		if len(arr) < 2 {
			continue
		}
		peer := udppunch.NewKeyFromStr(arr[0])
		endpoint := arr[1]
		if endpoint == "(none)" {
			endpoint = ""
		}
		peers[peer] = endpoint
	}
	return peers, nil
}

func SetPeerEndpoint(iface string, peer udppunch.Key, endpoint string) error {
	_, err := run(
		"wg", "set", iface,
		"peer", base64.StdEncoding.EncodeToString(peer[:]),
		"persistent-keepalive", "10",
		"endpoint", endpoint,
	)
	if err != nil {
		return err
	}
	return nil
}
