// Memstore prototype
// (c) 2023 Laurent Demailly - all rights reserved
// Apache 2.0 License

package main

import (
	"fortio.org/dflag"
	"fortio.org/memstore/mstore"
	"fortio.org/scli"
)

func main() {
	dflag.Flag("peers", mstore.Peers)
	scli.ServerMain()
	mstore.Start()
	scli.UntilInterrupted()
}
