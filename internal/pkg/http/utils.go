package http

import "net/http"

func ParseIP(request *http.Request) string {
	ip := request.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	ip = request.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	return "127.0.0.1"
}
