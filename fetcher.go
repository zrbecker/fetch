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
	var options = FetchOptions{
		Method: http.MethodGet,
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

	u.RawQuery = options.Params.Encode()

	req, err := http.NewRequestWithContext(ctx, options.Method, u.String(), options.Body)
	if err != nil {
		return err
	}
	req.Header = options.Header

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
