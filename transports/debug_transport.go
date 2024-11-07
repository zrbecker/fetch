package transports

import (
	"log"
	"net/http"
	"strings"
)

type DebugTransport struct {
	wrappedTransport http.RoundTripper
}

func NewDebugTransport(wrappedTransport http.RoundTripper) *DebugTransport {
	return &DebugTransport{wrappedTransport}
}

func (rt *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Println("Method:", req.Method)
	log.Println("URL:", req.URL.String())
	for key, values := range req.Header {
		log.Printf("Header %s: %s", key, strings.Join(values, ","))
	}
	return rt.wrappedTransport.RoundTrip(req)
}
