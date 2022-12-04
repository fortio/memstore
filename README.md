# memstore
Distributed HA in memory store for Golang

# Config Input
(all [dflags](https://github.com/fortio/fortio/tree/master/dflag#fortio-dynamic-flags-was-go-flagz) so can be changed without restart)

- List of DNS names, IPs (in Kubernetes you'd pass just a headless service name)
- Refresh frequency for DNS to IP

# Communication

Should we
- use some broadcasting/bus
- ring
- tcp or http or grpc
