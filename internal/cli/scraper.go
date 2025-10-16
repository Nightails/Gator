package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Nightails/gator/internal/database"
	"github.com/Nightails/gator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// scrapeFeeds scrapes the next feed to fetch from the database and saves the posts to the database.
func scrapeFeeds(s *state) error {
	feedToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	ctx := context.Background()
	_ = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		ID: feedToFetch.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})

	rssFeed, err := rss.FetchFeed(ctx, feedToFetch.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	rssFeed.UnescapeString()
	if err := savePostsToDB(rssFeed, feedToFetch.ID, s); err != nil {
		return err
	}

	return nil
}

// savePostsToDB saves the posts of the given RSS feed to the database.
func savePostsToDB(feed *rss.Feed, feedID uuid.UUID, s *state) error {
	for _, item := range feed.Channel.Item {
		publishedAt, err := parseTime(item.PubDate)
		if err != nil {
			fmt.Printf("Error: failed to parse publish date '%s': %v\n", item.PubDate, err)
			continue
		}

		if _, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      feedID,
		}); err != nil {
			// Check if the error is a unique constraint violation on URL
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				// Silently skip duplicate URLs
				continue
			}
			// Log all other errors
			fmt.Printf("Error: failed to save post: %v\n", err)
		}
	}

	return nil
}

// parseTime parses the given date string into a time.Time.
func parseTime(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		"2006-01-02T15:04:05Z07:00", // ISO 8601
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse date: %s", dateStr)
}

// printRSSFeed prints the given RSS feed to the console.
func printRSSFeed(rssFeed *rss.Feed) {
	fmt.Printf("Title: %s\n", rssFeed.Channel.Title)
	fmt.Printf("Description: %s\n", rssFeed.Channel.Description)
	fmt.Println("Items:")
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("- Title: %s\n", item.Title)
		fmt.Printf("- Link: %s\n", item.Link)
		fmt.Printf("- Description: %s\n", item.Description)
		fmt.Println()
	}
}
