package mstore

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"fortio.org/dflag"
	"fortio.org/log"
	"fortio.org/sets"
)

var (
	DNSWatchSleepTime = dflag.New(15*time.Second, "Sleep time between DNS resolution")
	// StatefulSet special handling (reverse DNS to name-0, name-1, etc.).
	StatefulSet = dflag.NewBool(false,
		"Deployment is a stateful set, we will reverse the dns to get the name-0, name-1... name-n peers")
	peerIPs   sets.Set[string]
	peerNames sets.Set[string]
	// In stateful set mode we do extra reverse lookup to get names and last index.
	numPeers int
	myName   string
	cancel   context.CancelFunc
	wg       sync.WaitGroup
)

func reverseDNS(ips sets.Set[string]) (sets.Set[string], bool) {
	allNames := sets.Set[string]{}
	hasError := false
	for ip := range ips {
		names, err := net.LookupAddr(ip)
		if err != nil {
			log.Errf("Error resolving IP %q: %v", ip, err)
			hasError = true
			continue
		}
		if len(names) == 0 {
			log.Errf("No names found for IP %q", ip)
			continue
		}
		log.LogVf("Names found for IP %q: %v", ip, names)
		// Pick first one and remove anything after first dot (e.g memstore-0.memstore.memstore.svc.cluster.local -> memstore-0)
		shortName := names[0]
		dotIndex := strings.Index(shortName, ".")
		if dotIndex != -1 {
			shortName = shortName[:dotIndex]
		}
		allNames.Add(shortName)
	}
	if !allNames.Has(myName) {
		log.Errf("My name %q not found in the reverse DNS list: %v", myName, allNames)
	}
	return allNames, hasError
}

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
	if newIPs.Equals(peerIPs) {
		log.LogVf("No change in peers: %v", peerIPs)
		return
	}
	if StatefulSet.Get() {
		var hasErrors bool
		peerNames, hasErrors = reverseDNS(newIPs)
		numPeers = len(peerNames)
		log.Infof("StatefulSet mode, errs %v found %d peers (including ourselves %q) for %q: %v",
			hasErrors, numPeers, myName, serviceName, peerNames)
		if hasErrors {
			return
		}
	}
	peerIPs = newIPs.Clone()
	log.Infof("Updated peers: %v", peerIPs)
	_ = Peers.SetV(peerIPs)
}

func dnsWatcher(ctx context.Context, serviceName string) {
	defer wg.Done()
	checkDNS(serviceName) // first time, without waiting or cancel check
	for {
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
		peerIPs.Clear()
		peerNames.Clear()
		numPeers = 0
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
