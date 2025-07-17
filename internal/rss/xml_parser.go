package rss

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, bytes.NewBuffer(make([]byte, 0)))
	if err != nil {
		return nil, fmt.Errorf("unable to form the http request: %v", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error when making the request: %v", err)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %v", err)
	}

	var result RSSFeed
	if err := xml.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("unable to interpret response: %v", err)
	}

	cleanFeed(&result)

	return &result, nil
}

func cleanFeed(feed *RSSFeed) error {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for _, item := range feed.Channel.Item {
		item.Description = html.UnescapeString(item.Description)
		item.Title = html.UnescapeString(item.Title)
	}

	return nil
}

func (f RSSFeed) PrintFeed() error {
	fmt.Println("\n========================================")
	fmt.Printf("\t%v:\n", f.Channel.Title)
	fmt.Printf("\t\t%v\n\n", f.Channel.Description)

	for _, item := range f.Channel.Item {
		fmt.Printf("\t*\t%v (%v) - \n\t\t%v\n\n", item.Title, item.PubDate, item.Description)
	}

	return nil
}