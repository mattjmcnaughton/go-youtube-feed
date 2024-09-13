package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ChannelData represents the structure of the response we expect from YouTube
type ChannelData struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Id string `json:"id"`
	} `json:"items"`
}

const (
	BASE_API_URL       = "https://www.googleapis.com"
	BASE_ATOM_FEED_URL = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"
)

type YoutubeClient struct {
	apiKey          string
	baseApiUrl      string
	baseAtomFeedUrl string
}

func NewYoutubeClient(apiKey string) *YoutubeClient {
	return &YoutubeClient{
		apiKey:          apiKey,
		baseApiUrl:      BASE_API_URL,
		baseAtomFeedUrl: BASE_ATOM_FEED_URL,
	}
}

func (yc *YoutubeClient) GenerateAtomFeedURL(ctx context.Context, handle string) (string, error) {
	// Parameters are encrypted, so not the end of the world to include the `apiKey`.
	url := fmt.Sprintf("%s/youtube/v3/channels/?part=id&forHandle=%s&key=%s", yc.baseApiUrl, handle, yc.apiKey)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating new request w/ context: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %w", err)
	}

	var data ChannelData
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("Error decoding response: %w", err)
	}

	if data.PageInfo.TotalResults > 1 {
		return "", fmt.Errorf("Found %d results, and only expected 1", data.PageInfo.TotalResults)
	}

	channelId := data.Items[0].Id

	return fmt.Sprintf(yc.baseAtomFeedUrl, channelId), nil
}
