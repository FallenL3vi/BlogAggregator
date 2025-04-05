package main

import (
	"encoding/xml"
	"net/http"
	"errors"
	"io"
	"html"
	"context"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(context.Background(), "GET", feedURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	var DefaultClient = &http.Client{}

	response, err := DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("Resposne code not OK")
	}

	body, err := io.ReadAll(response.Body)

	response.Body.Close()

	if err != nil {
		return nil, err
	}

	var feed RSSFeed

	err = xml.Unmarshal(body, &feed)

	if err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item:= range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}

	return &feed, nil

}