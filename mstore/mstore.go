package mstore

import (
	"fortio.org/dflag"
	"fortio.org/log"
	"fortio.org/sets"
)

var Peers = dflag.New(
	sets.Set[string]{},
	"Peers to connect to (comma separated set)",
).WithNotifier(peerChange)

func connect(p string) {
	log.Infof("Connecting to peer     : %q", p)
}

func disconnect(p string) {
	log.Infof("Disconnecting from peer: %q", p)
}

func peerChange(oldValue, newValue sets.Set[string]) {
	log.Infof("Peer set changed from %v to %v", oldValue, newValue)
	sets.RemoveCommon(oldValue, newValue)
	for _, p := range sets.Sort(newValue) {
		connect(p)
	}
	for _, p := range sets.Sort(oldValue) {
		disconnect(p)
	}
}

func Start() {
	log.Infof("memstore Start()")
	// peerChange does get call for even initial flag value
}
