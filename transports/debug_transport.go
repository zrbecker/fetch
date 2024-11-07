package transports

import (
	"bytes"
	"io"
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
	log.Println("Request Method:", req.Method)
	log.Println("Request URL:", req.URL.String())
	for key, values := range req.Header {
		log.Printf("Request Header %s: %s", key, strings.Join(values, ","))
	}
	if req.Body != nil {
		defer req.Body.Close()
		bodyBz, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		log.Println("Request Body:", string(bodyBz))

		req.Body = io.NopCloser(bytes.NewReader(bodyBz))
	}

	res, err := rt.wrappedTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	log.Println("Response Status:", res.Status)
	for key, values := range req.Header {
		log.Printf("Response Header %s: %s", key, strings.Join(values, ","))
	}
	if res.Body != nil {
		defer res.Body.Close()
		bodyBz, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		log.Println("Response Body:", string(bodyBz))

		res.Body = io.NopCloser(bytes.NewReader(bodyBz))
	}

	return res, nil
}
