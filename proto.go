// Memstore prototype
// (c) 2023 Laurent Demailly - all rights reserved
// Apache 2.0 License

package main

import (
	"os"
	"os/signal"

	"fortio.org/dflag"
	"fortio.org/log"
	"fortio.org/memstore/mstore"
	"fortio.org/scli"
)

func main() {
	dflag.Flag("peers", mstore.Peers)
	scli.ServerMain()
	mstore.Start()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	// Block until an INT signal is received
	<-sigChan
	log.Warnf("\nReceived INT signal, shutting down...")
}
