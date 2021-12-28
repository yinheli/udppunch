# udppunch

udp punch for wireguard, inspired by [natpunch-go](https://github.com/malcolmseyd/natpunch-go)

## usage

server side

```bash
./punch-server-linux-amd64 -port 19993
```

client side

> make sure wireguard is up

```bash
./dist/punch-client-linux-amd64 -server xxxx:19993 -iface wg0
```

## resource

- [natpunch-go](https://github.com/malcolmseyd/natpunch-go) (because of [#7](https://github.com/malcolmseyd/natpunch-go/issues/7) not support macOS, so I build this)
- [wireguard-vanity-address](https://github.com/yinheli/wireguard-vanity-address) generate keypairs with a given prefix
- [UDP hole punching](https://en.wikipedia.org/wiki/UDP_hole_punching)
