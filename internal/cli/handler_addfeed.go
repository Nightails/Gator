package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nightails/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("missing feed name and url")
	}

	ctx := context.Background()

	// get current user
	user, err := s.db.GetUserByName(ctx, s.cfg.UserName)
	if err != nil {
		return err
	}

	// create the feed
	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(ctx, feedParams)
	if err != nil {
		return err
	}

	fmt.Printf("added feed: %s\n", feed.Name)
	fmt.Printf("feed url: %s\n", feed.Url)
	fmt.Printf("created at: %v\n", feed.CreatedAt)

	// follow the feed
	ffParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	ffRecord, err := s.db.CreateFeedFollow(ctx, ffParams)
	if err != nil {
		return errors.New("unable to follow this feed")
	}
	fmt.Printf("%s follow %s\n", ffRecord.UserName, ffRecord.FeedName)

	return nil
}
