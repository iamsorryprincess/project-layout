package http

import (
	"fmt"
	"net/http"
	"net/url"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ref := `*`
		if refURL, err := url.Parse(request.Referer()); err == nil {
			ref = fmt.Sprintf("%s://%s", refURL.Scheme, refURL.Host)
		}

		writer.Header().Set("Access-Control-Allow-Origin", ref)
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")

		if request.Method == http.MethodOptions {
			writer.WriteHeader(200)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
