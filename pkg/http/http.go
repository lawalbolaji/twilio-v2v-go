package http

import (
	"io"
	"net/http"
)

type Client struct{}

func (c *Client) Post(url string, payload io.Reader, headers map[string]string) (int, []byte, error) {
	r, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return 0, nil, err
	}

	/* populate request headers */
	for key, val := range headers {
		r.Header.Add(key, val)
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return 0, nil, err
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}

	return res.StatusCode, b, nil
}
