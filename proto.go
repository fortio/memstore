// Memstore prototype
// (c) 2023 Laurent Demailly - all rights reserved
// Apache 2.0 License

package main

import (
	"fortio.org/dflag"
	"fortio.org/log"
	"fortio.org/memstore/mstore"
	"fortio.org/scli"
)

func main() {
	dflag.Flag("peers", mstore.Peers)
	dflag.Flag("dns", mstore.DNSWatch)
	dflag.Flag("dns-interval", mstore.DNSWatchSleepTime)
	scli.ServerMain()
	log.Infof("Starting memstore prototype...")
	mstore.Start()
	scli.UntilInterrupted()
	mstore.Stop()
}
