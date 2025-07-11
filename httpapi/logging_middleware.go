package httpapi

import (
	"log"
	"net"
	"net/http"
	"strings"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HTTP %s %s from %s", r.Method, r.URL.Path, getRemoteAddress(r))
		next.ServeHTTP(w, r)
	})
}

func getRemoteAddress(r *http.Request) string {
	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		return strings.TrimSpace(strings.Split(forwardedFor, ",")[0])
	}

	realIp := r.Header.Get("X-Real-IP")
	if realIp != "" {
		return realIp
	}

	if r.RemoteAddr != "" {
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			return host
		} else {
			return r.RemoteAddr
		}
	}

	return "unknown"
}
