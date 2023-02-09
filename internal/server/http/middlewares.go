package http

import (
	"net"
	"net/http"
)

// SubnetClientFilter is a middleware, that allows requests only from trusted subnet
func SubnetClientFilter(subnet net.IPNet) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := net.ParseIP(r.RemoteAddr)
			if clientIP == nil {
				http.Error(w, "Unknown client's IP (is headers set?)", http.StatusForbidden)
				return
			}

			if !subnet.Contains(clientIP) {
				http.Error(w, "Client's IP is not in trusted subnet", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
