// Memstore prototype
// (c) 2023 Laurent Demailly - all rights reserved
// Apache 2.0 License

package main

import (
	"os"

	"fortio.org/dflag"
	"fortio.org/log"
	"fortio.org/memstore/mstore"
	"fortio.org/scli"
)

func main() {
	dflag.Flag("peers", mstore.Peers)
	dflag.Flag("dns", mstore.DNSWatch)
	dflag.Flag("dns-interval", mstore.DNSWatchSleepTime)
	dflag.FlagBool("statefulset", mstore.StatefulSet)
	scli.ServerMain()
	if mstore.StatefulSet.Get() && mstore.DNSWatch.Get() == "" {
		log.Fatalf("StatefulSet mode needs -dns to be set")
	}
	if mstore.DNSWatch.Get() != "" && mstore.Peers.Get().Len() > 0 {
		log.Fatalf("Can only have either -peers or -dns set, not both")
	}
	myName, found := os.LookupEnv("NAME")
	if mstore.StatefulSet.Get() && !found {
		log.Fatalf("No NAME env var found for statefulset mode (to this pod's name)")
	}
	mstore.Start(myName)
	scli.UntilInterrupted()
	mstore.Stop()
}
