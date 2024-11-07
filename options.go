package fetch

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type FetcherOption func(*Fetcher) error

func Client(client *http.Client) FetcherOption {
	return func(f *Fetcher) error {
		f.client = client
		return nil
	}
}

type FetchOption func(*fetchOptions) error

type fetchOptions struct {
	method  string
	params  []Param
	headers []Header
	body    io.Reader
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
	return func(options *fetchOptions) error {
		options.method = method
		return nil
	}
}

func Params(params ...Param) FetchOption {
	return func(options *fetchOptions) error {
		options.params = append(options.params, params...)
		return nil
	}
}

func Headers(headers ...Header) FetchOption {
	return func(options *fetchOptions) error {
		options.headers = append(options.headers, headers...)
		return nil
	}
}

func Body(body interface{}) FetchOption {
	return func(options *fetchOptions) error {
		bodyReader, ok := body.(io.Reader)
		if !ok {
			jsonBz, err := json.Marshal(body)
			if err != nil {
				return err
			}
			bodyReader = io.NopCloser(bytes.NewReader(jsonBz))
		}
		options.body = bodyReader
		return nil
	}
}
