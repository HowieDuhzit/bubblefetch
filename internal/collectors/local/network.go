package local

import (
	"net"

	"github.com/yourusername/bubblefetch/internal/collectors"
)

// detectNetwork gathers network interface information
func detectNetwork() []collectors.NetworkInfo {
	var networks []collectors.NetworkInfo

	interfaces, err := net.Interfaces()
	if err != nil {
		return networks
	}

	for _, iface := range interfaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		netInfo := collectors.NetworkInfo{
			Interface: iface.Name,
			MAC:       iface.HardwareAddr.String(),
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if ipNet.IP.To4() != nil {
				netInfo.IPv4 = ipNet.IP.String()
			} else if ipNet.IP.To16() != nil {
				netInfo.IPv6 = ipNet.IP.String()
			}
		}

		// Only add if we have at least an IPv4 address
		if netInfo.IPv4 != "" || netInfo.IPv6 != "" {
			networks = append(networks, netInfo)
		}
	}

	return networks
}

// getLocalIP returns the primary local IP address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return ""
}
