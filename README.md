# udppunch

udp punch for wireguard

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
