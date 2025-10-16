package cli

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Nightails/gator/internal/database"
	"github.com/Nightails/gator/internal/rss"
)

func scrapeFeeds(db *database.Queries) error {
	feedToFetch, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	ctx := context.Background()
	_ = db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
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
	printRSSFeed(rssFeed)

	return nil
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
