// Memstore prototype
// (c) 2023 Laurent Demailly - all rights reserved
// Apache 2.0 License

package main

import (
	"flag"
	"os"

	"fortio.org/dflag"
	"fortio.org/fortio/fhttp"
	"fortio.org/log"
	"fortio.org/memstore/mstore"
	"fortio.org/memstore/probes"
	"fortio.org/scli"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	dflag.Flag("peers", mstore.Peers)
	dflag.Flag("dns", mstore.DNSWatch)
	dflag.Flag("dns-interval", mstore.DNSWatchSleepTime)
	dflag.FlagBool("statefulset", mstore.StatefulSet)
	dflag.FlagBool("ready", probes.ReadyFlag)
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
	epoch := os.Getenv("EPOCH")
	log.Infof("Starting memstore with name %q and epoch %s", myName, epoch)
	mstore.Start(myName)
	mux, addr := fhttp.HTTPServer("memstore", *port)
	if addr == nil {
		log.Fatalf("Failed to start http server")
	}
	probes.Setup(mux)
	probes.State.SetLive(true)
	probes.State.SetStarted(true)
	// For testing/changing we can use curl to set flags podip:7999/set?name=ready&value=true
	/*
		time.Sleep(50 * time.Second) // give time for the probes to be ready
		log.Warnf("Switching back to not ready")
		probes.State.SetReady(false)
	*/
	scli.UntilInterrupted()
	mstore.Stop()
}
