package utils

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/mssola/user_agent"
)

// GetAgent returns raw User-Agent string
func GetAgent(r *http.Request) string {
	return r.Header.Get("User-Agent")
}

// GetUserAgent returns parsed Browser and OS
func GetUserAgent(r *http.Request) (string, string) {
	uaString := r.Header.Get("User-Agent")
	ua := user_agent.New(uaString)

	// Ambil nama browser dan versi
	browserName, _ := ua.Browser()
	// Ambil OS
	os := ua.OS()

	return browserName, os
}

// GetClientIP returns the real IP address of the client
func GetClientIP(r *http.Request) string {
	// Try to get the IP from the X-Forwarded-For header (when behind a proxy)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// The header could contain multiple IPs, so split by comma and return the first one
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Try to get the IP from the X-Real-IP header (common in reverse proxy setups)
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to the remote address from the request's RemoteAddr field
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// GetOriginFromRef parses referer url
func GetOriginFromRef(ref string) (string, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", nil
	}

	// Jika user mengirim tanpa scheme
	if !strings.Contains(ref, "://") {
		if strings.HasPrefix(ref, "//") {
			ref = "https:" + ref
		} else {
			ref = "https://" + ref
		}
	}

	u, err := url.Parse(ref)
	if err != nil {
		return "", err
	}

	if u.Scheme == "" || u.Host == "" {
		return "", nil
	}

	return u.Scheme + "://" + u.Host, nil
}
