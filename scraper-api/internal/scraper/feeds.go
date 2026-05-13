package scraper

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Tags  []TagConfig  `json:"tags"`
	Feeds []FeedConfig `json:"feeds"`
	Links []LinkConfig `json:"links"`
}

type TagConfig struct {
	Tag   string `json:"tag"`
	Label string `json:"label"`
}

type FeedConfig struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Tag  string `json:"tag"`
}

type LinkConfig struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Summary     string    `json:"summary"`
	SourceName  string    `json:"source_name"`
	Tag         string    `json:"tag"`
	PublishedAt time.Time `json:"published_at"`
}

func LoadConfig(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read resource config: %w", err)
	}

	if strings.HasPrefix(strings.TrimSpace(string(content)), "[") {
		var feeds []FeedConfig
		if err := json.Unmarshal(content, &feeds); err != nil {
			return Config{}, fmt.Errorf("parse legacy feed config: %w", err)
		}

		config := Config{
			Tags:  tagsFromFeeds(feeds),
			Feeds: feeds,
		}
		if err := validateConfig(config, false); err != nil {
			return Config{}, err
		}
		return config, nil
	}

	var config Config
	if err := json.Unmarshal(content, &config); err != nil {
		return Config{}, fmt.Errorf("parse resource config: %w", err)
	}
	if err := validateConfig(config, true); err != nil {
		return Config{}, err
	}

	return config, nil
}

// LoadFeeds reads the user-editable feed list. Keeping this as JSON makes it
// easy for the team to add or remove sources without recompiling the API.
func LoadFeeds(path string) ([]FeedConfig, error) {
	config, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}
	return config.Feeds, nil
}

func validateConfig(config Config, requireTags bool) error {
	tagSet := map[string]bool{}
	for i, tag := range config.Tags {
		if strings.TrimSpace(tag.Tag) == "" || strings.TrimSpace(tag.Label) == "" {
			return fmt.Errorf("tag config at index %d must include tag and label", i)
		}
		if tagSet[tag.Tag] {
			return fmt.Errorf("duplicate tag %q", tag.Tag)
		}
		tagSet[tag.Tag] = true
	}

	if requireTags && len(config.Tags) == 0 {
		return fmt.Errorf("resource config must include at least one tag")
	}

	for i, feed := range config.Feeds {
		if feed.Name == "" || feed.URL == "" || feed.Tag == "" {
			return fmt.Errorf("feed config at index %d must include name, url, and tag", i)
		}
		if requireTags && !tagSet[feed.Tag] {
			return fmt.Errorf("feed config at index %d references unknown tag %q", i, feed.Tag)
		}
	}

	for i, link := range config.Links {
		if link.Title == "" || link.Link == "" || link.Summary == "" || link.SourceName == "" || link.Tag == "" || link.PublishedAt.IsZero() {
			return fmt.Errorf("link config at index %d must include title, link, summary, source_name, tag, and published_at", i)
		}
		if requireTags && !tagSet[link.Tag] {
			return fmt.Errorf("link config at index %d references unknown tag %q", i, link.Tag)
		}
	}

	return nil
}

func tagsFromFeeds(feeds []FeedConfig) []TagConfig {
	seen := map[string]bool{}
	tags := []TagConfig{}
	for _, feed := range feeds {
		if feed.Tag == "" || seen[feed.Tag] {
			continue
		}
		seen[feed.Tag] = true
		tags = append(tags, TagConfig{
			Tag:   feed.Tag,
			Label: labelFromTag(feed.Tag),
		})
	}
	return tags
}

func labelFromTag(tag string) string {
	switch tag {
	case "ai":
		return "AI"
	case "aws":
		return "AWS"
	case "mcp":
		return "MCP"
	case "nextjs":
		return "Next.js"
	default:
		words := strings.Fields(strings.ReplaceAll(tag, "-", " "))
		for i, word := range words {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
		return strings.Join(words, " ")
	}
}
