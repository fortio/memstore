# memstore
Distributed HA in memory store for Golang

## Config Input
(all [dflags](https://github.com/fortio/dflag#fortio-dynamic-flags) so can be changed without restart)

- List of DNS names, IPs (in Kubernetes you'd pass just a headless service name)
- Refresh frequency for DNS to IP

### Prototype
```
go run . -config-port 7999
```

Then go change the `peers` on https://localhost:7999 to see:
```
20:26:10 I mstore.go:17> Peer set changed from  to a,b,c,z
20:26:10 I mstore.go:23> Connecting to added Peer        : "a"
20:26:10 I mstore.go:23> Connecting to added Peer        : "b"
20:26:10 I mstore.go:23> Connecting to added Peer        : "c"
20:26:10 I mstore.go:23> Connecting to added Peer        : "z"
```
and
```
20:26:31 I mstore.go:17> Peer set changed from a,b,c,z to d,a,b,z
20:26:31 I mstore.go:20> Disconnecting from removed peer : "c"
20:26:31 I mstore.go:23> Connecting to added Peer        : "d"
```

or similar
```
make test
```

## Communication

Should we
- use some broadcasting/bus
- ring
- tcp or http or grpc

Let's use a fully mesh broadcast using point2point h2.

## Embedded or separate

Why not both

## Protocol

- Zookeeper
- Raft
- Something wrong but simpler (*)

## Persistence

- Dump to disk (Persistent Volume in k8s) periodically

## CircularBuffer

Both pub/sub thread safe blocking version and pure FIFO queue with set capacity versions:

See [cb/cb.go](cb/cb.go)
