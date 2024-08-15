package utils

import (
	"net"
	"testing"
)

func isValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
  return parsedIP.To4() != nil
}

func TestGetLocalIP(t *testing.T){
    ip, err := GetLocalIP()
    if err != nil {
        t.Error(err)
    }
    if !isValidIPv4(ip){
        t.Errorf("invalid ip v4, got: %s", ip)
    }
}
