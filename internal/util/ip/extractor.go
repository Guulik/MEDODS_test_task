package ip

import (
	"net"
	"net/http"
)

func GetIPAddress(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		if net.ParseIP(forwarded) != nil {
			return forwarded
		} else {
			return ""
		}
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
