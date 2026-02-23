package discovery

import (
	"context"
	"strings"

	"github.com/grandcat/zeroconf"
)

type DiscoveredDevice struct {
	Name string
	IP   string
	Port int
}

// Scan searches for Google Home/Chromecast devices on the local network.
// It returns a channel that emits devices as they are discovered.
// The channel is closed when the context expires and all entries are drained.
func Scan(ctx context.Context) (<-chan DiscoveredDevice, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	results := make(chan DiscoveredDevice)

	go func() {
		defer close(results)
		for entry := range entries {
			// Skip entries without IPv4 addresses
			if len(entry.AddrIPv4) == 0 {
				continue
			}

			// Try to get friendly name from TXT records
			// Google Cast devices advertise "fn=Friendly Name" in TXT
			name := extractFriendlyName(entry.Text)
			if name == "" {
				// Fallback to instance name
				name = entry.Instance
				if name == "" {
					name = entry.ServiceInstanceName()
				}
				name = strings.TrimSuffix(name, "._googlecast._tcp.local.")
			}

			results <- DiscoveredDevice{
				Name: name,
				IP:   entry.AddrIPv4[0].String(),
				Port: entry.Port,
			}
		}
	}()

	// Browse for Google Cast devices
	err = resolver.Browse(ctx, "_googlecast._tcp", "local.", entries)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// extractFriendlyName looks for "fn=..." in TXT records
func extractFriendlyName(txt []string) string {
	for _, record := range txt {
		if strings.HasPrefix(record, "fn=") {
			return strings.TrimPrefix(record, "fn=")
		}
	}
	return ""
}
