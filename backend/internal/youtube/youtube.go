package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
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

const BASE_API_URL = "https://www.googleapis.com"
const BASE_ATOM_FEED_URL = "https://www.youtube.com/feeds/videos.xml?channel_id=%s"

func GenerateAtomFeedURL(ctx context.Context, handle string, viperConfig *viper.Viper) (string, error) {
	apiKey := viperConfig.GetString("YOUTUBE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("`YOUTUBE_API_KEY` must be defined in config.")
	}

	// Parameters are encrypted, so not the end of the world to include the `apiKey`.
	url := fmt.Sprintf("%s/youtube/v3/channels/?part=id&forHandle=%s&key=%s", BASE_API_URL, handle, apiKey)
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

	return fmt.Sprintf(BASE_ATOM_FEED_URL, channelId), nil
}
