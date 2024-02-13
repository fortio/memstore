package mstore

import (
	"context"
	"net"
	"sync"
	"time"

	"fortio.org/dflag"
	"fortio.org/log"
	"fortio.org/sets"
)

var (
	peers             sets.Set[string]
	cancel            context.CancelFunc
	DNSWatchSleepTime = dflag.New(15*time.Second, "Sleep time between DNS resolution")
	wg                sync.WaitGroup
)

// Resolve the service name to a list of IPs.
func checkDNS(serviceName string) {
	ips, err := net.LookupHost(serviceName)
	if err != nil {
		log.Errf("Error resolving service %q: %v", serviceName, err)
		return
	}
	newIPs := sets.FromSlice(ips)
	log.LogVf("Resolved %q to %v", serviceName, newIPs)
	// If the list changes, update the peers list
	if newIPs.Equals(peers) {
		log.LogVf("No change in peers: %v", peers)
		return
	}
	peers = newIPs.Clone()
	log.Infof("Updated peers: %v", peers)
	_ = Peers.SetV(peers)
}

func dnsWatcher(ctx context.Context, serviceName string) {
	defer wg.Done()
	for {
		checkDNS(serviceName) // first time, without waiting or cancel check
		select {
		case <-ctx.Done():
			log.Warnf("DNS Watcher for %q exiting", serviceName)
			return
		case <-time.After(DNSWatchSleepTime.Get()):
			checkDNS(serviceName)
		}
	}
}

func DNSWatcher(serviceName string) context.CancelFunc {
	ctx := context.Background()
	ctx, cancel = context.WithCancel(ctx)
	wg.Add(1)
	go dnsWatcher(ctx, serviceName)
	return cancel
}

func StartDNSWatch(serviceName string) {
	if cancel != nil {
		cancel()
	}
	cancel = DNSWatcher(serviceName)
}

func StopDNSWatch() {
	if cancel != nil {
		cancel()
		cancel = nil
		wg.Wait()
	}
}
