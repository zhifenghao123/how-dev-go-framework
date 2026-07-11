package ihttp

import "net/http"

// Response struct of response
type Response struct {
	Code   int
	Data   interface{}
	Header map[string]string
	Writer http.ResponseWriter
}

// SetHeader set header of response
func (r *Response) SetHeader(key string, value string) {
	r.Header[key] = value
}
