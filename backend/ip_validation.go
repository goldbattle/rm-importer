package backend

import (
	"net"
)

func IsIpValid(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil
}
