package scraper

import (
	"encoding/json"
	"fmt"
	"os"
)

type FeedConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Tag  string `json:"tag"`
}

func LoadFeeds(path string) ([]FeedConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read feed config: %w", err)
	}

	var feeds []FeedConfig
	if err := json.Unmarshal(content, &feeds); err != nil {
		return nil, fmt.Errorf("parse feed config: %w", err)
	}

	for i, feed := range feeds {
		if feed.Name == "" || feed.URL == "" || feed.Tag == "" {
			return nil, fmt.Errorf("feed config at index %d must include name, url, and tag", i)
		}
	}

	return feeds, nil
}
