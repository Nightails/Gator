package cli

import (
	"context"
	"errors"
	"fmt"
)

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("too many arguments")
	}

	fmt.Println("listing feeds:")
	ctx := context.Background()

	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("- %s\n", feed.Name)
		fmt.Printf("- url: %s\n", feed.Url)
		// get the user that created the feed
		user, err := s.db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("- created by: %s\n", user.Name)
		fmt.Println()
	}

	return nil
}
