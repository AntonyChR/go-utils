package utils

import (
	"fmt"
	"net"
)

// GetLocalIPv4 retrieves the first non-loopback IPv4 address of the local machine.
// It iterates through the network interfaces, skipping inactive or loopback interfaces,
// and checks their associated addresses to find a valid IPv4 address.
//
// Returns:
// - A string containing the IPv4 address if found.
// - An error if no IPv4 address is found or if there is an issue accessing the network interfaces.
func GetLocalIPv4() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// Ignore invactive interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no se encontró dirección IPv4 en la red local")
}
