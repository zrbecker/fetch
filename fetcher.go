package fetch

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Fetcher struct {
	baseURL string
	client  *http.Client
}

func NewFetcher(baseURL string, opts ...FetcherOption) *Fetcher {
	c := &Fetcher{baseURL: baseURL, client: http.DefaultClient}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Fetcher) Fetch(ctx context.Context, path string, response interface{}, opts ...FetchOption) error {
	var options = fetchOptions{
		method: http.MethodGet,
	}

	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return err
		}
	}

	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return err
	}

	params := url.Values{}
	for _, param := range options.params {
		params.Add(param.Key, param.Value)
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, options.method, u.String(), options.body)
	if err != nil {
		return err
	}

	for _, header := range options.headers {
		req.Header.Add(header.Key, header.Value)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if err := checkHTTPResponse(res); err != nil {
		return err
	}

	if response != nil {
		httpRes, ok := response.(*http.Response)
		if ok {
			*httpRes = *res
		} else {
			defer res.Body.Close()
			bodyBz, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(bodyBz, response); err != nil {
				return err
			}
		}
	}

	return nil
}
