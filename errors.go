package fetch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type HTTPError struct {
	StatusCode int
	Status     string
	Body       []byte
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("%s: %s", e.Status, string(e.Body))
}

func checkHTTPResponse(res *http.Response) error {
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		res.Body = io.NopCloser(bytes.NewReader(body))

		return HTTPError{StatusCode: res.StatusCode, Status: res.Status, Body: body}
	}

	return nil
}
