package networks

import (
	"net/http"
	"strings"
)

type WebRequest struct {
	*http.Request
	ContentType     string
	Accept          string
	AcceptLanguages WebAcceptLanguages
}

func (r *WebRequest) GetRequest() *http.Request {
	return r.Request
}

func (r *WebRequest) Apply(req *http.Request) {
	r.Request = req
	r.Reset()

	// X-HTTP-Method-Override
	if method := r.Header.Get("X-HTTP-Method-Override"); method != "" && r.Method == "POST" {
		r.Method = method
	}
}

func (r *WebRequest) Reset() {
	r.ContentType = ResolveContentType(r.Request)
	r.Accept = r.Request.Header.Get("accept")
	r.AcceptLanguages = ResolveAcceptLanguage(r.Request)
}

func NewWebRequest() *WebRequest {
	return &WebRequest{}
}

func ResolveContentType(r *http.Request) string {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return "text/html"
	}

	return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}
