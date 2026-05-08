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
			log.Printf("scrape feed %q failed: %v", feed.Name, err)
			continue
		}
		total += count
	}

	return total, nil
}

func (s *Scraper) scrapeFeed(ctx context.Context, feed FeedConfig) (int, error) {
	parsedFeed, err := s.parser.ParseURLWithContext(feed.URL, ctx)
	if err != nil {
		return 0, fmt.Errorf("parse feed: %w", err)
	}

	count := 0
	for _, item := range parsedFeed.Items {
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
