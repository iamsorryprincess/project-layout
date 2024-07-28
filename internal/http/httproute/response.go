package httproute

import (
	"encoding/json"
	"net/http"

	"github.com/iamsorryprincess/project-layout/internal/log"
)

type Response struct {
	logger log.Logger
	writer http.ResponseWriter
}

func newResponse(writer http.ResponseWriter, logger log.Logger) *Response {
	return &Response{
		writer: writer,
		logger: logger,
	}
}

func (r *Response) Status(statusCode int) {
	r.writer.WriteHeader(statusCode)
}

func (r *Response) Header() http.Header {
	return r.writer.Header()
}

func (r *Response) Write(data []byte) (int, error) {
	return r.writer.Write(data)
}

func (r *Response) Text(code int, text string) error {
	r.writer.Header().Set("Content-Type", "text/plain")
	r.writer.WriteHeader(code)
	_, err := r.writer.Write([]byte(text))
	return err
}

func (r *Response) JSON(code int, data interface{}) error {
	r.writer.Header().Set("Content-Type", "application/json")
	r.writer.WriteHeader(code)
	return json.NewEncoder(r.writer).Encode(data)
}
