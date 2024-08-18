package feed

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title string `xml:"title"`
	Link  Link   `xml:"link"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

func ValidateAtomFeed(ctx context.Context, atomFeedURL string) error {
	if _, err := parseAtomFeed(ctx, atomFeedURL); err != nil {
		return fmt.Errorf("Error parsing atom feed: %w", err)
	}

	return nil
}

func ListEntry(ctx context.Context, atomFeedURL string, numEntries int) ([]Entry, error) {
	parsedFeed, err := parseAtomFeed(ctx, atomFeedURL)
	if err != nil {
		return nil, fmt.Errorf("Error parsing atom feed: %w", err)
	}

	return parsedFeed.Entries[:numEntries], nil
}

func parseAtomFeed(ctx context.Context, atomFeedURL string) (*Feed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", atomFeedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating new request w/ context: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}

	var feed Feed
	if err = xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("Error parsing Atom feed: %w", err)
	}

	return &feed, nil
}
