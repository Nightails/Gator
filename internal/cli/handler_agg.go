package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nightails/gator/internal/rss"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("too many arguments")
	}

	// temporary url, to be replaced with config
	url := "https://www.wagslane.dev/index.xml"
	ctx := context.Background()
	rssf, err := rss.FetchFeed(ctx, url)
	if err != nil {
		return err
	}
	rssf.UnescapeString()
	printRSSFeed(rssf)

	return nil
}

func printRSSFeed(rssf *rss.RSSFeed) {
	fmt.Printf("Title: %s\n", rssf.Channel.Title)
	fmt.Printf("Description: %s\n", rssf.Channel.Description)
	fmt.Println("Items:")
	for _, item := range rssf.Channel.Item {
		fmt.Printf("- Title: %s\n", item.Title)
		fmt.Printf("- Link: %s\n", item.Link)
		fmt.Printf("- Description: %s\n", item.Description)
		fmt.Println()
	}
}
