package middleware

import (
	"fmt"
	httpnet "net/http"
	"net/url"

	"github.com/iamsorryprincess/project-layout/internal/http"
	"github.com/iamsorryprincess/project-layout/internal/http/request"
	"github.com/iamsorryprincess/project-layout/internal/http/response"
)

func Cors(next http.HandlerFunc) http.HandlerFunc {
	return func(request *request.Request, response *response.Response) {
		ref := `*`
		if refURL, err := url.Parse(request.Referer()); err == nil {
			ref = fmt.Sprintf("%s://%s", refURL.Scheme, refURL.Host)
		}

		response.Header().Set("Access-Control-Allow-Origin", ref)
		response.Header().Set("Access-Control-Allow-Credentials", "true")
		response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")

		if request.Method == httpnet.MethodOptions {
			response.Status(httpnet.StatusOK)
			return
		}

		next(request, response)
	}
}
