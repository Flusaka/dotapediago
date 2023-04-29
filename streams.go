package dotapediago

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
)

func (client *Client) GetStreamURL(liquipediaUrl string) (string, error) {
	// If it exists in the cache, return that
	streamUrl, ok := client.streamCache.GetResolvedStreamURL(liquipediaUrl)
	if !ok {
		// If not found, start executing the http calls to get the stream URL, then add it to cache

		// Load the "special" stream page
		pageData, err := client.getPage(liquipediaUrl)
		if err != nil {
			return "", err
		}

		reader := bytes.NewReader(pageData)
		doc, err := goquery.NewDocumentFromReader(reader)

		if err != nil {
			return "", err
		}

		// First one is the Twitch stream
		streamPortal := doc.Find("iframe").First()
		src, exists := streamPortal.Attr("src")
		if !exists {
			return "", errors.New("no iframe found")
		}

		// Parse the url and get the "channel" query param
		stream, err := url.Parse(src)
		if err != nil {
			return "", errors.New("url not valid")
		}

		queryMap := stream.Query()
		channel := queryMap.Get("channel")
		streamUrl = fmt.Sprintf("https://twitch.tv/%s", channel)
		client.streamCache.SetResolvedStreamURL(liquipediaUrl, streamUrl)
	}
	return streamUrl, nil
}
