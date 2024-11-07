package fetch

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type FetcherOption func(*Fetcher) error

func Client(client *http.Client) FetcherOption {
	return func(f *Fetcher) error {
		f.client = client
		return nil
	}
}

type FetchOption func(*FetchOptions) error

type FetchOptions struct {
	Method string
	Params url.Values
	Header http.Header
	Body   io.Reader
}

type Param struct {
	Key   string
	Value string
}

type Header struct {
	Key   string
	Value string
}

func Method(method string) FetchOption {
	return func(options *FetchOptions) error {
		options.Method = method
		return nil
	}
}

func Params(params ...Param) FetchOption {
	return func(options *FetchOptions) error {
		for _, param := range params {
			options.Params.Add(param.Key, param.Value)
		}
		return nil
	}
}

func Headers(headers ...Header) FetchOption {
	return func(options *FetchOptions) error {
		if options.Header == nil {
			options.Header = make(http.Header)
		}
		for _, header := range headers {
			options.Header.Add(header.Key, header.Value)
		}
		return nil
	}
}

func Body(body interface{}) FetchOption {
	return func(options *FetchOptions) error {
		bodyReader, ok := body.(io.Reader)
		if !ok {
			if err := Headers(Header{
				Key:   "Content-Type",
				Value: "application/json; charset=UTF-8",
			})(options); err != nil {
				return err
			}

			jsonBz, err := json.Marshal(body)
			if err != nil {
				return err
			}

			bodyReader = io.NopCloser(bytes.NewReader(jsonBz))
		}
		options.Body = bodyReader
		return nil
	}
}
