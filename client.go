package dotapediago

import "net/http"

type Client struct {
	httpClient *http.Client
}

type customTransport struct {
	headers             map[string]string
	underlyingTransport http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range t.headers {
		req.Header.Set(key, value)
	}
	return t.underlyingTransport.RoundTrip(req)
}

func NewClient(userAgent string) *Client {
	client := &http.Client{
		Transport: &customTransport{
			headers: map[string]string{
				"User-Agent":      userAgent,
				"Accept-Encoding": "gzip",
			},
			underlyingTransport: http.DefaultTransport,
		},
	}
	return &Client{
		httpClient: client,
	}
}
