package discovery

import (
	"context"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
)

type DiscoveredDevice struct {
	Name string
	IP   string
	Port int
}

// Scan searches for Google Home/Chromecast devices on the local network
func Scan() ([]DiscoveredDevice, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	var devices []DiscoveredDevice

	// Scan for 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
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

			devices = append(devices, DiscoveredDevice{
				Name: name,
				IP:   entry.AddrIPv4[0].String(),
				Port: entry.Port,
			})
		}
	}()

	// Browse for Google Cast devices
	err = resolver.Browse(ctx, "_googlecast._tcp", "local.", entries)
	if err != nil {
		return nil, err
	}

	// Wait for timeout
	<-ctx.Done()

	return devices, nil
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
