package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"frontend-rss-hub/scraper-api/internal/models"
)

type ResourceStore interface {
	UpsertResource(resource models.Resource) error
}

type Scraper struct {
	parser *gofeed.Parser
	store  ResourceStore
}

func New(store ResourceStore) *Scraper {
	return &Scraper{
		parser: gofeed.NewParser(),
		store:  store,
	}
}

func (s *Scraper) Scrape(ctx context.Context, feeds []FeedConfig) (int, error) {
	total := 0

	for _, feed := range feeds {
		count, err := s.scrapeFeed(ctx, feed)
		if err != nil {
			// A single broken feed should not block the rest of the weekly read list.
			log.Printf("scrape feed %q failed: %v", feed.Name, err)
			continue
		}
		total += count
	}

	return total, nil
}

func (s *Scraper) ScrapeConfig(ctx context.Context, config Config) (int, error) {
	staticCount, err := s.UpsertLinks(config.Links)
	if err != nil {
		return staticCount, err
	}

	feedCount, err := s.Scrape(ctx, config.Feeds)
	return staticCount + feedCount, err
}

func (s *Scraper) UpsertLinks(links []LinkConfig) (int, error) {
	count := 0
	for _, link := range links {
		resource := models.Resource{
			Title:       strings.TrimSpace(link.Title),
			Link:        strings.TrimSpace(link.Link),
			Summary:     strings.TrimSpace(link.Summary),
			PublishedAt: link.PublishedAt.UTC(),
			SourceName:  strings.TrimSpace(link.SourceName),
			Tag:         strings.TrimSpace(link.Tag),
			CreatedAt:   time.Now().UTC(),
		}

		if err := s.store.UpsertResource(resource); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func (s *Scraper) scrapeFeed(ctx context.Context, feed FeedConfig) (int, error) {
	parsedFeed, err := s.parser.ParseURLWithContext(feed.URL, ctx)
	if err != nil {
		return 0, fmt.Errorf("parse feed: %w", err)
	}

	count := 0
	for _, item := range parsedFeed.Items {
		// Skip malformed feed entries; link/title are the minimum useful card data.
		if item == nil || item.Link == "" || item.Title == "" {
			continue
		}

		publishedAt := item.PublishedParsed
		if publishedAt == nil {
			publishedAt = item.UpdatedParsed
		}

		resource := models.Resource{
			Title:       strings.TrimSpace(item.Title),
			Link:        strings.TrimSpace(item.Link),
			Summary:     summaryFor(item),
			PublishedAt: timeOrNow(publishedAt),
			SourceName:  feed.Name,
			Tag:         feed.Tag,
			CreatedAt:   time.Now().UTC(),
		}

		if err := s.store.UpsertResource(resource); err != nil {
			return count, err
		}
		count++
	}

	return count, nil
}

func summaryFor(item *gofeed.Item) string {
	// Prefer RSS descriptions because they are usually shorter snippets than content.
	if strings.TrimSpace(item.Description) != "" {
		return strings.TrimSpace(item.Description)
	}
	return strings.TrimSpace(item.Content)
}

func timeOrNow(value *time.Time) time.Time {
	if value == nil || value.IsZero() {
		return time.Now().UTC()
	}
	return value.UTC()
}
