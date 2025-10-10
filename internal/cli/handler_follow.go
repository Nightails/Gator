package cli

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nightails/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing url")
	}

	ctx := context.Background()

	user, _ := s.db.GetUserByName(ctx, s.cfg.UserName)
	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return errors.New("this feed does not exist")
	}
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

	fmt.Printf("%s follow follow follow %s\n", ffRecord.UserName, ffRecord.FeedName)
	return nil
}
