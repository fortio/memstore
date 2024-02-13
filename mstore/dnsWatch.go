package mstore

import (
	"context"
	"net"
	"time"

	"fortio.org/dflag"
	"fortio.org/log"
	"fortio.org/sets"
)

var (
	peers             sets.Set[string]
	cancel            context.CancelFunc
	DNSWatchSleepTime = dflag.New(15*time.Second, "Sleep time between DNS resolution")
)

func dnsWatcher(ctx context.Context, serviceName string) {
	for {
		select {
		case <-ctx.Done():
			log.Warnf("DNS Watcher for %q exiting", serviceName)
			return
		default:
			// Resolve the service name to a list of IPs
			ips, err := net.LookupHost(serviceName)
			if err != nil {
				log.Errf("Error resolving service %q: %v", serviceName, err)
				time.Sleep(DNSWatchSleepTime.Get()) // Sleep for a minute before resolving again
				continue
			}
			newIPs := sets.FromSlice(ips)
			log.LogVf("Resolved %q to %v", serviceName, newIPs)
			// If the list changes, update the peers list
			if !newIPs.Equals(peers) {
				peers = newIPs.Clone()
				log.Infof("Updated peers: %v", peers)
				_ = Peers.SetV(peers)
			} else {
				log.LogVf("No change in peers: %v", peers)
			}
			time.Sleep(DNSWatchSleepTime.Get()) // Sleep for a minute before resolving again
		}
	}
}

func DNSWatcher(serviceName string) context.CancelFunc {
	ctx := context.Background()
	ctx, cancel = context.WithCancel(ctx)
	go dnsWatcher(ctx, serviceName)
	return cancel
}

func StartDNSWatch(serviceName string) {
	if cancel != nil {
		cancel()
	}
	cancel = DNSWatcher(serviceName)
}
