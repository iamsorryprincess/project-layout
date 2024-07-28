package middleware

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/iamsorryprincess/project-layout/internal/http/httproute"
)

func Cors(next httproute.HandlerFunc) httproute.HandlerFunc {
	return func(request *httproute.Request, response *httproute.Response) {
		ref := `*`
		if refURL, err := url.Parse(request.Referer()); err == nil {
			ref = fmt.Sprintf("%s://%s", refURL.Scheme, refURL.Host)
		}

		response.Header().Set("Access-Control-Allow-Origin", ref)
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")

		if request.Method == http.MethodOptions {
			response.Status(http.StatusOK)
			return
		}

		next(request, response)
	}
}
