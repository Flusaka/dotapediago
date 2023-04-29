package dotapediago

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"github.com/Flusaka/dotapediago/cache"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	Root struct {
		Text struct {
			Document string `json:"*"`
		} `json:"text"`
	} `json:"parse"`
}

type Client struct {
	httpClient  *http.Client
	streamCache *cache.StreamCache
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
		httpClient:  client,
		streamCache: cache.NewStreamCache(),
	}
}

func (client *Client) getPage(endpoint string) ([]byte, error) {
	res, err := client.httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}

	var responseReader io.Reader = res.Body
	if res.Header.Get("Content-Encoding") == "gzip" {
		responseReader, err = gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
	}

	return io.ReadAll(responseReader)
}

func (client *Client) getDocument(endpoint string) (*goquery.Document, error) {
	responseData, err := client.getPage(endpoint)
	if err != nil {
		return nil, err
	}

	var parsedResponse Response
	err = json.Unmarshal(responseData, &parsedResponse)
	if err != nil {
		return nil, err
	}

	reader := strings.NewReader(parsedResponse.Root.Text.Document)
	return goquery.NewDocumentFromReader(reader)
}
