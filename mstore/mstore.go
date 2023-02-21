package mstore

import (
	"fortio.org/dflag"
	"fortio.org/log"
	"fortio.org/sets"
)

var (
	Peers = dflag.New(
		sets.Set[string]{},
		"Peers to connect to (comma separated set)",
	).WithNotifier(peerChange)
)

func peerChange(oldValue, newValue sets.Set[string]) {
	log.Infof("Peer set changed from %v to %v", oldValue, newValue)
	sets.RemoveCommon(oldValue, newValue)
	for _, p := range sets.Sort(oldValue) {
		log.Infof("Disconnecting from removed peer : %q", p)
	}
	for _, p := range sets.Sort(newValue) {
		log.Infof("Connecting to added peer        : %q", p)
	}
}

func Start() {
	log.Infof("memstore Start()")
	for p := range Peers.Get() {
		log.Infof("Connecting to Peer %q", p)
	}
}
