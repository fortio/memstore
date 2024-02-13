package mstore

import (
	"fortio.org/dflag"
	"fortio.org/log"
	"fortio.org/sets"
)

var (
	// Either direct peer lists (yet dynamic).
	Peers = dflag.New(
		sets.Set[string]{},
		"Peers to connect to (comma separated set)",
	).WithNotifier(peerChange)
	// or DNS based discovery/watch.
	DNSWatch = dflag.New("", "DNS service name to watch for peers").WithNotifier(dnsChange)
)

func connect(p string) {
	log.Infof("Connecting to peer     : %q", p)
}

func disconnect(p string) {
	log.Infof("Disconnecting from peer: %q", p)
}

func peerChange(oldValue, newValue sets.Set[string]) {
	log.Infof("Peer set changed from %v to %v", oldValue, newValue)
	// Make copy of newValue so we don't mutate the flag's value
	newValue = newValue.Clone()
	sets.RemoveCommon(oldValue, newValue)
	for _, p := range sets.Sort(newValue) {
		connect(p)
	}
	for _, p := range sets.Sort(oldValue) {
		disconnect(p)
	}
}

func dnsChange(oldValue, newValue string) {
	log.Infof("DNSWatch changed from %q to %q", oldValue, newValue)
	StartDNSWatch(newValue)
}

func Start() {
	log.Infof("memstore Start()")
	// peerChange does get call for even initial flag value
}
