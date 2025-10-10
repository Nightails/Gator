package cli

import (
	"context"
	"errors"
	"fmt"
)

func handlerFollowing(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return errors.New("too many arguments")
	}

	ctx := context.Background()

	user, _ := s.db.GetUserByName(ctx, s.cfg.UserName)
	feeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return errors.New("unable to retrieve following feeds")
	}

	fmt.Printf("%s following:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("- %s\n", feed.FeedName)
	}

	return nil
}
