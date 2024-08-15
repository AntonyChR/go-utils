package utils

import (
	"errors"
	"net"
)

func GetLocalIP() (string, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("error getting local ip address")

}
